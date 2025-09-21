package idp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	admin "google.golang.org/api/admin/directory/v1"
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

//go:generate go tool mockgen -package=mocks -destination=../../mocks/idp/idp_mocks.go -source=idp.go GoogleProviderService

// GoogleProviderService is the interface that wraps the Google Provider Service methods.
type GoogleProviderService interface {
	ListUsers(ctx context.Context, query []string) ([]*admin.User, error)
	ListGroups(ctx context.Context, query []string) ([]*admin.Group, error)
	ListGroupMembers(ctx context.Context, groupID string, queries ...google.GetGroupMembersOption) ([]*admin.Member, error)
	GetUser(ctx context.Context, userID string) (*admin.User, error)

	// Batch operations for performance optimization
	GetUsersBatch(ctx context.Context, emails []string) ([]*admin.User, error)
	ListGroupMembersBatch(ctx context.Context, groupIDs []string, queries ...google.GetGroupMembersOption) (map[string][]*admin.Member, error)
}

// IdentityProvider is the Identity Provider service that implements the core.IdentityProvider interface and consumes the pkg.google methods.
type IdentityProvider struct {
	ps GoogleProviderService
}

// NewIdentityProvider returns a new instance of the Identity Provider service.
func NewIdentityProvider(gps GoogleProviderService) (*IdentityProvider, error) {
	if gps == nil {
		return nil, ErrDirectoryServiceNil
	}

	return &IdentityProvider{
		ps: gps,
	}, nil
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

	// Pre-allocate with exact capacity and filter nil users upfront
	syncUsers := make([]*model.User, 0, len(pUsers))
	for _, usr := range pUsers {
		if gu := buildUser(usr); gu != nil {
			syncUsers = append(syncUsers, gu)
		}
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

	slog.Debug("idp: GetGroupMembers()", "members", len(syncMembers))

	return syncMembersResult, nil
}

// GetUsersByGroupsMembers returns a list of users from the Identity Provider API.
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
	uniqueEmails := make(map[string]struct{})
	for _, groupMembers := range gmr.Resources {
		for _, member := range groupMembers.Resources {
			uniqueEmails[member.Email] = struct{}{}
		}
	}

	// Convert to slice for batch processing
	emails := make([]string, 0, len(uniqueEmails))
	for email := range uniqueEmails {
		emails = append(emails, email)
	}

	// Process emails in chunks of 500 (Google's recommended batch size)
	const batchSize = 500
	pUsers := make([]*model.User, 0, len(emails))

	emailChunks := chunkEmails(emails, batchSize)
	for _, emailBatch := range emailChunks {
		query := buildEmailQuery(emailBatch)

		users, err := i.ps.ListUsers(ctx, []string{query})
		if err != nil {
			return nil, fmt.Errorf("idp: error getting users batch: %w", err)
		}

		for _, usr := range users {
			if gu := buildUser(usr); gu != nil {
				pUsers = append(pUsers, gu)
				slog.Debug("idp: GetUsersByGroupsMembers()", "user", gu.Email)
			}
		}
	}

	pUsersResult := model.UsersResultBuilder().WithResources(pUsers).Build()
	slog.Debug("idp: GetUsersByGroupsMembers()", "users", len(pUsers))

	return pUsersResult, nil
}

// GetGroupsMembers return the members of the groups with parallel processing for improved performance
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

	// Add timeout to prevent hanging operations
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Use worker pool for concurrent processing
	const maxWorkers = 10 // Adjust based on Google API rate limits
	workers := maxWorkers
	if l < workers {
		workers = l
	}

	type groupMemberJob struct {
		index int
		group *model.Group
	}

	type groupMemberResult struct {
		index       int
		groupMember *model.GroupMembers
		err         error
	}

	jobs := make(chan groupMemberJob, l)
	results := make(chan groupMemberResult, l)

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					results <- groupMemberResult{job.index, nil, ctx.Err()}
					return
				default:
					members, err := i.GetGroupMembers(ctx, job.group.IPID)
					if err != nil {
						results <- groupMemberResult{job.index, nil, err}
						continue
					}

					ggm := model.GroupBuilder().
						WithIPID(job.group.IPID).
						WithName(job.group.Name).
						WithEmail(job.group.Email).
						Build()

					groupMember := model.GroupMembersBuilder().
						WithGroup(ggm).
						WithResources(members.Resources).
						Build()

					results <- groupMemberResult{job.index, groupMember, nil}
				}
			}
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i, group := range gr.Resources {
			select {
			case <-ctx.Done():
				return
			case jobs <- groupMemberJob{i, group}:
			}
		}
	}()

	// Collect results
	groupMembers := make([]*model.GroupMembers, l)
	for i := 0; i < l; i++ {
		select {
		case <-ctx.Done():
			// Wait for workers to finish
			wg.Wait()
			return nil, fmt.Errorf("idp: context cancelled while getting group members: %w", ctx.Err())
		case result := <-results:
			if result.err != nil {
				// Wait for workers to finish
				wg.Wait()
				return nil, fmt.Errorf("idp: error getting group members: %w", result.err)
			}
			groupMembers[result.index] = result.groupMember
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	groupsMembersResult := &model.GroupsMembersResult{
		Items:     len(groupMembers),
		Resources: groupMembers,
	}
	groupsMembersResult.SetHashCode()

	slog.Debug("idp: GetGroupsMembers()", "groups", len(groupMembers))

	return groupsMembersResult, nil
}
