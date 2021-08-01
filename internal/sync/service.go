package sync

import (
	"context"
	"errors"
	"sync"

	"github.com/slashdevops/aws-sso-gws-sync/internal/config"
)

var (
	ErrNilContext         = errors.New("context cannot be nil")
	ErrProviderServiceNil = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
)

type IdentityProviderService interface {
	GetGroups(ctx context.Context, filter []string) (*GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*UsersResult, error)
	GetGroupMembers(ctx context.Context, groupID string) (*MembersResult, error)
	GetUsersFromGroupMembers(ctx context.Context, members *MembersResult) (*UsersResult, error)
}

type SCIMService interface {
	GetGroups(ctx context.Context, filter []string) (*GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*UsersResult, error)
	CreateOrUpdateGroups(ctx context.Context, gr *GroupsResult) error
	CreateOrUpdateUsers(ctx context.Context, ur *UsersResult) error
	DeleteGroups(ctx context.Context, gr *GroupsResult) error
	DeleteUsers(ctx context.Context, ur *UsersResult) error
}

type SyncServiceOption func(*syncService)

func WithIdentityProviderGroupsFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provGroupsFilter = filter
	}
}

func WithIdentityProviderUsersFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provUsersFilter = filter
	}
}

type SyncService interface {
	SyncGroupsAndTheirMembers() error
	SyncGroupsAndUsers() error
}

type syncService struct {
	ctx              context.Context
	config           config.Config
	mu               *sync.Mutex
	prov             IdentityProviderService
	provGroupsFilter []string
	provUsersFilter  []string
	scim             SCIMService
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, opts ...SyncServiceOption) (SyncService, error) {

	if ctx == nil {
		return nil, ErrNilContext
	}
	if prov == nil {
		return nil, ErrProviderServiceNil
	}
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}

	return &syncService{
		ctx:  ctx,
		mu:   &sync.Mutex{},
		prov: prov,
		scim: scim,
	}, nil
}

func (ss *syncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	pGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return err
	}

	pUsers := make([]*User, 0)
	for _, pGroup := range pGroupsResult.Resources {

		pGroupMembers, err := ss.prov.GetGroupMembers(ss.ctx, pGroup.Id.IdentityProvider)
		if err != nil {
			return err
		}

		pUsersFromMembers, err := ss.prov.GetUsersFromGroupMembers(ss.ctx, pGroupMembers)
		if err != nil {
			return err
		}

		pUsers = append(pUsers, pUsersFromMembers.Resources...)
	}

	pUsersResult := &UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}

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
