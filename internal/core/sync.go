package core

import (
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrNilContext                    = errors.New("context cannot be nil")
	ErrProviderServiceNil            = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil                = errors.New("SCIM service cannot be nil")
	ErrGettingGroups                 = errors.New("error getting groups")
	ErrRepositoryNil                 = errors.New("repository cannot be nil")
	ErrSyncStateNil                  = errors.New("sync state cannot be nil")
	ErrGettingState                  = errors.New("error getting state")
	ErrGettingGetUsersAndGroupsUsers = errors.New("error getting users and groups and their users")
	ErrGettingGroupsMetadata         = errors.New("error getting groups metadata")
	ErrGettingStateMetadata          = errors.New("error getting state metadata")
	ErrGettingStateGroups            = errors.New("error getting state groups")
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
	if ctx == nil {
		return nil, ErrNilContext
	}
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
		return ErrGettingGroups
	}

	if pGroupsResult.Items == 0 {
		log.Warn("identity provider groups empty")
		return nil
	}

	pUsersResult, pGroupsUsersResult, err := ss.prov.GetUsersAndGroupsUsers(ss.ctx, pGroupsResult)
	if err != nil {
		return ErrGettingGetUsersAndGroupsUsers
	}

	if pGroupsUsersResult.Items == 0 {
		log.Warn("identity provider groups members empty")
	}

	if pUsersResult.Items == 0 {
		log.Warn("identity provider users empty")
	}

	// get the state from the repository
	stateMetadata, err := ss.repo.GetStateMetadata()
	if err != nil {
		return ErrGettingStateMetadata
	}

	// first time syncing
	if stateMetadata.LastSync == "" {

		log.Info("state without lastsync time, first time syncing")

		// Check SCIM side to see if there are any elelemnts to be
		// reconciled. I mean SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side
		// in case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side.
		sGroupsResult, err := ss.scim.GetGroups(ss.ctx)
		if err != nil {
			return err
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
					return err
				}
			}

			if gUpdate.Items == 0 {
				log.Info("no groups to be updated")
			} else {
				if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
					return err
				}
			}

			if gDelete.Items == 0 {
				log.Info("no groups to be deleted")
			} else {
				if err := ss.scim.DeleteGroups(ss.ctx, gDelete); err != nil {
					return err
				}
			}
		}

		// the groups here could be empty, but we need to check if users exist even if there are not groups
		sUsersResult, sGroupsUsersResult, err := ss.scim.GetUsersAndGroupsUsers(ss.ctx, sGroupsResult)
		if err != nil {
			return err
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
					return err
				}
			}

			if uUpdate.Items == 0 {
				log.Info("no users to be updated")
			} else {
				if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
					return err
				}
			}

			if uDelete.Items == 0 {
				log.Info("no users to be deleted")
			} else {
				if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
					return err
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
					return err
				}
			}

			if ugDelete.Items == 0 {
				log.Info("no groups to be removed")
			} else {
				if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
					return err
				}
			}
		}

	} else { // This is not the first time syncing
		log.WithField("lastsync", stateMetadata.LastSync).Info("state with lastsync time, it is not first time syncing")

		stateGroupsMetadata, err := ss.repo.GetGroupsMetadata()
		if err != nil {
			return ErrGettingGroupsMetadata
		}

		if pGroupsResult.HashCode == stateGroupsMetadata.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")
		} else {
			log.Info("provider groups and state groups are diferent")
			// now here we have the google fresh data and the last sync data state
			// we need to compare the data and decide what to do
			// see differences between the two data sets
			stateGroups, err := ss.repo.GetGroups()
			if err != nil {
				return ErrGettingStateGroups
			}

			gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, stateGroups)

			if gCreate.Items == 0 {
				log.Info("no groups to be created")
			} else {
				if err := ss.scim.CreateGroups(ss.ctx, gCreate); err != nil {
					return err
				}
			}

			if gUpdate.Items == 0 {
				log.Info("no groups to be updated")
			} else {
				if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
					return err
				}
			}

			if gDelete.Items == 0 {
				log.Info("no groups to be deleted")
			} else {
				if err := ss.scim.DeleteGroups(ss.ctx, gDelete); err != nil {
					return err
				}
			}
		}

		stateUsersMetadata, err := ss.repo.GetUsersMetadata()
		if err != nil {
			return ErrGettingGroupsMetadata
		}

		if pUsersResult.HashCode == stateUsersMetadata.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")
		} else {
			log.Info("provider users and state users are diferent")

			stateUsers, err := ss.repo.GetUsers()
			if err != nil {
				return ErrGettingStateGroups
			}

			uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, stateUsers)

			if uCreate.Items == 0 {
				log.Info("no users to be created")
			} else {
				if err := ss.scim.CreateUsers(ss.ctx, uCreate); err != nil {
					return err
				}
			}

			if uUpdate.Items == 0 {
				log.Info("no users to be updated")
			} else {
				if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
					return err
				}
			}

			if uDelete.Items == 0 {
				log.Info("no users to be deleted")
			} else {
				if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
					return err
				}
			}
		}

		stateGroupsUsersMetadata, err := ss.repo.GetGroupsUsersMetadata()
		if err != nil {
			return ErrGettingGroupsMetadata
		}

		if pGroupsUsersResult.HashCode == stateGroupsUsersMetadata.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")
		} else {
			log.Info("provider groups-members and state groups-members are diferent")

			stateGroupsUsers, err := ss.repo.GetGroupsUsers()
			if err != nil {
				return ErrGettingStateGroups
			}

			ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, stateGroupsUsers)

			if ugCreate.Items == 0 {
				log.Info("no groups-members to be created")
			} else {
				if err := ss.scim.CreateMembers(ss.ctx, ugCreate); err != nil {
					return err
				}
			}

			if ugDelete.Items == 0 {
				log.Info("no groups-members to be deleted")
			} else {
				if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
					return err
				}
			}
		}
	}

	// after be sure all the SCIM part is aligned with the Identity Provider part
	// we can update the state with the identity provider data
	newState := &model.State{
		LastSync: time.Now().Format(time.RFC3339),
		Resources: model.StateResources{
			Groups:      *pGroupsResult,
			Users:       *pUsersResult,
			GroupsUsers: *pGroupsUsersResult,
		},
	}

	if err := ss.repo.UpdateState(newState); err != nil {
		return err
	}

	return nil
}

func (ss *SyncService) SyncGroupsAndUsers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	pGroups, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return err
	}

	pUsers, err := ss.prov.GetUsers(ss.ctx, ss.provUsersFilter)
	if err != nil {
		return err
	}

	if err := ss.scim.DeleteGroups(ss.ctx, pGroups); err != nil {
		return err
	}

	if err := ss.scim.DeleteUsers(ss.ctx, pUsers); err != nil {
		return err
	}

	return nil
}
