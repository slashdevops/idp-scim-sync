package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/core/scim_mocks.go -source=scim.go

// SCIMService is the interface that needs to be implemented by the
// SCIM Provider service.
type SCIMService interface {
	GetGroups(ctx context.Context) (*model.GroupsResult, error)
	CreateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)
	UpdateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error)
	DeleteGroups(ctx context.Context, gr *model.GroupsResult) error

	GetUsers(ctx context.Context) (*model.UsersResult, error)
	CreateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)
	UpdateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error)
	DeleteUsers(ctx context.Context, ur *model.UsersResult) error

	GetUsersAndGroupsUsers(ctx context.Context) (*model.UsersResult, *model.GroupsUsersResult, error)

	CreateGroupsMembers(ctx context.Context, gur *model.GroupsUsersResult) error
	DeleteGroupsMembers(ctx context.Context, gur *model.GroupsUsersResult) error
}
