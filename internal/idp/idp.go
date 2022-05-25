package idp

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
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

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/idp/idp_mocks.go -source=idp.go GoogleProviderService

// GoogleProviderService is the interface that wraps the Google Provider Service methods.
type GoogleProviderService interface {
	ListUsers(ctx context.Context, query []string) ([]*admin.User, error)
	ListGroups(ctx context.Context, query []string) ([]*admin.Group, error)
	ListGroupMembers(ctx context.Context, groupID string, queries ...google.GetGroupMembersOption) ([]*admin.Member, error)
	GetUser(ctx context.Context, userID string) (*admin.User, error)
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
	uniqueGroups := make(map[string]struct{})
	syncGroups := make([]*model.Group, 0)

	pGroups, err := i.ps.ListGroups(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error listing groups: %w", err)
	}

	for _, grp := range pGroups {
		// this is a hack to avoid the second, third, etc repetition of the same group name
		if _, ok := uniqueGroups[grp.Name]; !ok {
			uniqueGroups[grp.Name] = struct{}{}

			e := &model.Group{
				IPID:  grp.Id,
				Name:  grp.Name,
				Email: grp.Email,
			}
			e.SetHashCode()

			syncGroups = append(syncGroups, e)
		} else {
			log.WithFields(log.Fields{
				"id":    grp.Id,
				"name":  grp.Name,
				"email": grp.Email,
			}).Warning("idp: group already exists with the same name, this group will be avoided, please make your groups uniques by name!")
		}
	}

	syncResult := &model.GroupsResult{
		Items:     len(syncGroups),
		Resources: syncGroups,
	}
	syncResult.SetHashCode()

	return syncResult, nil
}

// GetUsers returns a list of users from the Identity Provider API.
//
// The filter parameter is a list of strings that can be used to filter the users
// according to the Identity Provider API.
func (i *IdentityProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
	syncUsers := make([]*model.User, 0)

	pUsers, err := i.ps.ListUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error listing users: %w", err)
	}

	for _, usr := range pUsers {
		e := &model.User{
			IPID:        usr.Id,
			Name:        model.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		}
		e.SetHashCode()

		syncUsers = append(syncUsers, e)
	}

	uResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: syncUsers,
	}
	uResult.SetHashCode()

	return uResult, nil
}

// GetGroupMembers returns a list of members from the Identity Provider API.
func (i *IdentityProvider) GetGroupMembers(ctx context.Context, groupID string) (*model.MembersResult, error) {
	if groupID == "" {
		return nil, ErrGroupIDNil
	}

	syncMembers := make([]*model.Member, 0)

	pMembers, err := i.ps.ListGroupMembers(ctx, groupID, google.WithIncludeDerivedMembership(true))
	if err != nil {
		return nil, fmt.Errorf("idp: error listing group members: %w", err)
	}

	for _, member := range pMembers {
		// avoid nested groups, but members are included thanks to the google.WithIncludeDerivedMembership option above
		if member.Type == "GROUP" {
			log.WithFields(log.Fields{
				"id":    member.Id,
				"email": member.Email,
			}).Warn("skipping member because is a group, but group members will be included")
			continue
		}
		e := &model.Member{
			IPID:  member.Id,
			Email: member.Email,
		}
		e.SetHashCode()

		syncMembers = append(syncMembers, e)
	}

	syncMembersResult := &model.MembersResult{
		Items:     len(pMembers),
		Resources: syncMembers,
	}
	syncMembersResult.SetHashCode()

	return syncMembersResult, nil
}

// GetUsersByGroupsMembers returns a list of users from the Identity Provider API.
func (i *IdentityProvider) GetUsersByGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) (*model.UsersResult, error) {
	pUsers := make([]*model.User, 0)
	uniqUsers := make(map[string]struct{})

	for _, groupMembers := range gmr.Resources {
		for _, member := range groupMembers.Resources {
			u, err := i.ps.GetUser(ctx, member.Email)
			if err != nil {
				return nil, fmt.Errorf("idp: error getting user: %+v, email: %s, error: %w", member.IPID, member.Email, err)
			}

			e := &model.User{
				IPID:        u.Id,
				Name:        model.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
				DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
				Active:      !u.Suspended,
				Email:       u.PrimaryEmail,
			}
			e.SetHashCode()
			if _, ok := uniqUsers[e.Email]; !ok {
				uniqUsers[e.Email] = struct{}{}
				pUsers = append(pUsers, e)
			}
		}
	}
	pUsersResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}
	pUsersResult.SetHashCode()

	return pUsersResult, nil
}

// GetGroupsMembers return the members of the groups
func (i *IdentityProvider) GetGroupsMembers(ctx context.Context, gr *model.GroupsResult) (*model.GroupsMembersResult, error) {
	if gr == nil {
		return nil, ErrGroupResultNil
	}

	groupMembers := make([]*model.GroupMembers, 0)

	for _, group := range gr.Resources {
		members, err := i.GetGroupMembers(ctx, group.IPID)
		if err != nil {
			return nil, fmt.Errorf("idp: error getting group members: %w", err)
		}

		e := model.Group{
			IPID:  group.IPID,
			Name:  group.Name,
			Email: group.Email,
		}
		e.SetHashCode()

		if members.Items > 0 {
			groupMember := &model.GroupMembers{
				Items:     len(members.Resources),
				Group:     e,
				Resources: members.Resources,
			}
			groupMember.SetHashCode()

			groupMembers = append(groupMembers, groupMember)
		} else {
			groupMember := &model.GroupMembers{
				Items:     0,
				Group:     e,
				Resources: make([]*model.Member, 0),
			}
			groupMember.SetHashCode()

			groupMembers = append(groupMembers, groupMember)
		}
	}

	groupsMembersResult := &model.GroupsMembersResult{
		Items:     len(groupMembers),
		Resources: groupMembers,
	}
	groupsMembersResult.SetHashCode()

	return groupsMembersResult, nil
}
