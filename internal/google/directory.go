package google

import (
	"context"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// https://cloud.google.com/storage/docs/json_api
	groupsRequiredFields  googleapi.Field = "id,name,email"
	membersRequiredFields googleapi.Field = "id,email"
	usersRequiredFields   googleapi.Field = "id,name,primaryEmail,suspended"
)

// DirectoryService is a facade of Google Directory API client.
type DirectoryService interface {
	ListUsers(query []string) ([]*admin.User, error)
	ListGroups(query []string) ([]*admin.Group, error)
	ListGroupMembers(groupID string) ([]*admin.Member, error)
	GetUser(userID string) (*admin.User, error)
}

type directoryService struct {
	ctx context.Context
	svc *admin.Service
}

// NewService create a Google Directory Service.
// References:
// - https://pkg.go.dev/google.golang.org/api/admin/directory/v1
// Examples of scope:
// - admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope, admin.AdminDirectoryUserReadonlyScope
func NewService(ctx context.Context, UserEmail string, ServiceAccount []byte, scope ...string) (*admin.Service, error) {

	config, err := google.JWTConfigFromJSON(ServiceAccount, scope...)
	if err != nil {
		return nil, err
	}

	config.Subject = UserEmail

	ts := config.TokenSource(ctx)

	svc, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// NewDirectoryService create a Google Directory API client.
// References:
// - https://developers.google.com/admin-sdk/directory/v1/guides/delegation?utm_source=pocket_mylist#go
func NewDirectoryService(ctx context.Context, svc *admin.Service) (DirectoryService, error) {
	return &directoryService{
		ctx: ctx,
		svc: svc,
	}, nil
}

// ListUsers list all users in a Google Directory filtered by query.
func (d *directoryService) ListUsers(query []string) ([]*admin.User, error) {
	u := make([]*admin.User, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = d.svc.Users.List().Query(q).Customer("my_customer").Fields(usersRequiredFields).Pages(d.ctx, func(users *admin.Users) error {
				u = append(u, users.Users...)
				return nil
			})
		}
	} else {
		err = d.svc.Users.List().Customer("my_customer").Fields(usersRequiredFields).Pages(d.ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
	}
	return u, err
}

// ListGroups list all groups in a Google Directory filtered by query.
// References:
// - https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups
func (d *directoryService) ListGroups(query []string) ([]*admin.Group, error) {
	g := make([]*admin.Group, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = d.svc.Groups.List().Customer("my_customer").Query(q).Fields(groupsRequiredFields).Pages(d.ctx, func(groups *admin.Groups) error {
				g = append(g, groups.Groups...)
				return nil
			})
		}
	} else {
		err = d.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(d.ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})

	}
	return g, err
}

func (d *directoryService) ListGroupMembers(groupID string) ([]*admin.Member, error) {
	m := make([]*admin.Member, 0)

	err := d.svc.Members.List(groupID).Fields(membersRequiredFields).Pages(d.ctx, func(members *admin.Members) error {
		m = append(m, members.Members...)
		return nil
	})

	return m, err
}

func (d *directoryService) GetUser(userID string) (*admin.User, error) {

	u, err := d.svc.Users.Get(userID).Fields(usersRequiredFields).Context(d.ctx).Do()

	return u, err
}
