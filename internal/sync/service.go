package sync

import (
	"context"
	"errors"
	"sync"

	"github.com/slashdevops/aws-sso-gws-sync/internal/config"
)

var (
	ErrNilContext         = errors.New("context cannot be nil")
	ErrProviderServiceNil = errors.New("provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
)

type ProviderService interface {
	GetGoups(*context.Context, []string) (*GroupResult, error)
	GetUsers(*context.Context, []string) (*UserResult, error)
	GetGroupsMembers(*context.Context, *GroupResult) (*MemberResult, error)
	GetUsersFromGroupsMembers(*context.Context, []string, *MemberResult) (*UserResult, error)
}

type SCIMService interface {
	GetGroups(*context.Context, []string) (*GroupResult, error)
	GetUsers(*context.Context, []string) (*UserResult, error)
	CreateOrUpdateGroups(*context.Context, *GroupResult) error
	CreateOrUpdateUsers(*context.Context, *UserResult) error
	DeleteGroups(*context.Context, *GroupResult) error
	DeleteUsers(*context.Context, *UserResult) error
}

type SyncServiceOption func(*syncService)

func WithProviderGroupsFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provGroupsFilter = filter
	}
}

func WithProviderUsersFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provUsersFilter = filter
	}
}

type SyncService interface {
	SyncGroupsAndTheirMembers() error
	SyncGroupsAndUsers() error
}

type syncService struct {
	ctx              *context.Context
	config           config.Config
	mu               *sync.Mutex
	prov             ProviderService
	provGroupsFilter []string
	provUsersFilter  []string
	scim             SCIMService
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx *context.Context, prov ProviderService, scim SCIMService, opts ...SyncServiceOption) (SyncService, error) {

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

	pGroups, err := ss.prov.GetGoups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return err
	}

	pGroupsMembers, err := ss.prov.GetGroupsMembers(ss.ctx, pGroups)
	if err != nil {
		return err
	}

	pUsers, err := ss.prov.GetUsersFromGroupsMembers(ss.ctx, ss.provUsersFilter, pGroupsMembers)
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

func (ss *syncService) SyncGroupsAndUsers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	pGroups, err := ss.prov.GetGoups(ss.ctx, ss.provGroupsFilter)
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
