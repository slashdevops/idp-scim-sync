package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/slashdevops/aws-sso-gws-sync/internal/sync"
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

func NewGoogleIdentityProvider(ds DirectoryService) (sync.IdentityProviderService, error) {
	if ds == nil {
		return nil, ErrDirectoryServiceNil
	}

	return &googleProvider{
		ds: ds,
	}, nil
}

func (g *googleProvider) GetGroups(ctx context.Context, filter []string) (*sync.GroupsResult, error) {
	syncGroups := make([]*sync.Group, 0)

	googleGroups, err := g.ds.ListGroups(filter)
	if err != nil {
		return nil, ErrListingGroups
	}

	for _, grp := range googleGroups {
		syncGroups = append(syncGroups, &sync.Group{
			Id:    sync.Id{IdentityProvider: grp.Id},
			Name:  grp.Name,
			Email: grp.Email,
		})
	}

	syncResult := &sync.GroupsResult{
		Items:     len(googleGroups),
		Resources: syncGroups,
	}

	return syncResult, nil
}

func (g *googleProvider) GetUsers(ctx context.Context, filter []string) (*sync.UsersResult, error) {
	syncUsers := make([]*sync.User, 0)

	googleUsers, err := g.ds.ListUsers(filter)
	if err != nil {
		return nil, ErrListingUsers
	}

	for _, usr := range googleUsers {
		syncUsers = append(syncUsers, &sync.User{
			Id:          sync.Id{IdentityProvider: usr.Id},
			Name:        sync.Name{FamilyName: usr.Name.FamilyName, GivenName: usr.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName),
			Active:      !usr.Suspended,
			Email:       usr.PrimaryEmail,
		})
	}

	uResult := &sync.UsersResult{
		Items:     len(googleUsers),
		Resources: syncUsers,
	}

	return uResult, nil
}

func (g *googleProvider) GetGroupMembers(ctx context.Context, id string) (*sync.MembersResult, error) {
	syncMembers := make([]*sync.Member, 0)

	googleMembers, err := g.ds.ListGroupMembers(id)
	if err != nil {
		return nil, ErrListingGroupMembers
	}

	for _, member := range googleMembers {
		syncMembers = append(syncMembers, &sync.Member{
			Id:    sync.Id{IdentityProvider: member.Id},
			Email: member.Email,
		})
	}

	syncMembersResult := &sync.MembersResult{
		Items:     len(googleMembers),
		Resources: syncMembers,
	}

	return syncMembersResult, nil
}

func (g *googleProvider) GetUsersFromGroupMembers(ctx context.Context, mbr *sync.MembersResult) (*sync.UsersResult, error) {
	syncUsers := make([]*sync.User, 0)

	for _, member := range mbr.Resources {
		u, err := g.ds.GetUser(member.Id.IdentityProvider)
		if err != nil {
			return nil, ErrGettingUser
		}

		syncUsers = append(syncUsers, &sync.User{
			Id:          sync.Id{IdentityProvider: u.Id},
			Name:        sync.Name{FamilyName: u.Name.FamilyName, GivenName: u.Name.GivenName},
			DisplayName: fmt.Sprintf("%s %s", u.Name.GivenName, u.Name.FamilyName),
			Active:      !u.Suspended,
			Email:       u.PrimaryEmail,
		})
	}

	syncUsersResult := &sync.UsersResult{
		Items:     len(syncUsers),
		Resources: syncUsers,
	}

	return syncUsersResult, nil
}
