package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/core/scim_mocks.go -source=scim.go

// ScimService is the interface that needs to be implemented by the SCIM service.
type SCIMService interface {
	GetGroups(ctx context.Context) (*model.GroupsResult, error)
	GetUsers(ctx context.Context) (*model.UsersResult, error)
	GetUsersAndGroupsUsers(ctx context.Context) (*model.UsersResult, *model.GroupsUsersResult, error)

	CreateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)
	CreateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)
	CreateMembers(ctx context.Context, ur *model.GroupsUsersResult) (*model.GroupsUsersResult, error)

	UpdateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)
	UpdateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)

	DeleteGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)
	DeleteUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)
	DeleteMembers(ctx context.Context, ur *model.GroupsUsersResult) (*model.GroupsUsersResult, error)
}
