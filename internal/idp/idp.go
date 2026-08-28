package idp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	admin "google.golang.org/api/admin/directory/v1"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
)

// This implement core.IdentityProviderService interface

var (
	// ErrDirectoryServiceNil is returned when the GoogleProviderService is nil.
	ErrDirectoryServiceNil = errors.New("provider: directory service is nil")

	// ErrGroupIDNil is returned when the groupID is nil.
	ErrGroupIDNil = errors.New("provider: group id is nil")

	// ErrGroupResultNil is returned when the group result is nil.
	ErrGroupResultNil = errors.New("provider: group result is nil")
)

// getUsersConcurrency caps the number of in-flight GetUser requests issued by
// GetUsersByGroupsMembers. Ten keeps the fan-out well inside the Google
// Directory API's per-user quota while still overlapping round-trip latency.
const getUsersConcurrency = 10

//go:generate go tool mockgen -package=mocks -destination=../../mocks/idp/idp_mocks.go -source=idp.go GoogleProviderService

// GoogleProviderService is the interface that wraps the Google Provider Service methods.
type GoogleProviderService interface {
	ListUsers(ctx context.Context, query []string) ([]*admin.User, error)
	ListGroups(ctx context.Context, query []string) ([]*admin.Group, error)
	ListGroupMembers(ctx context.Context, groupID string, queries ...google.GetGroupMembersOption) ([]*admin.Member, error)
	GetUser(ctx context.Context, userID string) (*admin.User, error)

	// Batch operations for performance optimization
	ListGroupMembersBatch(ctx context.Context, groupIDs []string, queries ...google.GetGroupMembersOption) (map[string][]*admin.Member, error)
}

// IdentityProvider is the Identity Provider service that implements the core.IdentityProvider interface and consumes the pkg.google methods.
type IdentityProvider struct {
	ps           GoogleProviderService
	syncFieldSet *model.SyncFieldSet
}

// IdentityProviderOption is a function that configures an IdentityProvider.
type IdentityProviderOption func(*IdentityProvider)

// WithSyncFieldSet configures which optional user fields are included in the sync.
// When the field set is nil or empty, all fields are synced (default behavior).
func WithSyncFieldSet(fields *model.SyncFieldSet) IdentityProviderOption {
	return func(ip *IdentityProvider) {
		ip.syncFieldSet = fields
	}
}

// NewIdentityProvider returns a new instance of the Identity Provider service.
func NewIdentityProvider(gps GoogleProviderService, opts ...IdentityProviderOption) (*IdentityProvider, error) {
	if gps == nil {
		return nil, ErrDirectoryServiceNil
	}

	ip := &IdentityProvider{
		ps: gps,
	}

	for _, opt := range opts {
		opt(ip)
	}

	return ip, nil
}

// GetGroups returns a list of groups from the Identity Provider API.
//
// The filter parameter is a list of strings that can be used to filter the groups
// according to the Identity Provider API.
//
// This method checks the names of the groups and avoid the second, third, etc repetition of the same group name.
func (i *IdentityProvider) GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error) {
	pGroups, err := i.ps.ListGroups(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error getting groups: %w", err)
	}

	if len(pGroups) == 0 {
		syncGroups := make([]*model.Group, 0)
		gResult := model.GroupsResultBuilder().WithResources(syncGroups).Build()
		return gResult, nil
	}

	uniqueGroups := make(map[string]struct{}, len(pGroups))
	syncGroups := make([]*model.Group, 0, len(pGroups))
	for _, grp := range pGroups {
		// this is a hack to avoid the second, third, etc repetition of the same group name
		if _, ok := uniqueGroups[grp.Name]; !ok {
			uniqueGroups[grp.Name] = struct{}{}

			gg := model.GroupBuilder().
				WithIPID(grp.Id).
				WithName(grp.Name).
				WithEmail(grp.Email).
				Build()

			syncGroups = append(syncGroups, gg)
		} else {
			slog.Warn("idp: group already exists with the same name, this group will be avoided, please make your groups uniques by name!",
				"id", grp.Id,
				"name", grp.Name,
				"email", grp.Email,
			)
		}
	}

	syncResult := model.GroupsResultBuilder().WithResources(syncGroups).Build()
	slog.Debug("idp: GetGroups()", "groups", len(syncGroups))

	return syncResult, nil
}

