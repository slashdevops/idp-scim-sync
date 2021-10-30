package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/core/idp_mocks.go -source=idp.go

// IdentityProviderService is the interface that needs to be implemented by the
// Identity Provider service.
type IdentityProviderService interface {
	// GetGroups returns the groups filtered by the given filter in the Identity provider side.
	GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error)

	// GetUsers returns the users filtered by the given filter in the Identity provider side.
	GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error)

	// GetGroupMembers returns the members of the given group in the Identity provider side.
	GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error)

	// GetUsersByGroupMembers returns the users belongs to the given group in the Identity provider side.
	GetUsersByGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error)

	// GetUsersAndGroupsUsers returns the users and the groups and their users of the given groups in the Identity provider side.
	GetUsersAndGroupsUsers(ctx context.Context, gr *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error)
}
