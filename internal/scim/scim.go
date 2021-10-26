package scim

import (
	"context"
	"errors"
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
	CreateUser(ctx context.Context, u *aws.CreateUserRequest) (*aws.CreateUserResponse, error)
	DeleteUser(ctx context.Context, id string) error

	ListGroups(ctx context.Context, filter string) (*aws.ListGroupsResponse, error)
	CreateGroup(ctx context.Context, g *aws.CreateGroupRequest) (*aws.CreateGroupResponse, error)
	DeleteGroup(ctx context.Context, id string) error
	PatchGroup(ctx context.Context, pgr *aws.PatchGroupRequest) error
}

var ErrAWSSCIMProviderNil = fmt.Errorf("scim: AWS SCIM Provider is nil")

type SCIMProvider struct {
	scim AWSSCIMProvider
}

func NewSCIMProvider(scim AWSSCIMProvider) (*SCIMProvider, error) {
	if scim == nil {
		return nil, ErrAWSSCIMProviderNil
	}

	return &SCIMProvider{scim: scim}, nil
}

func (s *SCIMProvider) CreateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error) {
	groups := make([]model.Group, 0)

	for _, group := range gr.Resources {

		groupRequest := &aws.CreateGroupRequest{
			DisplayName: group.Name,
		}

		r, err := s.scim.CreateGroup(ctx, groupRequest)
		if err != nil {
			return nil, fmt.Errorf("scim: error creating group: %w", err)
		}

		e := group
		e.SCIMID = r.ID
		e.HashCode = hash.Get(e)
		groups = append(groups, e)
	}

	ret := &model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	ret.HashCode = hash.Get(ret)

	return ret, nil
}

func (s *SCIMProvider) GetGroups(ctx context.Context) (*model.GroupsResult, error) {
	groupsResponse, err := s.scim.ListGroups(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("scim: error listing groups: %w", err)
	}

	groups := make([]model.Group, 0)
	for _, group := range groupsResponse.Resources {
		e := model.Group{
			SCIMID: group.ID,
			Name:   group.DisplayName,
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

func (s *SCIMProvider) UpdateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error) {
	return nil, errors.New("not implemented")
}

func (s *SCIMProvider) DeleteGroups(ctx context.Context, gr *model.GroupsResult) error {
	for _, group := range gr.Resources {
		// TODO: implement a delay to avoid AWS throttling
		if err := s.scim.DeleteGroup(ctx, group.SCIMID); err != nil {
			return fmt.Errorf("scim: error deleting group: %s, %w", group.SCIMID, err)
		}
	}
	return nil
}

func (s *SCIMProvider) GetUsers(ctx context.Context) (*model.UsersResult, error) {
	usersResponse, err := s.scim.ListUsers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("scim: error listing users: %w", err)
	}

	users := make([]model.User, 0)
	for _, user := range usersResponse.Resources {
		e := model.User{
			SCIMID: user.ID,
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

func (s *SCIMProvider) CreateUsers(ctx context.Context, usrs *model.UsersResult) (*model.UsersResult, error) {
	users := make([]model.User, 0)

	for _, user := range usrs.Resources {

		userRequest := &aws.CreateUserRequest{
			DisplayName: user.DisplayName,
			ExternalId:  user.IPID,
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

		r, err := s.scim.CreateUser(ctx, userRequest)
		if err != nil {
			return nil, fmt.Errorf("scim: error creating user: %w", err)
		}

		e := user
		e.SCIMID = r.ID
		e.HashCode = hash.Get(e)
		users = append(users, e)
	}

	ret := &model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	ret.HashCode = hash.Get(ret)

	return ret, nil
}

func (s *SCIMProvider) UpdateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error) {
	return nil, errors.New("not implemented")
}

func (s *SCIMProvider) DeleteUsers(ctx context.Context, ur *model.UsersResult) error {
	for _, user := range ur.Resources {
		// TODO: implement a delay to avoid AWS throttling
		if err := s.scim.DeleteUser(ctx, user.SCIMID); err != nil {
			return fmt.Errorf("scim: error deleting user: %s, %w", user.SCIMID, err)
		}
	}
	return nil
}

func (s *SCIMProvider) CreateGroupsMembers(ctx context.Context, gur *model.GroupsUsersResult) error {
	for _, groupUsers := range gur.Resources {
		usersID := make([]string, 0)
		for _, user := range groupUsers.Resources {
			usersID = append(usersID, user.SCIMID)
		}

		patchGroupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          groupUsers.Group.SCIMID,
				DisplayName: groupUsers.Group.Name,
			},
			Patch: aws.Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []aws.Operation{
					{
						OP:    "add",
						Path:  "members",
						Value: usersID,
					},
				},
			},
		}

		if err := s.scim.PatchGroup(ctx, patchGroupRequest); err != nil {
			return fmt.Errorf("scim: error patching group: %w", err)
		}
	}

	return nil
}

func (s *SCIMProvider) DeleteGroupsMembers(ctx context.Context, gur *model.GroupsUsersResult) error {
	for _, groupUsers := range gur.Resources {
		usersID := make([]string, 0)
		for _, user := range groupUsers.Resources {
			usersID = append(usersID, user.SCIMID)
		}

		patchGroupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          groupUsers.Group.SCIMID,
				DisplayName: groupUsers.Group.Name,
			},
			Patch: aws.Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []aws.Operation{
					{
						OP:    "remove",
						Path:  "members",
						Value: usersID,
					},
				},
			},
		}

		if err := s.scim.PatchGroup(ctx, patchGroupRequest); err != nil {
			return fmt.Errorf("scim: error patching group: %w", err)
		}
	}

	return nil
}

func (s *SCIMProvider) GetUsersAndGroupsUsers(ctx context.Context) (*model.UsersResult, *model.GroupsUsersResult, error) {
	usersResult, err := s.GetUsers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("scim: error getting users: %w", err)
	}
	usersResult.HashCode = hash.Get(usersResult)

	groupsIDUsers := make(map[string][]model.User)
	groupsData := make(map[string]model.Group)

	// inefficient but it is the only way to do that because AWS API Doesn't have efficient
	// way to get the members of groups
	for _, user := range usersResult.Resources {

		// https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
		f := fmt.Sprintf("members eq \"%s\"", user.SCIMID)
		sGroupsResponse, err := s.scim.ListGroups(ctx, f)
		if err != nil {
			return nil, nil, fmt.Errorf("scim: error listing groups: %w", err)
		}

		for _, grp := range sGroupsResponse.Resources {
			groupsIDUsers[grp.ID] = append(groupsIDUsers[grp.ID], user)

			// only one time assignment for users in different groups
			if _, ok := groupsData[grp.ID]; !ok {
				e := model.Group{
					SCIMID: grp.ID,
					Name:   grp.DisplayName,
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
