package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	ErrDirectoryServiceNil = errors.New("directory service is nil")
	ErrListingGroups       = errors.New("error listing groups")
	ErrListingUsers        = errors.New("error listing users")
	ErrListingGroupMembers = errors.New("error listing group members")
	ErrGettingUser         = errors.New("error getting user")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/provider/identity_provider_mocks.go -source=identity_provider.go GoogleProviderService

// This implement core.IdentityProviderService interface
// and as a consumer define GoogleProviderService to use pkg.google methods

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

// GetGroups returns a list of groups from the Google Directory API.
//
// The filter parameter is a list of strings that can be used to filter the groups
// according to the Google Directory API.
func (i *IdentityProvider) GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error) {
	syncGroups := make([]model.Group, 0)

	googleGroups, err := i.ps.ListGroups(ctx, filter)
	if err != nil {
		return nil, ErrListingGroups
	}

	for _, grp := range googleGroups {

		e := model.Group{
			ID:    grp.Id,
			Name:  grp.Name,
			Email: grp.Email,
		}
		e.HashCode = hash.Get(e)

		syncGroups = append(syncGroups, e)
	}

	syncResult := &model.GroupsResult{
		Items:     len(googleGroups),
		Resources: syncGroups,
	}

	syncResult.HashCode = hash.Get(syncResult)

	return syncResult, nil
}

// GetUsers returns a list of users from the Google Directory API.
//
// The filter parameter is a list of strings that can be used to filter the users
// according to the Google Directory API.
func (i *IdentityProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
	syncUsers := make([]model.User, 0)

	googleUsers, err := i.ps.ListUsers(ctx, filter)
	if err != nil {
		return nil, ErrListingUsers
	}

	for _, usr := range googleUsers {

		e := model.User{
			ID:          usr.Id,
			Name:        model.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		}
		e.HashCode = hash.Get(e)

		syncUsers = append(syncUsers, e)
	}

	uResult := &model.UsersResult{
		Items:     len(googleUsers),
		Resources: syncUsers,
	}
	uResult.HashCode = hash.Get(uResult)

	return uResult, nil
}

func (i *IdentityProvider) GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error) {
	syncMembers := make([]model.Member, 0)

	googleMembers, err := i.ps.ListGroupMembers(ctx, id)
	if err != nil {
		return nil, ErrListingGroupMembers
	}

	for _, member := range googleMembers {
		e := model.Member{
			ID:    member.Id,
			Email: member.Email,
		}
		e.HashCode = hash.Get(e)

		syncMembers = append(syncMembers, e)
	}

	syncMembersResult := &model.MembersResult{
		Items:     len(googleMembers),
		Resources: syncMembers,
	}
	syncMembersResult.HashCode = hash.Get(syncMembersResult)

	return syncMembersResult, nil
}

func (i *IdentityProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error) {
	syncUsers := make([]model.User, 0)

	for _, member := range mbr.Resources {
		u, err := i.ps.GetUser(ctx, member.ID)
		if err != nil {
			return nil, ErrGettingUser
		}

		e := model.User{
			ID:          u.Id,
			Name:        model.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
			Active:      !u.Suspended,
			Email:       u.PrimaryEmail,
		}
		e.HashCode = hash.Get(e)

		syncUsers = append(syncUsers, e)
	}

	syncUsersResult := &model.UsersResult{
		Items:     len(syncUsers),
		Resources: syncUsers,
	}
	syncUsersResult.HashCode = hash.Get(syncUsersResult)

	return syncUsersResult, nil
}

// GetUsersAndGroupsUsers returns a model.UsersResult and model.GroupsUsersResult data structures with the users and groups
func (i *IdentityProvider) GetUsersAndGroupsUsers(ctx context.Context, groups *model.GroupsResult) (*model.UsersResult, *model.GroupsUsersResult, error) {
	// make pUsers unique
	userSet := make(map[string]struct{})

	pUsers := make([]model.User, 0)
	pGroupsUsers := make([]model.GroupUsers, 0)

	for _, pGroup := range groups.Resources {

		pMembers, err := i.GetGroupMembers(ctx, pGroup.ID)
		if err != nil {
			return nil, nil, err
		}

		pUsersFromMembers, err := i.GetUsersFromGroupMembers(ctx, pMembers)
		if err != nil {
			return nil, nil, err
		}

		for _, pUser := range pUsersFromMembers.Resources {
			if _, ok := userSet[pUser.ID]; !ok {
				pUsers = append(pUsers, pUser)
				userSet[pUser.ID] = struct{}{}
			}
		}

		pGroupUsers := model.GroupUsers{
			Items: len(pMembers.Resources),
			Group: model.Group{
				ID:    pGroup.ID,
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
