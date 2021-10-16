package scim

import (
	"context"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// This implement core.SCIMService interface

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/scim/scim_mocks.go -source=scim.go AWSSCIMProvider

// Define AWSSCIMProvider interface to use aws.aws methods
type AWSSCIMProvider interface {
	ListUsers(ctx context.Context, filter string) (*aws.ListUsersResponse, error)
	ListGroups(ctx context.Context, filter string) (*aws.ListGroupsResponse, error)

	CreateUser(ctx context.Context, u *aws.CreateUserRequest) (*aws.CreateUserResponse, error)
	CreateGroup(ctx context.Context, g *aws.CreateGroupRequest) (*aws.CreateGroupResponse, error)
}

var ErrAWSSCIMProviderNil = fmt.Errorf("AWSSCIMProvider is nil")

type SCIMProvider struct {
	scim AWSSCIMProvider
}

func NewSCIMProvider(scim AWSSCIMProvider) (*SCIMProvider, error) {
	if scim == nil {
		return nil, ErrAWSSCIMProviderNil
	}

	return &SCIMProvider{scim: scim}, nil
}

func (s *SCIMProvider) GetGroups(ctx context.Context) (*model.GroupsResult, error) {
	sGroupsResponse, err := s.scim.ListGroups(ctx, "")
	if err != nil {
		return nil, err
	}

	groups := make([]model.Group, 0)
	for _, group := range sGroupsResponse.Resources {
		e := model.Group{
			ID:   group.ID,
			Name: group.DisplayName,
		}
		e.HashCode = hash.Get(e)

		groups = append(groups, e)
	}

	groupsResult := &model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	groupsResult.HashCode = hash.Get(groupsResult)

	return groupsResult, nil
}

func (s *SCIMProvider) GetUsers(ctx context.Context) (*model.UsersResult, error) {
	UsersResponse, err := s.scim.ListUsers(ctx, "")
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	for _, user := range UsersResponse.Resources {
		e := model.User{
			ID: user.ID,
			Name: model.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			DisplayName: user.DisplayName,
			Active:      user.Active,
			Email:       user.Emails[0].Value,
		}
		e.HashCode = hash.Get(e)

		users = append(users, e)
	}

	usersResult := &model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	usersResult.HashCode = hash.Get(usersResult)

	return usersResult, nil
}

//
func (s *SCIMProvider) GetUsersAndGroupsUsers(ctx context.Context) (*model.UsersResult, *model.GroupsUsersResult, error) {
	usersResult, err := s.GetUsers(ctx)
	if err != nil {
		return nil, nil, err
	}
	usersResult.HashCode = hash.Get(usersResult)

	groupsIDUsers := make(map[string][]model.User)
	groupsData := make(map[string]model.Group)

	// inefficient but it is the only way to do that because AWS API Doesn't have efficient
	// way to get the members of groups
	for _, user := range usersResult.Resources {

		// https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
		f := fmt.Sprintf("members eq \"%s\"", user.ID)
		sGroupsResponse, err := s.scim.ListGroups(ctx, f)
		if err != nil {
			return nil, nil, err
		}

		for _, grp := range sGroupsResponse.Resources {
			groupsIDUsers[grp.ID] = append(groupsIDUsers[grp.ID], user)

			// only one time assignment
			if _, ok := groupsData[grp.ID]; !ok {
				e := model.Group{
					ID:   grp.ID,
					Name: grp.DisplayName,
				}
				e.HashCode = hash.Get(e)

				groupsData[grp.ID] = e
			}
		}
	}

	groupsUsers := make([]model.GroupUsers, 0)

	for groupID, users := range groupsIDUsers {
		e := model.GroupUsers{
			Items:     len(users),
			Group:     groupsData[groupID],
			Resources: users,
		}
		e.HashCode = hash.Get(e)

		groupsUsers = append(groupsUsers, e)
	}

	groupsUsersResult := &model.GroupsUsersResult{
		Items:     len(groupsUsers),
		Resources: groupsUsers,
	}
	groupsUsersResult.HashCode = hash.Get(groupsUsersResult)

	return usersResult, groupsUsersResult, nil
}

func (s *SCIMProvider) CreateGroups(ctx context.Context, gr *model.GroupsResult) error {
	for _, group := range gr.Resources {

		sGroupRequest := &aws.CreateGroupRequest{
			DisplayName: group.Name,
		}

		_, err := s.scim.CreateGroup(ctx, sGroupRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SCIMProvider) CreateUsers(ctx context.Context, usrs *model.UsersResult) error {
	for _, user := range usrs.Resources {

		sUserRequest := &aws.CreateUserRequest{
			DisplayName: user.DisplayName,
			ExternalId:  user.ID,
			Name: aws.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			Emails: []aws.Email{
				{
					Value: user.Email,
					Type:  "work",
				},
			},
			Active: user.Active,
		}

		_, err := s.scim.CreateUser(ctx, sUserRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SCIMProvider) CreateMembers(ctx context.Context, ur *model.GroupsUsersResult) error {
	return nil
}

func (s *SCIMProvider) UpdateGroups(ctx context.Context, gr *model.GroupsResult) error {
	return nil
}

func (s *SCIMProvider) UpdateUsers(ctx context.Context, ur *model.UsersResult) error {
	return nil
}

func (s *SCIMProvider) DeleteGroups(ctx context.Context, gr *model.GroupsResult) error {
	return nil
}

func (s *SCIMProvider) DeleteUsers(ctx context.Context, ur *model.UsersResult) error {
	return nil
}

func (s *SCIMProvider) DeleteMembers(ctx context.Context, ur *model.GroupsUsersResult) error {
	return nil
}
