package idp

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

var ErrDirectoryServiceNil = errors.New("provoder: directory service is nil")

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/idp/idp_mocks.go -source=idp.go GoogleProviderService

// This implement core.IdentityProviderService interface
// and as a consumer define GoogleProviderService to use pkg/google methods

type GoogleProviderService interface {
	ListUsers(ctx context.Context, query []string) ([]*admin.User, error)
	ListGroups(ctx context.Context, query []string) ([]*admin.Group, error)
	ListGroupMembers(ctx context.Context, groupID string) ([]*admin.Member, error)
	GetUser(ctx context.Context, userID string) (*admin.User, error)
}

type IdentityProvider struct {
	ps GoogleProviderService
}

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
func (i *IdentityProvider) GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error) {
	syncGroups := make([]model.Group, 0)

	pGroups, err := i.ps.ListGroups(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error listing groups: %w", err)
	}

	for _, grp := range pGroups {

		e := model.Group{
			IPID:  grp.Id,
			Name:  grp.Name,
			Email: grp.Email,
		}
		e.HashCode = hash.Get(e)

		syncGroups = append(syncGroups, e)
	}

	syncResult := &model.GroupsResult{
		Items:     len(pGroups),
		Resources: syncGroups,
	}

	syncResult.HashCode = hash.Get(syncResult)

	return syncResult, nil
}

// GetUsers returns a list of users from the Identity Provider API.
//
// The filter parameter is a list of strings that can be used to filter the users
// according to the Identity Provider API.
func (i *IdentityProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
	syncUsers := make([]model.User, 0)

	pUsers, err := i.ps.ListUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("idp: error listing users: %w", err)
	}

	for _, usr := range pUsers {

		e := model.User{
			IPID:        usr.Id,
			Name:        model.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		}
		e.HashCode = hash.Get(e)

		syncUsers = append(syncUsers, e)
	}

	uResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: syncUsers,
	}
	uResult.HashCode = hash.Get(uResult)

	return uResult, nil
}

// GetGroupMembers returns a list of members from the Identity Provider API.
func (i *IdentityProvider) GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error) {
	syncMembers := make([]model.Member, 0)

	pMembers, err := i.ps.ListGroupMembers(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("idp: error listing group members: %w", err)
	}

	for _, member := range pMembers {
		e := model.Member{
			IPID:  member.Id,
			Email: member.Email,
		}
		e.HashCode = hash.Get(e)

		syncMembers = append(syncMembers, e)
	}

	syncMembersResult := &model.MembersResult{
		Items:     len(pMembers),
		Resources: syncMembers,
	}
	syncMembersResult.HashCode = hash.Get(syncMembersResult)

	return syncMembersResult, nil
}

// GetUsersFromGroupMembers returns a list of users from the Identity Provider API.
func (i *IdentityProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error) {
	pUsers := make([]model.User, 0)

	for _, member := range mbr.Resources {
		u, err := i.ps.GetUser(ctx, member.IPID)
		if err != nil {
			return nil, fmt.Errorf("idp: error getting user: %w", err)
		}

		e := model.User{
			IPID:        u.Id,
			Name:        model.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
			Active:      !u.Suspended,
			Email:       u.PrimaryEmail,
		}
		e.HashCode = hash.Get(e)

		pUsers = append(pUsers, e)
	}

	pUsersResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}
	pUsersResult.HashCode = hash.Get(pUsersResult)

	return pUsersResult, nil
}

// GetUsersAndGroupsUsers returnpUserss a model.UsersResult and model.GroupsUsersResult data structures with the users and groups
func (i *IdentityProvider) GetUsersAndGroupsUsers(ctx context.Context, groups *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error) {
	// make pUsers unique
	userSet := make(map[string]struct{})

	pUsers := make([]model.User, 0)
	pGroupsUsers := make([]model.GroupUsers, 0)

	for _, pGroup := range groups.Resources {

		pMembers, err := i.GetGroupMembers(ctx, pGroup.IPID)
		if err != nil {
			return nil, nil, fmt.Errorf("idp: error getting group members: %w", err)
		}

		pUsersFromMembers, err := i.GetUsersFromGroupMembers(ctx, pMembers)
		if err != nil {
			return nil, nil, fmt.Errorf("idp: error getting users from group members: %w", err)
		}

		for _, pUser := range pUsersFromMembers.Resources {
			if _, ok := userSet[pUser.IPID]; !ok {
				pUsers = append(pUsers, pUser)
				userSet[pUser.IPID] = struct{}{}
			}
		}

		pGroupUsers := model.GroupUsers{
			Items: len(pMembers.Resources),
			Group: model.Group{
				IPID:  pGroup.IPID,
				Name:  pGroup.Name,
				Email: pGroup.Email,
			},
			Resources: pUsers,
		}
		pGroupUsers.HashCode = hash.Get(pGroupUsers)

		pGroupsUsers = append(pGroupsUsers, pGroupUsers)
	}

	usersResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}
	usersResult.HashCode = hash.Get(usersResult)

	groupsUsersResult := &model.GroupsUsersResult{
		Items:     len(pGroupsUsers),
		Resources: pGroupsUsers,
	}
	groupsUsersResult.HashCode = hash.Get(groupsUsersResult)

	return usersResult, groupsUsersResult, nil
}
