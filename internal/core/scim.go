package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/core/scim_mocks.go -source=scim.go

// SCIMService is the interface that needs to be implemented by the
// SCIM Provider service.
type SCIMService interface {
	// GetGroups returns a list of all groups from the SCIM service.
	GetGroups(ctx context.Context) (*model.GroupsResult, error)

	// CreateGroups create groups in the SCIM Service given a list of groups.
	CreateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)

	// UpdateGroups updates groups in the SCIM Service given a list of groups.
	UpdateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)

	// DeleteGroups deletes groups in the SCIM Service given a list of groups.
	DeleteGroups(ctx context.Context, gr *model.GroupsResult) error

	// GetUsers returns a list of all users from the SCIM service.
	GetUsers(ctx context.Context) (*model.UsersResult, error)

	// CreateUsers create users in the SCIM Service given a list of users.
	CreateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)

	// UpdateUsers updates users in the SCIM Service given a list of users.
	UpdateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)

	// DeleteUsers deletes users in the SCIM Service given a list of users.
	DeleteUsers(ctx context.Context, ur *model.UsersResult) error

	// GetGroupsMembers the Groups and their Members from the SCIM service.
	GetGroupsMembers(ctx context.Context, gr *model.GroupsResult) (*model.GroupsMembersResult, error)

	// CreateGroupsMembers create groups members in the SCIM Service given a list of groups members.
	CreateGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) (*model.GroupsMembersResult, error)

	// DeleteGroupsMembers deletes groups members in the SCIM Service given a list of groups members.
	DeleteGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) error
}
