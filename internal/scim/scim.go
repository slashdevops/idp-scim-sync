package scim

import (
	"context"
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"

	log "github.com/sirupsen/logrus"
)

// This implement core.SCIMService interface

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/scim/scim_mocks.go -source=scim.go AWSSCIMProvider

// AWSSCIMProvider interface to consume aws.aws methods
type AWSSCIMProvider interface {
	// ListUsers lists users in SCIM Provider
	ListUsers(ctx context.Context, filter string) (*aws.ListUsersResponse, error)

	// CreateUser creates a user in SCIM Provider
	CreateUser(ctx context.Context, u *aws.CreateUserRequest) (*aws.CreateUserResponse, error)

	// PutUser updates a user in SCIM Provider
	PutUser(ctx context.Context, usr *aws.PutUserRequest) (*aws.PutUserResponse, error)

	// DeleteUser deletes a user in SCIM Provider
	DeleteUser(ctx context.Context, id string) error

	// GetUser gets a user in SCIM Provider
	GetUser(ctx context.Context, userID string) (*aws.GetUserResponse, error)

	// GetUserByUserName gets a user in SCIM Provider
	GetUserByUserName(ctx context.Context, userName string) (*aws.GetUserResponse, error)

	// ListGroups lists groups in SCIM Provider
	ListGroups(ctx context.Context, filter string) (*aws.ListGroupsResponse, error)

	// CreateGroup creates a group in SCIM Provider
	CreateGroup(ctx context.Context, g *aws.CreateGroupRequest) (*aws.CreateGroupResponse, error)

	// DeleteGroup deletes a group in SCIM Provider
	DeleteGroup(ctx context.Context, id string) error

	// PatchGroup patches a group in SCIM Provider
	PatchGroup(ctx context.Context, pgr *aws.PatchGroupRequest) error
}

// ErrSCIMProviderNil is returned when the SCIMProvider is nil
var ErrSCIMProviderNil = fmt.Errorf("scim: SCIM Provider is nil")

// SCIMProvider represents a SCIM provider
type SCIMProvider struct {
	scim AWSSCIMProvider
}

// NewSCIMProvider creates a new SCIM provider
func NewSCIMProvider(scim AWSSCIMProvider) (*SCIMProvider, error) {
	if scim == nil {
		return nil, ErrSCIMProviderNil
	}

	return &SCIMProvider{scim: scim}, nil
}

// GetGroups returns groups from SCIM Provider
func (s *SCIMProvider) GetGroups(ctx context.Context) (*model.GroupsResult, error) {
	groupsResponse, err := s.scim.ListGroups(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("scim: error listing groups: %w", err)
	}

	groups := make([]*model.Group, 0)
	for _, group := range groupsResponse.Resources {
		e := &model.Group{
			SCIMID: group.ID,
			Name:   group.DisplayName,
			IPID:   group.ExternalID,
		}
		e.SetHashCode()

		groups = append(groups, e)
	}

	groupsResult := &model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	groupsResult.SetHashCode()

	return groupsResult, nil
}

// CreateGroups creates groups in SCIM Provider
func (s *SCIMProvider) CreateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error) {
	groups := make([]*model.Group, 0)

	for _, group := range gr.Resources {
		groupRequest := &aws.CreateGroupRequest{
			DisplayName: group.Name,
			ExternalID:  group.IPID,
		}

		log.WithFields(log.Fields{
			"group": group.Name,
			"idpid": group.IPID,
			"email": group.Email,
		}).Trace("creating group (details)")

		log.WithFields(log.Fields{
			"group": group.Name,
		}).Warn("creating group")

		r, err := s.scim.CreateGroup(ctx, groupRequest)
		if err != nil {
			return nil, fmt.Errorf("scim: error creating group: %w", err)
		}

		e := group
		e.SCIMID = r.ID
		e.SetHashCode()
		groups = append(groups, e)
	}

	ret := &model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	ret.SetHashCode()

	return ret, nil
}

