package scim

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

type SCIMService interface {
	GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error)
	CreateOrUpdateGroups(ctx context.Context, gr *model.GroupsResult) error
	CreateOrUpdateUsers(ctx context.Context, ur *model.UsersResult) error
	DeleteGroups(ctx context.Context, gr *model.GroupsResult) error
	DeleteUsers(ctx context.Context, ur *model.UsersResult) error
}
