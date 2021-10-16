package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrProviderServiceNil = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
	ErrRepositoryNil      = errors.New("repository cannot be nil")
)

type SyncService struct {
	ctx              context.Context
	mu               *sync.RWMutex
	provGroupsFilter []string
	provUsersFilter  []string
	prov             IdentityProviderService
	scim             SCIMService
	repo             StateRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, repo StateRepository, opts ...SyncServiceOption) (*SyncService, error) {
	if prov == nil {
		return nil, ErrProviderServiceNil
	}
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if repo == nil {
		return nil, ErrRepositoryNil
	}

	ss := &SyncService{
		ctx:              ctx,
		mu:               &sync.RWMutex{},
		prov:             prov,
		provGroupsFilter: []string{}, // fill in with the opts
		provUsersFilter:  []string{}, // fill in with the opts
		scim:             scim,
		repo:             repo,
	}

	for _, opt := range opts {
		opt(ss)
	}

	return ss, nil
}

func (ss *SyncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	log := log.WithFields(log.Fields{
		"filter": ss.provGroupsFilter,
	})

	// retrive data from the identity provider
	// always theses steps are necessary
	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting groups from the identity provider: %w", err)
	}

	if pGroupsResult.Items == 0 {
		log.Warn("identity provider groups empty")
		return nil
	}

	pUsersResult, pGroupsUsersResult, err := ss.prov.GetUsersAndGroupsUsers(ss.ctx, pGroupsResult)
	if err != nil {
		return fmt.Errorf("error getting users and groups and their users: %w", err)
	}

	if pGroupsUsersResult.Items == 0 {
		log.Warn("identity provider groups members empty")
	}

	if pUsersResult.Items == 0 {
		log.Warn("identity provider users empty")
	}

	// get the state metadata from the repository
	state, err := ss.repo.GetState(ss.ctx)
	if err != nil {
		return fmt.Errorf("error getting state from the repository: %w", err)
	}

	// first time syncing
	if state.LastSync == "" {

		log.Info("state without lastsync time, first time syncing")

		// Check SCIM side to see if there are any elelemnts to be
		// reconciled. I mean SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side
		// in case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side.
		sGroupsResult, err := ss.scim.GetGroups(ss.ctx)
		if err != nil {
			return fmt.Errorf("error getting groups from the SCIM service: %w", err)
		}

		if sGroupsResult.Items == 0 {
			log.Info("no SCIM groups to be reconciling")
		} else {
			log.Warn("reconciling the SCIM data with the Identity Provider data, the first syncing")

			// reconciling elements
			gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, sGroupsResult)

			if gCreate.Items == 0 {
				log.Info("no groups to be created")
			} else {
				if err := ss.scim.CreateGroups(ss.ctx, gCreate); err != nil {
					return fmt.Errorf("error creating groups in SCIM Provider: %w", err)
				}
			}

			if gUpdate.Items == 0 {
				log.Info("no groups to be updated")
			} else {
				if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
					return fmt.Errorf("error updating groups in SCIM Provider: %w", err)
				}
			}

			if gDelete.Items == 0 {
				log.Info("no groups to be deleted")
			} else {
				if err := ss.scim.DeleteGroups(ss.ctx, gDelete); err != nil {
					return fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
				}
			}
		}

		// the groups here could be empty, but we need to check if users exist even if there are not groups
		// here I return all the users and not only the members of groups
		// becuase the users and groups in the scim needs to be controlled by
		// the sync process
		sUsersResult, sGroupsUsersResult, err := ss.scim.GetUsersAndGroupsUsers(ss.ctx)
		if err != nil {
			return fmt.Errorf("error getting users and groups and their users from SCIM provider: %w", err)
		}

		if sUsersResult.Items == 0 {
			log.Info("no SCIM users to be reconciling")
		} else {

			// reconciling users
			uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, sUsersResult)

			if uCreate.Items == 0 {
				log.Info("no users to be created")
			} else {
				if err := ss.scim.CreateUsers(ss.ctx, uCreate); err != nil {
					return fmt.Errorf("error creating users in SCIM Provider: %w", err)
				}
			}

			if uUpdate.Items == 0 {
				log.Info("no users to be updated")
			} else {
				if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
					return fmt.Errorf("error updating users in SCIM Provider: %w", err)
				}
			}

			if uDelete.Items == 0 {
				log.Info("no users to be deleted")
			} else {
				if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
					return fmt.Errorf("error deleting users in SCIM Provider: %w", err)
				}
			}
		}

		if sGroupsUsersResult.Items == 0 {
			log.Info("no SCIM groups and members to be reconciling")
		} else {
			// reconciling groups-users --> groups members
			ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, sGroupsUsersResult)

			if ugCreate.Items == 0 {
				log.Info("no groups to be joined")
			} else {
				if err := ss.scim.CreateMembers(ss.ctx, ugCreate); err != nil {
					return fmt.Errorf("error creating groups members in SCIM Provider: %w", err)
				}
			}

			if ugDelete.Items == 0 {
				log.Info("no groups to be removed")
			} else {
				if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
					return fmt.Errorf("error deleting groups members in SCIM Provider: %w", err)
				}
			}
		}

	} else { // This is not the first time syncing
		log.WithField("lastsync", state.LastSync).Info("state with lastsync time, it is not first time syncing")

		if pGroupsResult.HashCode == state.Resources.Groups.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")
		} else {
			log.Info("provider groups and state groups are diferent")
			// now here we have the google fresh data and the last sync data state
			// we need to compare the data and decide what to do
			// see differences between the two data sets

			gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, &state.Resources.Groups)

			if gCreate.Items == 0 {
				log.Info("no groups to be created")
			} else {
				if err := ss.scim.CreateGroups(ss.ctx, gCreate); err != nil {
					return fmt.Errorf("error creating groups in SCIM Provider: %w", err)
				}
			}

			if gUpdate.Items == 0 {
				log.Info("no groups to be updated")
			} else {
				if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
					return fmt.Errorf("error updating groups in SCIM Provider: %w", err)
				}
			}

			if gDelete.Items == 0 {
				log.Info("no groups to be deleted")
			} else {
				if err := ss.scim.DeleteGroups(ss.ctx, gDelete); err != nil {
					return fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
				}
			}
		}

		if pUsersResult.HashCode == state.Resources.Users.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")
		} else {
			log.Info("provider users and state users are diferent")

			uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, &state.Resources.Users)

			if uCreate.Items == 0 {
				log.Info("no users to be created")
			} else {
				if err := ss.scim.CreateUsers(ss.ctx, uCreate); err != nil {
					return fmt.Errorf("error creating users in SCIM Provider: %w", err)
				}
			}

			if uUpdate.Items == 0 {
				log.Info("no users to be updated")
			} else {
				if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
					return fmt.Errorf("error updating users in SCIM Provider: %w", err)
				}
			}

			if uDelete.Items == 0 {
				log.Info("no users to be deleted")
			} else {
				if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
					return fmt.Errorf("error deleting users in SCIM Provider: %w", err)
				}
			}
		}

		if pGroupsUsersResult.HashCode == state.Resources.GroupsUsers.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")
		} else {
			log.Info("provider groups-members and state groups-members are diferent")

			ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, &state.Resources.GroupsUsers)

			if ugCreate.Items == 0 {
				log.Info("no groups-members to be created")
			} else {
				if err := ss.scim.CreateMembers(ss.ctx, ugCreate); err != nil {
					return fmt.Errorf("error creating groups members in SCIM Provider: %w", err)
				}
			}

			if ugDelete.Items == 0 {
				log.Info("no groups-members to be deleted")
			} else {
				if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
					return fmt.Errorf("error deleting groups members in SCIM Provider: %w", err)
				}
			}
		}
	}

	// after be sure all the SCIM side is aligned with the Identity Provider side
	// we can update the state with the identity provider data
	newState := &model.State{
		SchemaVersion: "1.0.0",
		CodeVersion:   "0.0.1",
		LastSync:      time.Now().Format(time.RFC3339),
		Resources: model.StateResources{
			Groups:      *pGroupsResult,
			Users:       *pUsersResult,
			GroupsUsers: *pGroupsUsersResult,
		},
	}

	if err := ss.repo.SaveState(ss.ctx, newState); err != nil {
		return fmt.Errorf("error saving state: %w", err)
	}

	return nil
}

func (ss *SyncService) SyncGroupsAndUsers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	pGroups, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting groups from provider: %w", err)
	}

	pUsers, err := ss.prov.GetUsers(ss.ctx, ss.provUsersFilter)
	if err != nil {
		return fmt.Errorf("error getting users from provider: %w", err)
	}

	if err := ss.scim.DeleteGroups(ss.ctx, pGroups); err != nil {
		return fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
	}

	if err := ss.scim.DeleteUsers(ss.ctx, pUsers); err != nil {
		return fmt.Errorf("error deleting users in SCIM Provider: %w", err)
	}

	return nil
}
