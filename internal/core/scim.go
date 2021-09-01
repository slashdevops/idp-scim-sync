package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// ScimService is the interface that needs to be implemented by the SCIM service.
type SCIMService interface {
	GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error)
	GetUsersAndGroupsUsers(ctx context.Context, groups *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error)

	CreateOrUpdateGroups(ctx context.Context, gr *model.GroupsResult) error
	CreateOrUpdateUsers(ctx context.Context, ur *model.UsersResult) error
	DeleteGroups(ctx context.Context, gr *model.GroupsResult) error
	DeleteUsers(ctx context.Context, ur *model.UsersResult) error
}