// GetUsers returns a list of users from the Identity Provider API.
//
// The filter parameter is a list of strings that can be used to filter the users
// according to the Identity Provider API.
func (i *IdentityProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
	pUsers, err := i.ps.ListUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error getting users: %w", err)
	}

	if len(pUsers) == 0 {
		syncUsers := make([]*model.User, 0)
		uResult := model.UsersResultBuilder().WithResources(syncUsers).Build()
		return uResult, nil
	}

	syncUsers := make([]*model.User, 0, len(pUsers))
	for _, usr := range pUsers {
		gu := buildUser(usr, i.syncFieldSet)
		if gu == nil {
			// buildUser rejects records the AWS SCIM API would refuse (no name,
			// given name, family name or primary email) and logs the reason.
			// Skip them: a nil element travels on into reconciliation and
			// hashing, where it is dereferenced.
			continue
		}
		syncUsers = append(syncUsers, gu)
	}
	uResult := model.UsersResultBuilder().WithResources(syncUsers).Build()
	slog.Debug("idp: GetUsers()", "users", len(syncUsers))

	return uResult, nil
}

// GetGroupMembers returns a list of members from the Identity Provider API.
func (i *IdentityProvider) GetGroupMembers(ctx context.Context, groupID string) (*model.MembersResult, error) {
	if groupID == "" {
		return nil, ErrGroupIDNil
	}

	pMembers, err := i.ps.ListGroupMembers(ctx, groupID, google.WithIncludeDerivedMembership(true))
	if err != nil {
		return nil, fmt.Errorf("idp: error getting group members: %w", err)
	}

	if len(pMembers) == 0 {
		syncMembers := make([]*model.Member, 0)
		membersResult := model.MembersResultBuilder().WithResources(syncMembers).Build()
		return membersResult, nil
	}

	syncMembers := make([]*model.Member, 0, len(pMembers))
	for _, member := range pMembers {
		// avoid nested groups, but members are included thanks to the google.WithIncludeDerivedMembership option above
		if member.Type == "GROUP" {
			slog.Warn("skipping member because is a group, but group members will be included",
				"id", member.Id,
				"email", member.Email,
			)
			continue
		}

		gm := model.MemberBuilder().
			WithIPID(member.Id).
			WithEmail(member.Email).
			WithStatus(member.Status).
			Build()

		syncMembers = append(syncMembers, gm)
	}

	syncMembersResult := model.MembersResultBuilder().WithResources(syncMembers).Build()

	return syncMembersResult, nil
}

