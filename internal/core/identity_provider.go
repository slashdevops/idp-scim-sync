package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate mockgen -package=mocks -destination=../mocks/core/identity_provider_mocks.go -source=identity_provider.go

// IdentityProviderService is the interface that needs to be implemented by the
// Identity Provider service.
type IdentityProviderService interface {
	GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error)
	GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error)
	GetUsersFromGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error)
	GetUsersAndGroupsUsers(ctx context.Context, groups *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error)
}
