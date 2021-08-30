package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrDirectoryServiceNil = errors.New("directory service is nil")
	ErrListingGroups       = errors.New("error listing groups")
	ErrListingUsers        = errors.New("error listing users")
	ErrListingGroupMembers = errors.New("error listing group members")
	ErrGettingUser         = errors.New("error getting user")
)

// This implement core.IdentityProviderService interface

type GoogleProvider struct {
	ds DirectoryService
}

func NewGoogleIdentityProvider(ds DirectoryService) (*GoogleProvider, error) {
	if ds == nil {
		return nil, ErrDirectoryServiceNil
	}

	return &GoogleProvider{
		ds: ds,
	}, nil
}

// GetGroups returns a list of groups from the Google Directory API.
//
// The filter parameter is a list of strings that can be used to filter the groups
// according to the Google Directory API.
func (g *GoogleProvider) GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error) {
	syncGroups := make([]*model.Group, 0)

	googleGroups, err := g.ds.ListGroups(filter)
	if err != nil {
		return nil, ErrListingGroups
	}

	for _, grp := range googleGroups {
		syncGroups = append(syncGroups, &model.Group{
			ID:    grp.Id,
			Name:  grp.Name,
			Email: grp.Email,
		})
	}

	syncResult := &model.GroupsResult{
		Items:     len(googleGroups),
		Resources: syncGroups,
	}

	return syncResult, nil
}

// GetUsers returns a list of users from the Google Directory API.
//
// The filter parameter is a list of strings that can be used to filter the users
// according to the Google Directory API.
func (g *GoogleProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
	syncUsers := make([]*model.User, 0)

	googleUsers, err := g.ds.ListUsers(filter)
	if err != nil {
		return nil, ErrListingUsers
	}

	for _, usr := range googleUsers {
		syncUsers = append(syncUsers, &model.User{
			ID:          usr.Id,
			Name:        model.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		})
	}

	uResult := &model.UsersResult{
		Items:     len(googleUsers),
		Resources: syncUsers,
	}

	return uResult, nil
}

func (g *GoogleProvider) GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error) {
	syncMembers := make([]*model.Member, 0)

	googleMembers, err := g.ds.ListGroupMembers(id)
	if err != nil {
		return nil, ErrListingGroupMembers
	}

	for _, member := range googleMembers {
		syncMembers = append(syncMembers, &model.Member{
			ID:    member.Id,
			Email: member.Email,
		})
	}

	syncMembersResult := &model.MembersResult{
		Items:     len(googleMembers),
		Resources: syncMembers,
	}

	return syncMembersResult, nil
}

func (g *GoogleProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error) {
	syncUsers := make([]*model.User, 0)

	for _, member := range mbr.Resources {
		u, err := g.ds.GetUser(member.ID)
		if err != nil {
			return nil, ErrGettingUser
		}

		syncUsers = append(syncUsers, &model.User{
			ID:          u.Id,
			Name:        model.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
			Active:      !u.Suspended,
			Email:       u.PrimaryEmail,
		})
	}

	syncUsersResult := &model.UsersResult{
		Items:     len(syncUsers),
		Resources: syncUsers,
	}

	return syncUsersResult, nil
}

func (g *GoogleProvider) GetUsersAndGroupsUsers(ctx context.Context, groups *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error) {
	pUsers := make([]*model.User, 0)
	pGroupsUsers := make([]*model.GroupUsers, 0)

	for _, pGroup := range groups.Resources {

		pMembers, err := g.GetGroupMembers(ctx, pGroup.ID)
		if err != nil {
			return nil, nil, err
		}

		pUsersFromMembers, err := g.GetUsersFromGroupMembers(ctx, pMembers)
		if err != nil {
			return nil, nil, err
		}
		pUsers = append(pUsers, pUsersFromMembers.Resources...)

		pGroupUsers := &model.GroupUsers{
			Items: len(pMembers.Resources),
			Group: model.Group{
				ID:    pGroup.ID,
				Name:  pGroup.Name,
				Email: pGroup.Email,
			},
			Resources: pUsers,
		}

		pGroupsUsers = append(pGroupsUsers, pGroupUsers)
	}

	usersResult := &model.UsersResult{
		Items:     len(pUsers),
		Resources: pUsers,
	}

	groupsUsersResult := &model.GroupsUsersResult{
		Items:     len(pGroupsUsers),
		Resources: pGroupsUsers,
	}

	return usersResult, groupsUsersResult, nil
}
