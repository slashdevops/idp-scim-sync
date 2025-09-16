package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go tool mockgen -package=mocks -destination=../../mocks/core/idp_mocks.go -source=idp.go

// IdentityProviderService is the interface consumed by the core services and
// needs to be implemented by the Identity Provider service.
type IdentityProviderService interface {
	// GetGroups returns the groups filtered by the given filter in the Identity provider side.
	GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error)

	// GetUsers returns the users filtered by the given filter in the Identity provider side.
	GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error)

	// GetGroupMembers returns the members of the given group in the Identity provider side.
	GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error)

	// GetUsersByGroupsMembers returns the users belongs to the given group in the Identity provider side.
	GetUsersByGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) (*model.UsersResult, error)

	// GetGroupsMembers returns the groups and their members from the Identity provider side.
	GetGroupsMembers(ctx context.Context, gr *model.GroupsResult) (*model.GroupsMembersResult, error)
}
