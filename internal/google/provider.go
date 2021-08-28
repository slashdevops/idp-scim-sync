package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/provider"
)

var (
	ErrDirectoryServiceNil = errors.New("directory service is nil")
	ErrListingGroups       = errors.New("error listing groups")
	ErrListingUsers        = errors.New("error listing users")
	ErrListingGroupMembers = errors.New("error listing group members")
	ErrGettingUser         = errors.New("error getting user")
)

type googleProvider struct {
	ds DirectoryService
}

func NewGoogleIdentityProvider(ds DirectoryService) (provider.IdentityProviderService, error) {
	if ds == nil {
		return nil, ErrDirectoryServiceNil
	}

	return &googleProvider{
		ds: ds,
	}, nil
}

func (g *googleProvider) GetGroups(ctx context.Context, filter []string) (*model.GroupsResult, error) {
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

	// TODO: Check groups are not repeated thanks to the filter

	syncResult := &model.GroupsResult{
		Items:     len(googleGroups),
		Resources: syncGroups,
	}

	return syncResult, nil
}

func (g *googleProvider) GetUsers(ctx context.Context, filter []string) (*model.UsersResult, error) {
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

	// TODO: Check users are not repeated thanks to the filter

	uResult := &model.UsersResult{
		Items:     len(googleUsers),
		Resources: syncUsers,
	}

	return uResult, nil
}

func (g *googleProvider) GetGroupMembers(ctx context.Context, id string) (*model.MembersResult, error) {
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

func (g *googleProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *model.MembersResult) (*model.UsersResult, error) {
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
