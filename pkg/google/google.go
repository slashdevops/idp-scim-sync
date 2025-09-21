package google

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// Base fields common to all objects
	baseFields = "id,etag"

	// Field definitions for specific object types
	userFields   = baseFields + ",primaryEmail,name,suspended,kind,emails,addresses,organizations,phones,languages,locations"
	groupFields  = baseFields + ",name,email"
	memberFields = baseFields + ",email,status,type"

	// Complete field specifications for API calls
	groupsRequiredFields    googleapi.Field = "nextPageToken, groups(" + groupFields + ")"
	membersRequiredFields   googleapi.Field = "nextPageToken, members(" + memberFields + ")"
	listUsersRequiredFields googleapi.Field = "nextPageToken, users(" + userFields + ")"
	getUsersRequiredFields  googleapi.Field = userFields
)

var (
	// ErrGoogleClientScopeNil is returned when the scope is nil.
	ErrGoogleClientScopeNil = fmt.Errorf("google: google client scope is required")

	// ErrUserIDNil is returned when the user ID is nil.
	ErrUserIDNil = fmt.Errorf("google: user id is required")

	// ErrUserEmailNil is returned when the user email is nil.
	ErrUserEmailNil = fmt.Errorf("google: user email is required")

	// ErrGroupIDNil is returned when the group ID is nil.
	ErrGroupIDNil = fmt.Errorf("google: group id is required")

	// ErrServiceAccountNil is returned when the service account credentials are nil.
	ErrServiceAccountNil = fmt.Errorf("google: service account credentials are required")

	// ErrUserAgentNil is returned when the user agent is nil.
	ErrUserAgentNil = fmt.Errorf("google: user agent is required")

	// ErrGoogleClientNil is returned when the google client is nil.
	ErrGoogleClientNil = fmt.Errorf("google: google client is required")
)

// DirectoryService represent the  Google Directory API client.
type DirectoryService struct {
	svc *admin.Service
}

type DirectoryServiceConfig struct {
	Client         *http.Client
	UserEmail      string
	ServiceAccount []byte
	Scopes         []string
	UserAgent      string
}

// NewService create a Google Directory Service.
// References:
// - https://pkg.go.dev/google.golang.org/api/admin/directory/v1
// Examples of scope:
// - "https://www.googleapis.com/auth/admin.directory.group.readonly"
// - "https://www.googleapis.com/auth/admin.directory.group.member.readonly"
// - "https://www.googleapis.com/auth/admin.directory.user.readonly"
func NewService(ctx context.Context, config DirectoryServiceConfig) (*admin.Service, error) {
	if config.Client == nil {
		return nil, ErrGoogleClientNil
	}

	if config.UserEmail == "" {
		return nil, ErrUserEmailNil
	}

	if config.ServiceAccount == nil {
		return nil, ErrServiceAccountNil
	}

	if len(config.Scopes) == 0 {
		return nil, ErrGoogleClientScopeNil
	}

	if config.UserAgent == "" {
		return nil, ErrUserAgentNil
	}

	creds, err := google.CredentialsFromJSONWithParams(ctx, config.ServiceAccount, google.CredentialsParams{
		Scopes:  config.Scopes,
		Subject: config.UserEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("google: %v", err)
	}

	config.Client.Transport = &oauth2.Transport{
		Source: creds.TokenSource,
		Base:   config.Client.Transport,
	}

	svc, err := admin.NewService(
		ctx,
		option.WithUserAgent(config.UserAgent),
		option.WithHTTPClient(config.Client),
	)
	if err != nil {
		return nil, fmt.Errorf("google: %v", err)
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// optimistic initial capacity
	u := make([]*admin.User, 0, 50)

	if len(query) > 0 {
		for _, q := range query {
			if q != "" {
				err := ds.svc.Users.List().Query(q).Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list users with query %q: %w", q, err)
				}
			} else {
				err := ds.svc.Users.List().Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list users: %w", err)
				}
			}
		}
	} else {
		err := ds.svc.Users.List().Customer("my_customer").Fields(listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
			u = append(u, users.Users...)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("google: failed to list users: %w", err)
		}
	}

	return u, nil
}

// ListGroups list all groups in a Google Directory filtered by query.
// References:
// - https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups
func (ds *DirectoryService) ListGroups(ctx context.Context, query []string) ([]*admin.Group, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// optimistic initial capacity
	g := make([]*admin.Group, 0, 50)

	if len(query) > 0 {
		for _, q := range query {
			if q != "" {
				err := ds.svc.Groups.List().Customer("my_customer").Query(q).Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
					g = append(g, groups.Groups...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list groups with query %q: %w", q, err)
				}
			} else {
				err := ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
					g = append(g, groups.Groups...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list groups: %w", err)
				}
			}
		}
	} else {
		err := ds.svc.Groups.List().Customer("my_customer").Fields(groupsRequiredFields).Pages(ctx, func(groups *admin.Groups) error {
			g = append(g, groups.Groups...)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("google: failed to list groups: %w", err)
		}
	}

	return g, nil
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

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	qs := getGroupMembersOptions{}
	for _, q := range queries {
		q(&qs)
	}

	// optimistic initial capacity
	m := make([]*admin.Member, 0, 20)
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
				slog.Warn("google: member not included in group because status is not ACTIVE", "email", member.Email, "status", member.Status, "groupID", groupID)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
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

// GetUsersBatch retrieves multiple users by their email addresses using batch queries.
// This method is optimized for performance by using the ListUsers method with email queries
// instead of making individual GetUser calls.
func (ds *DirectoryService) GetUsersBatch(ctx context.Context, emails []string) ([]*admin.User, error) {
	if len(emails) == 0 {
		return []*admin.User{}, nil
	}

	// Use ListUsers with email queries for batch retrieval
	emailQueries := make([]string, len(emails))
	for i, email := range emails {
		emailQueries[i] = fmt.Sprintf("email:%s", email)
	}

	// Join queries with OR to get all users in one request
	query := strings.Join(emailQueries, " OR ")

	users, err := ds.ListUsers(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("google: error in batch get users: %w", err)
	}

	return users, nil
}

// ListGroupMembersBatch retrieves members for multiple groups concurrently.
// Returns a map where keys are group IDs and values are slices of members.
func (ds *DirectoryService) ListGroupMembersBatch(ctx context.Context, groupIDs []string, queries ...GetGroupMembersOption) (map[string][]*admin.Member, error) {
	if len(groupIDs) == 0 {
		return make(map[string][]*admin.Member), nil
	}

	result := make(map[string][]*admin.Member, len(groupIDs))

	// Process groups concurrently with a reasonable limit
	const maxConcurrent = 10
	sem := make(chan struct{}, maxConcurrent)
	var mu sync.Mutex
	var wg sync.WaitGroup

	errChan := make(chan error, len(groupIDs))

	for _, groupID := range groupIDs {
		wg.Add(1)
		go func(gid string) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			members, err := ds.ListGroupMembers(ctx, gid, queries...)
			if err != nil {
				errChan <- fmt.Errorf("google: error getting members for group %s: %w", gid, err)
				return
			}

			mu.Lock()
			result[gid] = members
			mu.Unlock()
		}(groupID)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	if err := <-errChan; err != nil {
		return nil, err
	}

	return result, nil
}
