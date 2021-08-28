package core

import (
	"context"
	"errors"
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/provider"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/scim"
)

var (
	ErrNilContext         = errors.New("context cannot be nil")
	ErrProviderServiceNil = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
	ErrGettingGroups      = errors.New("error getting groups")
	ErrRepositoryNil      = errors.New("repository cannot be nil")
)

type SyncService interface {
	SyncGroupsAndTheirMembers() error
	SyncGroupsAndUsers() error
}

type syncService struct {
	ctx              context.Context
	mu               *sync.RWMutex
	prov             provider.IdentityProviderService
	provGroupsFilter []string
	provUsersFilter  []string
	scim             scim.SCIMService
	repo             repository.SyncRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov provider.IdentityProviderService, scim scim.SCIMService, repo repository.SyncRepository, opts ...SyncServiceOption) (SyncService, error) {
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

	ss := &syncService{
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

func (ss *syncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Check if is the first time we are syncing
	state, err := ss.repo.GetState()
	if err != nil {
		return err
	}

	// retrive data from provider
	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return ErrGettingGroups
	}

	pUsers := make([]*model.User, 0)
	var pGroupsMembers model.GroupsMembers

	for _, pGroup := range pGroupsResult.Resources {

		pMembers, err := ss.prov.GetGroupMembers(ss.ctx, pGroup.ID)
		if err != nil {
			return err
		}

		pGroupsMembers[pGroup.ID] = pMembers.Resources

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
		Resources: &pGroupsMembers,
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
		_ = sGroupsResults
		_ = sUsersResult
		_ = sGroupsMembersResult
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

func (ss *syncService) SyncGroupsAndUsers() error {
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