// UpdateGroups updates groups in SCIM Provider
func (s *SCIMProvider) UpdateGroups(ctx context.Context, gr *model.GroupsResult) (*model.GroupsResult, error) {
	groups := make([]*model.Group, 0)

	for _, group := range gr.Resources {
		groupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          group.SCIMID,
				DisplayName: group.Name,
			},
			Patch: aws.PatchGroup{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []aws.OperationGroup{
					{
						OP: "replace",
						Value: map[string]string{
							"id":         group.SCIMID,
							"ExternalID": group.IPID,
						},
					},
				},
			},
		}

		log.WithFields(log.Fields{
			"group":  group.Name,
			"idpid":  group.IPID,
			"scimid": group.SCIMID,
			"email":  group.Email,
		}).Trace("updating group (details)")

		log.WithFields(log.Fields{
			"group": group.Name,
			"email": group.Email,
		}).Warn("updating group")

		if err := s.scim.PatchGroup(ctx, groupRequest); err != nil {
			return nil, fmt.Errorf("scim: error updating groups: %w", err)
		}

		// return the same group
		e := &model.Group{
			SCIMID: group.SCIMID,
			Name:   group.Name,
			IPID:   group.IPID,
			Email:  group.Email,
		}
		e.SetHashCode()
		groups = append(groups, e)
	}

	ret := &model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	ret.SetHashCode()

	return ret, nil
}

// DeleteGroups deletes groups in SCIM Provider
func (s *SCIMProvider) DeleteGroups(ctx context.Context, gr *model.GroupsResult) error {
	for _, group := range gr.Resources {
		log.WithFields(log.Fields{
			"group":  group.Name,
			"idpid":  group.IPID,
			"scimid": group.SCIMID,
			"email":  group.Email,
		}).Trace("deleting group (details)")

		log.WithFields(log.Fields{
			"group": group.Name,
			"email": group.Email,
		}).Trace("deleting group")

		if err := s.scim.DeleteGroup(ctx, group.SCIMID); err != nil {
			return fmt.Errorf("scim: error deleting group: %s, %w", group.SCIMID, err)
		}
	}
	return nil
}

// GetUsers returns users from SCIM Provider
func (s *SCIMProvider) GetUsers(ctx context.Context) (*model.UsersResult, error) {
	usersResponse, err := s.scim.ListUsers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("scim: error listing users: %w", err)
	}

	users := make([]model.User, 0)
	for _, user := range usersResponse.Resources {
		e := model.User{
			IPID:   user.ExternalID,
			SCIMID: user.ID,
			Name: model.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			DisplayName: user.DisplayName,
			Active:      user.Active,
			Email:       user.Emails[0].Value,
		}
		e.SetHashCode()

		users = append(users, e)
	}

	usersResult := &model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	usersResult.SetHashCode()

	return usersResult, nil
}

