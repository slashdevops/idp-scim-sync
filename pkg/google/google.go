package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	log "github.com/sirupsen/logrus"
)

const (
	// https://cloud.google.com/storage/docs/json_api
	groupsRequiredFields    googleapi.Field = "nextPageToken, groups(id,name,email,etag)"
	membersRequiredFields   googleapi.Field = "nextPageToken, members(id,email,status,type,etag)"
	listUsersRequiredFields googleapi.Field = "nextPageToken, users(id,name,primaryEmail,suspended,etag,emails)"
	getUsersRequiredFields  googleapi.Field = "id,name,primaryEmail,suspended,etag"
)

var (
	// ErrGoogleClientScopeNil is returned when the scope is nil.
	ErrGoogleClientScopeNil = fmt.Errorf("google: google client scope is required")

	// ErrUserIDNil is returned when the user ID is nil.
	ErrUserIDNil = fmt.Errorf("google: user id is required")

	// ErrGroupIDNil is returned when the group ID is nil.
	ErrGroupIDNil = fmt.Errorf("google: group id is required")
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

	creds, err := google.CredentialsFromJSONWithParams(ctx, serviceAccount, google.CredentialsParams{
		Scopes:  scope,
		Subject: userEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("google: error getting config for Service Account: %v", err)
	}

	svc, err := admin.NewService(ctx, option.WithTokenSource(creds.TokenSource))
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
				err = ds.svc.Users.List().Query(q).Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
			} else {
				err = ds.svc.Users.List().Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
			}
			if err != nil {
				return nil, err
			}
		}
	} else {
		err = ds.svc.Users.List().Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
		if err != nil {
			return nil, err
		}
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
			if err != nil {
				return nil, err
			}
		}
	} else {
		err = ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return g, err
}

// ListGroupMembers return a list of all members given a group ID.
// references:
// - https://developers.google.com/admin-sdk/directory/reference/rest/v1/members/list
// - https://developers.google.com/admin-sdk/directory/v1/guides/manage-group-members
// - https://cloud.google.com/identity/docs/how-to/query-memberships
func (ds *DirectoryService) ListGroupMembers(ctx context.Context, groupID string, queries ...GetGroupMembersOption) ([]*admin.Member, error) {
	if groupID == "" {
		return nil, ErrGroupIDNil
	}

	qs := getGroupMembersOptions{}
	for _, q := range queries {
		q(&qs)
	}

	m := make([]*admin.Member, 0)
	mlc := ds.svc.Members.List(groupID)

	if qs.includeDerivedMembership {
		mlc = mlc.IncludeDerivedMembership(true)
	}
	if qs.maxResults > 0 {
		mlc = mlc.MaxResults(qs.maxResults)
	}
	if qs.pageToken != "" {
		mlc = mlc.PageToken(qs.pageToken)
	}
	if qs.roles != "" {
		mlc = mlc.Roles(qs.roles)
	}

	err := mlc.Fields(membersRequiredFields).Pages(ctx, func(members *admin.Members) error {
		for _, member := range members.Members {
			// Add only active members to list
			if member.Status == "ACTIVE" {
				m = append(m, member)
			} else {
				log.Warnf("google: not including %s to group %s members due to incorrect status %s (not ACTIVE)", member.Email, groupID, member.Status)
			}
		}
		return nil
	})

	return m, err
}

// GetUser return a user given a user ID.
// userID: the user's primary email address, alias email address, or unique user ID.
func (ds *DirectoryService) GetUser(ctx context.Context, userID string) (*admin.User, error) {
	if userID == "" {
		return nil, ErrUserIDNil
	}

	u, err := ds.svc.Users.Get(userID).Fields(getUsersRequiredFields).Context(ctx).Do()
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
