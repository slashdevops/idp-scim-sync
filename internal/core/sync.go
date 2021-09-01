package core

import (
	"context"
	"errors"
	"sync"

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
	repo             SyncRepository
	state            SyncState
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, repo SyncRepository, state SyncState, opts ...SyncServiceOption) (*SyncService, error) {
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
	if state == nil {
		return nil, ErrSyncStateNil
	}

	ss := &SyncService{
		ctx:              ctx,
		mu:               &sync.RWMutex{},
		prov:             prov,
		provGroupsFilter: []string{}, // fill in with the opts
		provUsersFilter:  []string{}, // fill in with the opts
		scim:             scim,
		repo:             repo,
		state:            state,
	}

	for _, opt := range opts {
		opt(ss)
	}

	return ss, nil
}

func (ss *SyncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// retrive data from provider
	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return ErrGettingGroups
	}

	if pGroupsResult.Items == 0 {
		log.Info("No groups to sync")
		return nil
	}

	pUsersResult, pGroupsUsersResult, err := ss.prov.GetUsersAndGroupsUsers(ss.ctx, pGroupsResult)
	if err != nil {
		return ErrGettingGetUsersAndGroupsUsers
	}

	if pUsersResult.Items == 0 {
		log.Info("No users to sync")
	}

	state, err := ss.repo.GetState(ss.state.GetName())
	if err != nil {
		return ErrGettingState
	}

	// first time syncing, no state hashcode
	if state.HashCode == "" {

		// Check SCIM side to see if there are any elelemnts to be
		// reconciled.
		sGroupsResult, err := ss.scim.GetGroups(ss.ctx, ss.provGroupsFilter)
		if err != nil {
			return err
		}

		if sGroupsResult.Items == 0 {
			log.Info("No groups to reconciliate")
		}

		// Create the sync state
		// store data to repository
		sStoreGroupsResult, err := ss.repo.StoreGroups(pGroupsResult)
		if err != nil {
			return err
		}

		sStoreGroupsUsersResult, err := ss.repo.StoreGroupsUsers(pGroupsUsersResult)
		if err != nil {
			return err
		}

		sStoreUsersResult, err := ss.repo.StoreUsers(pUsersResult)
		if err != nil {
			return err
		}

		_ = sStoreGroupsResult
		_ = sStoreGroupsUsersResult
		_ = sStoreUsersResult

		// // reusing the state variable
		// state, err = createSyncState(&sStoreGroupsResult, &sStoreGroupsUsersResult, &sStoreUsersResult)
		// if err != nil {
		// 	return err
		// }

		// sStoreStateResult, err := ss.repo.StoreState(&state)
		// if err != nil {
		// 	return err
		// }

		// // TODO: decide what to do with the result
		// _ = sStoreStateResult

	} else {
		sGroupsResults, err := ss.repo.GetGroups(state.Groups.Place)
		if err != nil {
			return err
		}

		sUsersResult, err := ss.repo.GetUsers(state.Users.Place)
		if err != nil {
			return err
		}

		sGroupsUsersResult, err := ss.repo.GetGroupsUsers(state.GroupsMembers.Place)
		if err != nil {
			return err
		}

		// now here we have the google fresh data and the last sync data state
		// we need to compare the data and decide what to do
		// see differences between the two data sets
		_, _, _, _ = groupsOperations(pGroupsResult, sGroupsResults)
		_, _, _, _ = usersOperations(pUsersResult, sUsersResult)
		_, _, _ = groupsUsersOperations(pGroupsUsersResult, sGroupsUsersResult)
	}

	// sync data to SCIM
	if err := ss.scim.CreateOrUpdateUsers(ss.ctx, pUsersResult); err != nil {
		return err
	}

	if err := ss.scim.CreateOrUpdateGroups(ss.ctx, pGroupsResult); err != nil {
		return err
	}

	if err := ss.scim.DeleteUsers(ss.ctx, pUsersResult); err != nil {
		return err
	}

	if err := ss.scim.DeleteGroups(ss.ctx, pGroupsResult); err != nil {
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

	if err := ss.scim.CreateOrUpdateGroups(ss.ctx, pGroups); err != nil {
		return err
	}

	if err := ss.scim.CreateOrUpdateUsers(ss.ctx, pUsers); err != nil {
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