// CreateUsers creates users in SCIM Provider
func (s *SCIMProvider) CreateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error) {
	users := make([]model.User, 0)

	for _, user := range ur.Resources {
		userRequest := &aws.CreateUserRequest{
			ID:          "",
			UserName:    user.Email,
			DisplayName: user.DisplayName,
			ExternalID:  user.IPID,
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

		log.WithFields(log.Fields{
			"user":  user.DisplayName,
			"email": user.Email,
			"ipdid": user.IPID,
		}).Trace("creating user")

		r, err := s.scim.CreateUser(ctx, userRequest)
		if err != nil {
			return nil, fmt.Errorf("scim: error creating user: %w", err)
		}

		e := model.User{
			IPID:   user.IPID,
			SCIMID: r.ID,
			Name: model.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			DisplayName: user.DisplayName,
			Active:      user.Active,
			Email:       user.Email,
		}
		e.SCIMID = r.ID
		e.SetHashCode()

		users = append(users, e)
	}

	ret := &model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	ret.SetHashCode()

	return ret, nil
}

// UpdateUsers updates users in SCIM Provider given a list of users
func (s *SCIMProvider) UpdateUsers(ctx context.Context, ur *model.UsersResult) (*model.UsersResult, error) {
	users := make([]model.User, 0)

	for _, user := range ur.Resources {
		userRequest := &aws.PutUserRequest{
			ID:          user.SCIMID,
			DisplayName: user.DisplayName,
			UserName:    user.Email,
			ExternalID:  user.IPID,
			Name: aws.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			Emails: []aws.Email{
				{
					Value:   user.Email,
					Type:    "work",
					Primary: true,
				},
			},
			Active: user.Active,
		}

		log.WithFields(log.Fields{
			"user":   user.DisplayName,
			"email":  user.Email,
			"ipdid":  user.IPID,
			"scimid": user.SCIMID,
		}).Trace("updating user (details)")

		log.Tracef("scim: updating user -> userRequest: %s", utils.ToJSON(userRequest))

		log.WithFields(log.Fields{
			"user":  user.DisplayName,
			"email": user.Email,
		}).Debug("updating user")

		r, err := s.scim.PutUser(ctx, userRequest)
		if err != nil {
			return nil, fmt.Errorf("scim: error updating user: %w", err)
		}

		e := model.User{
			IPID:   user.IPID,
			SCIMID: r.ID,
			Name: model.Name{
				FamilyName: user.Name.FamilyName,
				GivenName:  user.Name.GivenName,
			},
			DisplayName: user.DisplayName,
			Active:      user.Active,
			Email:       user.Email,
		}
		e.SCIMID = r.ID
		e.SetHashCode()
		users = append(users, e)
	}

	ret := &model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	ret.SetHashCode()

	return ret, nil
}

// DeleteUsers deletes users in SCIM Provider given a list of users
func (s *SCIMProvider) DeleteUsers(ctx context.Context, ur *model.UsersResult) error {
	for _, user := range ur.Resources {
		log.WithFields(log.Fields{
			"user":   user.DisplayName,
			"email":  user.Email,
			"scimid": user.SCIMID,
			"idpid":  user.IPID,
		}).Trace("deleting user (details)")

		log.WithFields(log.Fields{
			"user":  user.DisplayName,
			"email": user.Email,
		}).Warn("deleting user")

		if err := s.scim.DeleteUser(ctx, user.SCIMID); err != nil {
			return fmt.Errorf("scim: error deleting user: %s, %w", user.SCIMID, err)
		}
	}
	return nil
}

// CreateGroupsMembers creates groups members in SCIM Provider given a list of groups members
func (s *SCIMProvider) CreateGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) (*model.GroupsMembersResult, error) {
	groupsMembers := make([]*model.GroupMembers, 0)

	for _, groupMembers := range gmr.Resources {
		members := make([]model.Member, 0)

		// https://talks.golang.org/2012/10things.slide#2
		membersIDValue := []struct {
			Value string `json:"value"`
		}{}

		for _, member := range groupMembers.Resources {
			if member.SCIMID == "" {
				u, err := s.scim.GetUserByUserName(ctx, member.Email)
				if err != nil {
					return nil, fmt.Errorf("scim: error getting user by email: %w", err)
				}
				member.SCIMID = u.ID
			}

			membersIDValue = append(membersIDValue, struct {
				Value string `json:"value"`
			}{
				Value: member.SCIMID,
			})

			e := model.Member{
				IPID:   member.IPID,
				SCIMID: member.SCIMID,
				Email:  member.Email,
			}
			e.SetHashCode()
			members = append(members, e)

			log.WithFields(log.Fields{
				"group":  groupMembers.Group.Name,
				"idpid":  member.IPID,
				"scimid": member.SCIMID,
				"email":  member.Email,
			}).Trace("adding member to group (details)")

			log.WithFields(log.Fields{
				"group": groupMembers.Group.Name,
				"email": member.Email,
			}).Warn("adding member to group")
		}

		e := groupMembers
		e.SetHashCode()
		e.Resources = members

		groupsMembers = append(groupsMembers, e)

		patchGroupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          groupMembers.Group.SCIMID,
				DisplayName: groupMembers.Group.Name,
			},
			Patch: aws.PatchGroup{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []aws.OperationGroup{
					{
						OP:    "add",
						Path:  "members",
						Value: membersIDValue,
					},
				},
			},
		}

		if err := s.scim.PatchGroup(ctx, patchGroupRequest); err != nil {
			return nil, fmt.Errorf("scim: error patching group: %w", err)
		}
	}

	ret := &model.GroupsMembersResult{
		Items:     len(groupsMembers),
		Resources: groupsMembers,
	}
	ret.SetHashCode()

	return ret, nil
}

