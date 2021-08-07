package sync

import "context"

type IdentityProviderService interface {
	GetGroups(ctx context.Context, filter []string) (*GroupsResult, error)
	GetUsers(ctx context.Context, filter []string) (*UsersResult, error)
	GetGroupMembers(ctx context.Context, groupID string) (*MembersResult, error)
	GetUsersFromGroupMembers(ctx context.Context, members *MembersResult) (*UsersResult, error)
}