// GetUsersByGroupsMembers returns a list of users from the Identity Provider API.
// It fetches user details concurrently with a bounded number of goroutines.
func (i *IdentityProvider) GetUsersByGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) (*model.UsersResult, error) {
	if gmr == nil {
		return nil, ErrGroupResultNil
	}

	if len(gmr.Resources) == 0 {
		syncUsers := make([]*model.User, 0)
		uResult := model.UsersResultBuilder().WithResources(syncUsers).Build()
		return uResult, nil
	}

	// Collect unique emails first
	uniqEmails := make(map[string]string, len(gmr.Resources)) // email -> IPID
	for _, groupMembers := range gmr.Resources {
		for _, member := range groupMembers.Resources {
			if _, ok := uniqEmails[member.Email]; !ok {
				uniqEmails[member.Email] = member.IPID
			}
		}
	}

	// Fetch users concurrently. One errgroup context is shared by every
	// goroutine, so the first failure cancels the requests still in flight and
	// the queued ones never start: their results would be discarded anyway,
	// since the first error is what the caller receives.
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(getUsersConcurrency)

	var mu sync.Mutex
	pUsers := make([]*model.User, 0, len(uniqEmails))

	for email, ipid := range uniqEmails {
		g.Go(func() error {
			u, err := i.ps.GetUser(ctx, email)
			if err != nil {
				return fmt.Errorf("idp: error getting user: %+v, email: %s, error: %w", ipid, email, err)
			}

			gu := buildUser(u, i.syncFieldSet)
			if gu == nil {
				// Rejected by buildUser, which already logged why. Skipping
				// keeps a single malformed directory record from crashing the
				// whole sync.
				return nil
			}

			mu.Lock()
			pUsers = append(pUsers, gu)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// uniqEmails is a map, so the fan-out completes in an arbitrary order.
	// UsersResult.HashCode is order-insensitive (SetHashCode sorts a copy), but
	// the state file written to S3 preserves slice order, so an unsorted result
	// made every sync produce a large spurious diff. UserName is the user's
	// primary email and is always populated by buildUser.
	slices.SortFunc(pUsers, func(a, b *model.User) int {
		if c := strings.Compare(a.UserName, b.UserName); c != 0 {
			return c
		}
		return strings.Compare(a.IPID, b.IPID)
	})

	pUsersResult := model.UsersResultBuilder().WithResources(pUsers).Build()

	slog.Debug("idp: GetUsersByGroupsMembers()", "users", len(pUsers))

	return pUsersResult, nil
}

// GetGroupsMembers return the members of the groups
func (i *IdentityProvider) GetGroupsMembers(ctx context.Context, gr *model.GroupsResult) (*model.GroupsMembersResult, error) {
	if gr == nil {
		return nil, ErrGroupResultNil
	}

	l := len(gr.Resources)
	if l == 0 {
		groupsMembersResult := &model.GroupsMembersResult{
			Items:     l,
			Resources: make([]*model.GroupMembers, l),
		}
		groupsMembersResult.SetHashCode()

		return groupsMembersResult, nil
	}

	// Collect all group IDs for batch operation
	groupIDs := make([]string, l)
	groupsByID := make(map[string]*model.Group, l)
	for i, group := range gr.Resources {
		groupIDs[i] = group.IPID
		groupsByID[group.IPID] = group
	}

	// Use batch operation to get all group members at once
	membersMap, err := i.ps.ListGroupMembersBatch(ctx, groupIDs, google.WithIncludeDerivedMembership(true))
	if err != nil {
		return nil, fmt.Errorf("idp: error getting group members batch: %w", err)
	}

	// Process the results
	groupMembers := make([]*model.GroupMembers, 0, l)
	for groupID, pMembers := range membersMap {
		group := groupsByID[groupID]

		syncMembers := make([]*model.Member, 0, len(pMembers))
		for _, member := range pMembers {
			// avoid nested groups, but members are included thanks to the google.WithIncludeDerivedMembership option above
			if member.Type == "GROUP" {
				slog.Warn("skipping member because is a group, but group members will be included",
					"id", member.Id,
					"email", member.Email,
				)
				continue
			}

			gm := model.MemberBuilder().
				WithIPID(member.Id).
				WithEmail(member.Email).
				WithStatus(member.Status).
				Build()

			syncMembers = append(syncMembers, gm)
		}

		ggm := model.GroupBuilder().
			WithIPID(group.IPID).
			WithName(group.Name).
			WithEmail(group.Email).
			Build()

		groupMember := model.GroupMembersBuilder().
			WithGroup(ggm).
			WithResources(syncMembers).
			Build()

		groupMembers = append(groupMembers, groupMember)
	}

	groupsMembersResult := &model.GroupsMembersResult{
		Items:     len(groupMembers),
		Resources: groupMembers,
	}
	groupsMembersResult.SetHashCode()

	slog.Debug("idp: GetGroupsMembers()", "groups", len(groupMembers))

	return groupsMembersResult, nil
}