// DeleteGroupsMembers deletes groups members in SCIM Provider given a list of groups members
func (s *SCIMProvider) DeleteGroupsMembers(ctx context.Context, gmr *model.GroupsMembersResult) error {
	for _, groupMembers := range gmr.Resources {
		// https://talks.golang.org/2012/10things.slide#2
		membersIDValue := []struct {
			Value string `json:"value"`
		}{}

		for _, member := range groupMembers.Resources {
			membersIDValue = append(membersIDValue, struct {
				Value string `json:"value"`
			}{
				Value: member.SCIMID,
			})

			log.WithFields(log.Fields{
				"group":  groupMembers.Group.Name,
				"idpid":  member.IPID,
				"scimid": member.SCIMID,
				"email":  member.Email,
			}).Trace("removing member from group (details)")

			log.WithFields(log.Fields{
				"group": groupMembers.Group.Name,
				"email": member.Email,
			}).Warn("removing member from group")
		}

		patchGroupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          groupMembers.Group.SCIMID,
				DisplayName: groupMembers.Group.Name,
			},
			Patch: aws.PatchGroup{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []aws.OperationGroup{
					{
						OP:    "remove",
						Path:  "members",
						Value: membersIDValue,
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

// GetGroupsMembers returns a list of groups and their members from the SCIM Provider
// NOTE: this method doesn't work becuae unfortunately the SCIM API doesn't support
// list the members of a group, or get a group and their members at the same time
// reference: https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
func (s *SCIMProvider) GetGroupsMembers(ctx context.Context, gr *model.GroupsResult) (*model.GroupsMembersResult, error) {
	groupMembers := make([]*model.GroupMembers, 0)

	for _, group := range gr.Resources {
		// https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
		f := fmt.Sprintf("displayName eq %q", group.Name)
		lgr, err := s.scim.ListGroups(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("scim: error listing groups: %w", err)
		}
		log.Tracef("lgr: lgr : %s", utils.ToJSON(lgr))

		for _, gr := range lgr.Resources {
			members := make([]model.Member, 0)

			for _, member := range gr.Members {
				u, err := s.scim.GetUser(ctx, member.Value)
				if err != nil {
					return nil, fmt.Errorf("scim: error getting user: %s, error %w", member.Value, err)
				}

				m := model.Member{
					SCIMID: member.Value,
					Email:  u.Emails[0].Value,
				}
				m.SetHashCode()

				members = append(members, m)
			}

			e := &model.GroupMembers{
				Items:     len(members),
				Group:     *group,
				Resources: members,
			}
			e.SetHashCode()

			groupMembers = append(groupMembers, e)
		}
	}

	groupsMembersResult := &model.GroupsMembersResult{
		Items:     len(groupMembers),
		Resources: groupMembers,
	}
	groupsMembersResult.SetHashCode()

	return groupsMembersResult, nil
}

// GetGroupsMembersBruteForce returns a list of groups and their members from the SCIM Provider
// NOTE: this is an bad alternative to the method GetGroupsMembers,  because read the note in the method.
func (s *SCIMProvider) GetGroupsMembersBruteForce(ctx context.Context, gr *model.GroupsResult, ur *model.UsersResult) (*model.GroupsMembersResult, error) {
	groupMembers := make([]*model.GroupMembers, 0)

	// brute force implemented here thanks to the fxxckin' aws sso scim api
	for _, group := range gr.Resources {
		members := make([]model.Member, 0)

		for _, user := range ur.Resources {
			// https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
			f := fmt.Sprintf("id eq %q and members eq %q", group.SCIMID, user.SCIMID)
			lgr, err := s.scim.ListGroups(ctx, f)
			if err != nil {
				return nil, fmt.Errorf("scim: error listing groups: %w", err)
			}

			if lgr.TotalResults > 0 { // crazy thing of the AWS SSO SCIM API, it doesn't return the memnber into the Resources array
				m := model.Member{
					IPID:   user.IPID,
					SCIMID: user.SCIMID,
					Email:  user.Email,
				}
				m.SetHashCode()

				members = append(members, m)
			}
			e := &model.GroupMembers{
				Items:     len(members),
				Group:     *group,
				Resources: members,
			}
			e.SetHashCode()

			groupMembers = append(groupMembers, e)
		}
	}

	groupsMembersResult := &model.GroupsMembersResult{
		Items:     len(groupMembers),
		Resources: groupMembers,
	}
	groupsMembersResult.SetHashCode()

	return groupsMembersResult, nil
}
