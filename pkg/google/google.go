package google

import (
	"context"
	"errors"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// https://cloud.google.com/storage/docs/json_api
	groupsRequiredFields  googleapi.Field = "groups(id,name,email)"
	membersRequiredFields googleapi.Field = "members(id,email)"
	usersRequiredFields   googleapi.Field = "users(id,name,primaryEmail,suspended)"
)

var (
	ErrInvalidServiceAccount    = errors.New("invalid service account")
	ErrCreatingDirectoryService = errors.New("creating directory service")
)

type DirectoryService struct {
	ctx context.Context
	svc *admin.Service
}

// NewService create a Google Directory Service.
// References:
// - https://pkg.go.dev/google.golang.org/api/admin/directory/v1
// Examples of scope:
// - "https://www.googleapis.com/auth/admin.directory.group.readonly"
// - "https://www.googleapis.com/auth/admin.directory.group.member.readonly"
// - "https://www.googleapis.com/auth/admin.directory.user.readonly"
func NewService(ctx context.Context, UserEmail string, ServiceAccount []byte, scope ...string) (*admin.Service, error) {
	config, err := google.JWTConfigFromJSON(ServiceAccount, scope...)
	if err != nil {
		return nil, ErrInvalidServiceAccount
	}

	config.Subject = UserEmail
	ts := config.TokenSource(ctx)

	svc, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, ErrCreatingDirectoryService
	}

	return svc, nil
}

// NewDirectoryService create a Google Directory API client.
// References:
// - https://developers.google.com/admin-sdk/directory/v1/guides/delegation?utm_source=pocket_mylist#go
func NewDirectoryService(ctx context.Context, svc *admin.Service) (*DirectoryService, error) {
	return &DirectoryService{
		ctx: ctx,
		svc: svc,
	}, nil
}

// ListUsers list all users in a Google Directory filtered by query.
func (ds *DirectoryService) ListUsers(query []string) ([]*admin.User, error) {
	u := make([]*admin.User, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = ds.svc.Users.List().Query(q).Customer("my_customer").Fields(usersRequiredFields).Pages(ds.ctx, func(users *admin.Users) error {
				u = append(u, users.Users...)
				return nil
			})
		}
	} else {
		err = ds.svc.Users.List().Customer("my_customer").Fields(usersRequiredFields).Pages(ds.ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
	}
	return u, err
}

// ListGroups list all groups in a Google Directory filtered by query.
// References:
// - https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups
func (ds *DirectoryService) ListGroups(query []string) ([]*admin.Group, error) {
	g := make([]*admin.Group, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = ds.svc.Groups.List().Customer("my_customer").Query(q).Fields(groupsRequiredFields).Pages(ds.ctx, func(groups *admin.Groups) error {
				g = append(g, groups.Groups...)
				return nil
			})
		}
	} else {
		err = ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ds.ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})
	}
	return g, err
}

// ListGroupMembers list all members in a Google Directory group filtered by query.
func (ds *DirectoryService) ListGroupMembers(groupID string) ([]*admin.Member, error) {
	m := make([]*admin.Member, 0)

	err := ds.svc.Members.List(groupID).Fields(membersRequiredFields).Pages(ds.ctx, func(members *admin.Members) error {
		m = append(m, members.Members...)
		return nil
	})

	return m, err
}

// GetUser get a user in a Google Directory filtered by query.
func (ds *DirectoryService) GetUser(userID string) (*admin.User, error) {
	u, err := ds.svc.Users.Get(userID).Fields(usersRequiredFields).Context(ds.ctx).Do()

	return u, err
}

// GetGroups get a group in a Google Directory filtered by query.
func (ds *DirectoryService) GetGroup(groupID string) (*admin.Group, error) {
	g, err := ds.svc.Groups.Get(groupID).Fields(groupsRequiredFields).Context(ds.ctx).Do()

	return g, err
}
