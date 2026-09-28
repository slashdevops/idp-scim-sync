package google

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/pprof"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/slashdevops/idp-scim-sync/internal/model"
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
	ErrGoogleClientScopeNil = errors.New("google: google client scope is required")

	// ErrUserIDNil is returned when the user ID is nil.
	ErrUserIDNil = errors.New("google: user id is required")

	// ErrUserEmailNil is returned when the user email is nil.
	ErrUserEmailNil = errors.New("google: user email is required")

	// ErrGroupIDNil is returned when the group ID is nil.
	ErrGroupIDNil = errors.New("google: group id is required")

	// ErrServiceAccountNil is returned when the service account credentials are nil.
	ErrServiceAccountNil = errors.New("google: service account credentials are required")

	// ErrUserAgentNil is returned when the user agent is nil.
	ErrUserAgentNil = errors.New("google: user agent is required")

	// ErrGoogleClientNil is returned when the google client is nil.
	ErrGoogleClientNil = errors.New("google: google client is required")
)

// DirectoryService represent the  Google Directory API client.
type DirectoryService struct {
	svc                     *admin.Service
	listUsersRequiredFields googleapi.Field
	getUsersRequiredFields  googleapi.Field
}

// DirectoryServiceConfig carries everything NewService needs to build an
// authenticated Google Directory API client.
//
// All fields are required. ServiceAccount holds the raw service-account JSON
// (not a path), and UserEmail is the Workspace admin the service account
// impersonates via domain-wide delegation.
type DirectoryServiceConfig struct {
	Client         *http.Client
	UserEmail      string
	UserAgent      string
	ServiceAccount []byte
	Scopes         []string
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

	// TODO(#go-1.27-modernization): migrate off CredentialsFromJSONWithParams.
	//
	// It is deprecated for a security reason — it does not validate the
	// credential configuration — and here the service-account JSON arrives at
	// runtime from AWS Secrets Manager (or a file on disk for local runs) and is
	// never validated before use. The replacement is
	// cloud.google.com/go/auth/credentials.DetectDefault with CredentialsOptions,
	// or CredentialsFromJSON plus explicit up-front validation of type,
	// client_email and private_key.
	//
	// Deliberately not changed in this pass: this is the authentication path of
	// a production Lambda, so it needs its own change, its own tests, and a
	// manual run against a real Google Workspace tenant. Tracked as finding D1.
	//
	//nolint:staticcheck // SA1019: see the TODO above; migration is scheduled separately.
	creds, err := google.CredentialsFromJSONWithParams(ctx, config.ServiceAccount, google.CredentialsParams{
		Scopes:  config.Scopes,
		Subject: config.UserEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("google: %w", err)
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
		return nil, fmt.Errorf("google: %w", err)
	}

	return svc, nil
}

// DirectoryServiceOption is a function that configures a DirectoryService.
type DirectoryServiceOption func(*DirectoryService)

// WithSyncFieldSet configures the DirectoryService to only request fields
// needed for the configured sync field set from the Google API.
// When fields is nil or empty, all user fields are requested (default behavior).
func WithSyncFieldSet(fields *model.SyncFieldSet) DirectoryServiceOption {
	return func(ds *DirectoryService) {
		uf := buildUserFields(fields)
		ds.listUsersRequiredFields = googleapi.Field("nextPageToken, users(" + uf + ")")
		ds.getUsersRequiredFields = googleapi.Field(uf)
	}
}

// buildUserFields constructs the Google API fields parameter based on the configured field set.
func buildUserFields(fields *model.SyncFieldSet) string {
	// Always include required fields
	parts := []string{baseFields, "primaryEmail", "name", "suspended", "kind", "emails"}

	if fields.Includes(model.SyncUserFieldAddresses) {
		parts = append(parts, "addresses")
	}
	if fields.Includes(model.SyncUserFieldPhoneNumbers) {
		parts = append(parts, "phones")
	}
	if fields.Includes(model.SyncUserFieldPreferredLanguage) {
		parts = append(parts, "languages")
	}
	if fields.Includes(model.SyncUserFieldTitle) || fields.Includes(model.SyncUserFieldEnterpriseData) {
		parts = append(parts, "organizations")
	}
	if fields.Includes(model.SyncUserFieldEnterpriseData) {
		parts = append(parts, "relations")
	}
	// locations is currently not mapped to any SCIM attribute, but include it
	// when all fields are synced for backward compatibility
	if fields == nil || fields.IsEmpty() {
		parts = append(parts, "locations")
	}

	return strings.Join(parts, ",")
}

// NewDirectoryService create a Google Directory API client.
// References:
// - https://developers.google.com/admin-sdk/directory/v1/guides/delegation?utm_source=pocket_mylist#go
func NewDirectoryService(svc *admin.Service, opts ...DirectoryServiceOption) (*DirectoryService, error) {
	ds := &DirectoryService{
		svc:                     svc,
		listUsersRequiredFields: listUsersRequiredFields,
		getUsersRequiredFields:  getUsersRequiredFields,
	}

	for _, opt := range opts {
		opt(ds)
	}

	return ds, nil
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
				slog.Debug("google: Listing users with query", "query", q)
				err := ds.svc.Users.List().Query(q).Customer("my_customer").Fields(ds.listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					slog.Debug("google: Retrieved users page", "page_size", len(users.Users))
					u = append(u, users.Users...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list users with query %q: %w", q, err)
				}
			} else {
				err := ds.svc.Users.List().Customer("my_customer").Fields(ds.listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
					u = append(u, users.Users...)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("google: failed to list users: %w", err)
				}
			}
		}
	} else {
		err := ds.svc.Users.List().Customer("my_customer").Fields(ds.listUsersRequiredFields).Pages(ctx, func(users *admin.Users) error {
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

	u, err := ds.svc.Users.Get(userID).Fields(ds.getUsersRequiredFields).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("google: error getting user %s: %w", userID, err)
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
		return nil, fmt.Errorf("google: error getting group %s: %w", groupID, err)
	}

	return g, nil
}

// listGroupMembersBatchConcurrency caps the number of in-flight members-list
// requests issued by ListGroupMembersBatch. Ten keeps the fan-out well inside
// the Google Directory API's per-user quota while still overlapping the
// round-trip latency that dominates a sync.
const listGroupMembersBatchConcurrency = 10

// ListGroupMembersBatch retrieves members for multiple groups concurrently and
// returns a map keyed by group ID.
//
// The fan-out is bounded by listGroupMembersBatchConcurrency and shares one
// errgroup context, so the first failure cancels the requests still in flight
// and the queued ones never start. That matters because every request spends
// Google Directory API quota whose result would be discarded anyway: the first
// error is what the caller receives.
func (ds *DirectoryService) ListGroupMembersBatch(ctx context.Context, groupIDs []string, queries ...GetGroupMembersOption) (map[string][]*admin.Member, error) {
	if len(groupIDs) == 0 {
		return make(map[string][]*admin.Member), nil
	}

	result := make(map[string][]*admin.Member, len(groupIDs))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(listGroupMembersBatchConcurrency)

	var mu sync.Mutex

	for _, groupID := range groupIDs {
		g.Go(func() error {
			// See internal/idp.GetUsersByGroupsMembers for why this is
			// SetGoroutineLabels rather than pprof.Do.
			pprof.SetGoroutineLabels(pprof.WithLabels(ctx, pprof.Labels(
				"sync", "group-members",
				"group", groupID,
			)))

			members, err := ds.ListGroupMembers(ctx, groupID, queries...)
			if err != nil {
				return fmt.Errorf("google: error getting members for group %s: %w", groupID, err)
			}

			mu.Lock()
			result[groupID] = members
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
