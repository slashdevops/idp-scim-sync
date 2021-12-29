package google

import (
	"context"
	"fmt"

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
	ErrGoogleClientScopeNil = fmt.Errorf("google: google client scope is required")
	ErrUserIDNil            = fmt.Errorf("google: user id is required")
	ErrGroupIDNil           = fmt.Errorf("google: group id is required")
)

// DirectoryService represent the  Google Directory API client.
type DirectoryService struct {
	svc *admin.Service
}

// NewService create a Google Directory Service.
// References:
// - https://pkg.go.dev/google.golang.org/api/admin/directory/v1
// Examples of scope:
// - "https://www.googleapis.com/auth/admin.directory.group.readonly"
// - "https://www.googleapis.com/auth/admin.directory.group.member.readonly"
// - "https://www.googleapis.com/auth/admin.directory.user.readonly"
func NewService(ctx context.Context, userEmail string, serviceAccount []byte, scope ...string) (*admin.Service, error) {
	if len(scope) == 0 {
		return nil, ErrGoogleClientScopeNil
	}

	config, err := google.JWTConfigFromJSON(serviceAccount, scope...)
	if err != nil {
		return nil, fmt.Errorf("google: error getting JWT config from Service Account: %v", err)
	}

	config.Subject = userEmail
	ts := config.TokenSource(ctx)

	svc, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("google: error creating service: %v", err)
	}

	return svc, nil
}

// NewDirectoryService create a Google Directory API client.
// References:
// - https://developers.google.com/admin-sdk/directory/v1/guides/delegation?utm_source=pocket_mylist#go
func NewDirectoryService(svc *admin.Service) (*DirectoryService, error) {
	return &DirectoryService{
		svc: svc,
	}, nil
}

// ListUsers list all users in a Google Directory filtered by query.
func (ds *DirectoryService) ListUsers(ctx context.Context, query []string) ([]*admin.User, error) {
	u := make([]*admin.User, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			if q != "" {
				err = ds.svc.Users.List().Query(q).Customer("my_customer").Fields(usersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
			} else {
				err = ds.svc.Users.List().Customer("my_customer").Fields(usersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
			}
		}
	} else {
		err = ds.svc.Users.List().Customer("my_customer").Fields(usersRequiredFields).Pages(ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
	}
	return u, err
}

// ListGroups list all groups in a Google Directory filtered by query.
// References:
// - https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups
func (ds *DirectoryService) ListGroups(ctx context.Context, query []string) ([]*admin.Group, error) {
	g := make([]*admin.Group, 0)
	var err error

	if len(query) > 0 {
		for _, q := range query {
			if q != "" {
				err = ds.svc.Groups.List().Customer("my_customer").Query(q).Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
					g = append(g, groups.Groups...)
					return nil
				})
			} else {
				err = ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
					g = append(g, groups.Groups...)
					return nil
				})
			}
		}
	} else {
		err = ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})
	}
	return g, err
}

// ListGroupMembers return a list of all members given a group ID.
func (ds *DirectoryService) ListGroupMembers(ctx context.Context, groupID string) ([]*admin.Member, error) {
	if groupID == "" {
		return nil, ErrGroupIDNil
	}

	m := make([]*admin.Member, 0)

	err := ds.svc.Members.List(groupID).Fields(membersRequiredFields).Pages(ctx, func(members *admin.Members) error {
		m = append(m, members.Members...)
		return nil
	})

	return m, err
}

// GetUser return a user given a user ID.
func (ds *DirectoryService) GetUser(ctx context.Context, userID string) (*admin.User, error) {
	if userID == "" {
		return nil, ErrUserIDNil
	}

	// TODO: u, err := ds.svc.Users.Get(userID).Fields(usersRequiredFields).Context(ctx).Do()
	u, err := ds.svc.Users.Get(userID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("google: error getting user %s: %v", userID, err)
	}

	return u, nil
}

// GetGroup return a group given a group ID.
func (ds *DirectoryService) GetGroup(ctx context.Context, groupID string) (*admin.Group, error) {
	if groupID == "" {
		return nil, ErrGroupIDNil
	}

	g, err := ds.svc.Groups.Get(groupID).Fields(groupsRequiredFields).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("google: error getting group %s: %v", groupID, err)
	}

	return g, nil
}
