package core

import (
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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

	state, err := ss.repo.GetState()
	if err != nil {
		return ErrGettingState
	}

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

	// first time syncing
	if state.LastSync == "" {

		log.Info("State without last sync time, first time syncing")
		state.LastSync = time.Now().Format(time.RFC3339)

		// Check SCIM side to see if there are any elelemnts to be
		// reconciled. I mean SCIM is not clean until the first sync
		// and we need to reconcile the SCIM side with the identity provider side
		sGroupsResult, err := ss.scim.GetGroups(ss.ctx)
		if err != nil {
			return err
		}

		if sGroupsResult.Items == 0 {
			log.Info("No SCIM groups to be reconciliate")
		} else {
			log.Warn("Reconciliating the SCIM data with the Identity Provider data until the first syncing")

			// reconciliate elements
			gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, sGroupsResult)

			if gCreate.Items == 0 {
				log.Info("No groups to be created")
			} else {
				if err := ss.scim.CreateGroups(ss.ctx, gCreate); err != nil {
					return err
				}
			}

			if gUpdate.Items == 0 {
				log.Info("No groups to be updated")
			} else {
				if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
					return err
				}
			}

			if gDelete.Items == 0 {
				log.Info("No groups to be deleted")
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
			log.Info("No SCIM users to be reconciliate")
		} else {

			// reconciliate users
			uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, sUsersResult)

			if uCreate.Items == 0 {
				log.Info("No users to be created")
			} else {
				if err := ss.scim.CreateUsers(ss.ctx, uCreate); err != nil {
					return err
				}
			}

			if uUpdate.Items == 0 {
				log.Info("No users to be updated")
			} else {
				if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
					return err
				}
			}

			if uDelete.Items == 0 {
				log.Info("No users to be deleted")
			} else {
				if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
					return err
				}
			}
		}

		if sGroupsUsersResult.Items == 0 {
			log.Info("No SCIM groups and members to be reconciliate")
		} else {
			// reconciliate groups-users --> groups members
			ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, sGroupsUsersResult)

			if ugCreate.Items == 0 {
				log.Info("No groups-users to be created")
			} else {
				if err := ss.scim.CreateMembers(ss.ctx, ugCreate); err != nil {
					return err
				}
			}

			if ugDelete.Items == 0 {
				log.Info("No groups-users to be deleted")
			} else {
				if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
					return err
				}
			}
		}
	} else { // This is not the first time syncing
		log.Info("State with last sync time, not first time syncing")

		rGroupsResult, err := ss.repo.GetGroups()
		if err != nil {
			return err
		}

		rUsersResult, err := ss.repo.GetUsers()
		if err != nil {
			return err
		}

		rGroupsUsersResult, err := ss.repo.GetGroupsUsers()
		if err != nil {
			return err
		}

		// now here we have the google fresh data and the last sync data state
		// we need to compare the data and decide what to do
		// see differences between the two data sets
		gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, rGroupsResult)

		if gCreate.Items == 0 {
			log.Info("No groups to be created")
		} else {
			if err := ss.scim.CreateGroups(ss.ctx, gCreate); err != nil {
				return err
			}
		}

		if gUpdate.Items == 0 {
			log.Info("No groups to be updated")
		} else {
			if err := ss.scim.UpdateGroups(ss.ctx, gUpdate); err != nil {
				return err
			}
		}

		if gDelete.Items == 0 {
			log.Info("No groups to be deleted")
		} else {
			if err := ss.scim.DeleteGroups(ss.ctx, gDelete); err != nil {
				return err
			}
		}

		uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, rUsersResult)

		if uCreate.Items == 0 {
			log.Info("No users to be created")
		} else {
			if err := ss.scim.CreateUsers(ss.ctx, uCreate); err != nil {
				return err
			}
		}

		if uUpdate.Items == 0 {
			log.Info("No users to be updated")
		} else {
			if err := ss.scim.UpdateUsers(ss.ctx, uUpdate); err != nil {
				return err
			}
		}

		if uDelete.Items == 0 {
			log.Info("No users to be deleted")
		} else {
			if err := ss.scim.DeleteUsers(ss.ctx, uDelete); err != nil {
				return err
			}
		}

		ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, rGroupsUsersResult)

		if ugCreate.Items == 0 {
			log.Info("No groups-users to be created")
		} else {
			if err := ss.scim.CreateMembers(ss.ctx, ugCreate); err != nil {
				return err
			}
		}

		if ugDelete.Items == 0 {
			log.Info("No groups-users to be deleted")
		} else {
			if err := ss.scim.DeleteMembers(ss.ctx, ugDelete); err != nil {
				return err
			}
		}
	}

	// Create the sync state
	// store data to repository
	sStoreGroupsResult, err := ss.repo.StoreGroups(pGroupsResult)
	if err != nil {
		return err
	}

	sGroupsUsersResult, err := ss.repo.StoreGroupsUsers(pGroupsUsersResult)
	if err != nil {
		return err
	}

	sStoreUsersResult, err := ss.repo.StoreUsers(pUsersResult)
	if err != nil {
		return err
	}

	sStoreStateResult, err := ss.repo.StoreState(&sStoreGroupsResult, &sStoreUsersResult, &sGroupsUsersResult)
	if err != nil {
		return err
	}

	log.Infof("Sysced %s", sStoreStateResult.Location)
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
