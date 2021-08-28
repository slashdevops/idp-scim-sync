package core

import "context"

type SCIMService interface {
	GetGroups(ctx context.Context, filter []string) (*GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*UsersResult, error)
	CreateOrUpdateGroups(ctx context.Context, gr *GroupsResult) error
	CreateOrUpdateUsers(ctx context.Context, ur *UsersResult) error
	DeleteGroups(ctx context.Context, gr *GroupsResult) error
	DeleteUsers(ctx context.Context, ur *UsersResult) error
}
