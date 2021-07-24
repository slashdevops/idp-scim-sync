package google

import (
	"context"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type ClientService interface {
	ListGroups([]string) ([]*admin.Group, error)
}

type Client struct {
	ctx     context.Context
	service *admin.Service
}

// NewClient create a Google Directory API client.
// References:
// - https://developers.google.com/admin-sdk/directory/v1/guides/delegation?utm_source=pocket_mylist#go
func NewClient(ctx context.Context, AdminEmail string, ServiceAccount []byte) (*Client, error) {

	config, err := google.JWTConfigFromJSON(
		ServiceAccount,
		admin.AdminDirectoryGroupReadonlyScope,
		admin.AdminDirectoryGroupMemberReadonlyScope,
		admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, err
	}

	config.Subject = AdminEmail

	ts := config.TokenSource(ctx)

	srv, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:     ctx,
		service: srv,
	}, nil
}

// ListUsers list all users in a Google Directory filtered by query.
func (c *Client) ListUsers(query []string) ([]*admin.User, error) {
	u := make([]*admin.User, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = c.service.Users.List().Query(q).Customer("my_customer").Pages(c.ctx, func(users *admin.Users) error {
				u = append(u, users.Users...)
				return nil
			})
		}
	} else {
		err = c.service.Users.List().Customer("my_customer").Pages(c.ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
	}
	return u, err
}

// ListGroups list all groups in a Google Directory filtered by query.
func (c *Client) ListGroups(query []string) ([]*admin.Group, error) {
	g := make([]*admin.Group, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			err = c.service.Groups.List().Customer("my_customer").Query(q).Pages(c.ctx, func(groups *admin.Groups) error {
				g = append(g, groups.Groups...)
				return nil
			})
		}
	} else {
		err = c.service.Groups.List().Customer("my_customer").Pages(c.ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})

	}
	return g, err
}
