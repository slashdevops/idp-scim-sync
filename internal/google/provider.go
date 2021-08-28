package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/core"
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

func NewGoogleIdentityProvider(ds DirectoryService) (core.IdentityProviderService, error) {
	if ds == nil {
		return nil, ErrDirectoryServiceNil
	}

	return &googleProvider{
		ds: ds,
	}, nil
}

func (g *googleProvider) GetGroups(ctx context.Context, filter []string) (*core.GroupsResult, error) {
	syncGroups := make([]*core.Group, 0)

	googleGroups, err := g.ds.ListGroups(filter)
	if err != nil {
		return nil, ErrListingGroups
	}

	for _, grp := range googleGroups {
		syncGroups = append(syncGroups, &core.Group{
			ID:    grp.Id,
			Name:  grp.Name,
			Email: grp.Email,
		})
	}

	// TODO: Check groups are not repeated thanks to the filter

	syncResult := &core.GroupsResult{
		Items:     len(googleGroups),
		Resources: syncGroups,
	}

	return syncResult, nil
}

func (g *googleProvider) GetUsers(ctx context.Context, filter []string) (*core.UsersResult, error) {
	syncUsers := make([]*core.User, 0)

	googleUsers, err := g.ds.ListUsers(filter)
	if err != nil {
		return nil, ErrListingUsers
	}

	for _, usr := range googleUsers {
		syncUsers = append(syncUsers, &core.User{
			ID:          usr.Id,
			Name:        core.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		})
	}

	// TODO: Check users are not repeated thanks to the filter

	uResult := &core.UsersResult{
		Items:     len(googleUsers),
		Resources: syncUsers,
	}

	return uResult, nil
}

func (g *googleProvider) GetGroupMembers(ctx context.Context, id string) (*core.MembersResult, error) {
	syncMembers := make([]*core.Member, 0)

	googleMembers, err := g.ds.ListGroupMembers(id)
	if err != nil {
		return nil, ErrListingGroupMembers
	}

	for _, member := range googleMembers {
		syncMembers = append(syncMembers, &core.Member{
			ID:    member.Id,
			Email: member.Email,
		})
	}

	syncMembersResult := &core.MembersResult{
		Items:     len(googleMembers),
		Resources: syncMembers,
	}

	return syncMembersResult, nil
}

func (g *googleProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *core.MembersResult) (*core.UsersResult, error) {
	syncUsers := make([]*core.User, 0)

	for _, member := range mbr.Resources {
		u, err := g.ds.GetUser(member.ID)
		if err != nil {
			return nil, ErrGettingUser
		}

		syncUsers = append(syncUsers, &core.User{
			ID:          u.Id,
			Name:        core.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
			Active:      !u.Suspended,
			Email:       u.PrimaryEmail,
		})
	}

	syncUsersResult := &core.UsersResult{
		Items:     len(syncUsers),
		Resources: syncUsers,
	}

	return syncUsersResult, nil
}
