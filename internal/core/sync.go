package core

import (
	"context"
	"errors"
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrNilContext         = errors.New("context cannot be nil")
	ErrProviderServiceNil = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
	ErrGettingGroups      = errors.New("error getting groups")
	ErrRepositoryNil      = errors.New("repository cannot be nil")
	ErrGettingState       = errors.New("error getting state")
)

type SyncService struct {
	ctx              context.Context
	mu               *sync.RWMutex
	prov             IdentityProviderService
	provGroupsFilter []string
	provUsersFilter  []string
	scim             SCIMService
	repo             SyncRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, repo SyncRepository, opts ...SyncServiceOption) (*SyncService, error) {
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

	// Check if is the first time we are syncing
	state, err := ss.repo.GetState()
	if err != nil {
		return ErrGettingState
	}

	// retrive data from provider
	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return ErrGettingGroups
	}

	pUsers := make([]*model.User, 0)
	pGroupsMembers := make([]*model.GroupMembers, 0)

	for _, pGroup := range pGroupsResult.Resources {

		pMembers, err := ss.prov.GetGroupMembers(ss.ctx, pGroup.ID)
		if err != nil {
			return err
		}

		pGroupMembers := &model.GroupMembers{
			ID:        pGroup.ID,
			Email:     pGroup.Email,
			Items:     len(pMembers.Resources),
			Resources: pMembers.Resources,
		}
		pGroupsMembers = append(pGroupsMembers, pGroupMembers)

		pUsersFromMembers, err := ss.prov.GetUsersFromGroupMembers(ss.ctx, pMembers)
		if err != nil {
			return err
		}

		pUsers = append(pUsers, pUsersFromMembers.Resources...)
	}

	pUsersResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}

	pGroupsMembersResult := &model.GroupsMembersResult{
		Items:     len(pGroupsMembers),
		Resources: pGroupsMembers,
	}

	// First time we are syncing
	if state.Checksum == "" {

		// Create the sync state
		// store data to repository
		sStoreGroupsResult, err := ss.repo.StoreGroups(pGroupsResult)
		if err != nil {
			return err
		}

		sStoreGroupsMembersResult, err := ss.repo.StoreGroupsMembers(pGroupsMembersResult)
		if err != nil {
			return err
		}

		sStoreUsersResult, err := ss.repo.StoreUsers(pUsersResult)
		if err != nil {
			return err
		}

		// reusing the state variable
		state, err = createSyncState(&sStoreGroupsResult, &sStoreGroupsMembersResult, &sStoreUsersResult)
		if err != nil {
			return err
		}

		sStoreStateResult, err := ss.repo.StoreState(&state)
		if err != nil {
			return err
		}

		// TODO: decide what to do with the result
		_ = sStoreStateResult

	} else {
		sGroupsResults, err := ss.repo.GetGroups(state.Groups.Place)
		if err != nil {
			return err
		}

		sUsersResult, err := ss.repo.GetUsers(state.Users.Place)
		if err != nil {
			return err
		}

		sGroupsMembersResult, err := ss.repo.GetGroupsMembers(state.GroupsMembers.Place)
		if err != nil {
			return err
		}

		// now here we have the google fresh data and the last sync data state
		// we need to compare the data and decide what to do
		// see differences between the two data sets
		_, _, _, _ = groupsDifferences(pGroupsResult, sGroupsResults)
		_, _, _, _ = usersDifferences(pUsersResult, sUsersResult)
		_, _, _ = groupsMembersDifferences(pGroupsMembersResult, sGroupsMembersResult)
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
