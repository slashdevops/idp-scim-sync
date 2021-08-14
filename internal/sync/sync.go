package sync

import (
	"context"
	"errors"
	"sync"
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
	prov             IdentityProviderService
	provGroupsFilter []string
	provUsersFilter  []string
	scim             SCIMService
	repo             SyncRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, repo SyncRepository, opts ...SyncServiceOption) (SyncService, error) {

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

	// retrive data from provider
	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return ErrGettingGroups
	}

	pUsers := make([]*User, 0)
	var pGroupsMembers GroupsMembers

	for _, pGroup := range pGroupsResult.Resources {

		pMembers, err := ss.prov.GetGroupMembers(ss.ctx, pGroup.Id)
		if err != nil {
			return err
		}

		pGroupsMembers[pGroup.Id] = pMembers.Resources

		pUsersFromMembers, err := ss.prov.GetUsersFromGroupMembers(ss.ctx, pMembers)
		if err != nil {
			return err
		}

		pUsers = append(pUsers, pUsersFromMembers.Resources...)
	}

	pUsersResult := &UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}

	pGroupsMembersResult := &GroupsMembersResult{
		Items:     len(pGroupsMembers),
		Resources: &pGroupsMembers,
	}

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

	state, err := CreateSyncState(&sStoreGroupsResult, &sStoreGroupsMembersResult, &sStoreUsersResult)
	if err != nil {
		return err
	}

	sStoreStateResult, err := ss.repo.StoreState(&state)
	if err != nil {
		return err
	}
	_ = sStoreStateResult

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

func CreateSyncState(sgr *StoreGroupsResult, sgmr *StoreGroupsMembersResult, sur *StoreUsersResult) (SyncState, error) {
	return SyncState{
		Version:  "1.0.0",
		Checksum: "TBD",
	}, nil
}
