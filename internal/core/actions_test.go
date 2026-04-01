package core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/core"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func emptyGroupsResult() *model.GroupsResult {
	return model.GroupsResultBuilder().WithResources([]*model.Group{}).Build()
}

func emptyUsersResult() *model.UsersResult {
	return model.UsersResultBuilder().WithResources([]*model.User{}).Build()
}

func emptyGroupsMembersResult() *model.GroupsMembersResult {
	r := &model.GroupsMembersResult{
		Items:     0,
		Resources: make([]*model.GroupMembers, 0),
	}
	r.SetHashCode()
	return r
}

func TestSyncGroupsFromState(t *testing.T) {
	ctx := context.Background()

	t.Run("no changes when hashes match", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		groups := emptyGroupsResult()
		got, err := syncGroupsFromState(ctx, mockSCIM, groups, groups)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, groups.HashCode, got.HashCode)
	})

	t.Run("creates new groups when state is empty", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		idpGroups := model.GroupsResultBuilder().WithResources([]*model.Group{
			model.GroupBuilder().WithIPID("g1").WithName("Group1").WithEmail("g1@test.com").Build(),
		}).Build()
		stateGroups := emptyGroupsResult()

		createdGroup := model.GroupsResultBuilder().WithResources([]*model.Group{
			model.GroupBuilder().WithIPID("g1").WithSCIMID("scim-g1").WithName("Group1").WithEmail("g1@test.com").Build(),
		}).Build()

		mockSCIM.EXPECT().CreateGroups(gomock.Any(), gomock.Any()).Return(createdGroup, nil)
		// DeleteGroups is NOT called because there are no groups to delete (state was empty)

		got, err := syncGroupsFromState(ctx, mockSCIM, idpGroups, stateGroups)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 1, got.Items)
	})

	t.Run("returns error when create fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		idpGroups := model.GroupsResultBuilder().WithResources([]*model.Group{
			model.GroupBuilder().WithIPID("g1").WithName("Group1").WithEmail("g1@test.com").Build(),
		}).Build()
		stateGroups := emptyGroupsResult()

		mockSCIM.EXPECT().CreateGroups(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("create failed"))

		got, err := syncGroupsFromState(ctx, mockSCIM, idpGroups, stateGroups)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "create failed")
	})
}

func TestSyncUsersFromState(t *testing.T) {
	ctx := context.Background()

	t.Run("no changes when hashes match", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		users := emptyUsersResult()
		got, err := syncUsersFromState(ctx, mockSCIM, users, users)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, users.HashCode, got.HashCode)
	})

	t.Run("creates new users when state is empty", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		idpUsers := model.UsersResultBuilder().WithResources([]*model.User{
			model.UserBuilder().WithIPID("u1").WithUserName("user@test.com").WithDisplayName("User").
				WithName(model.NameBuilder().WithGivenName("U").WithFamilyName("Ser").Build()).
				WithActive(true).Build(),
		}).Build()
		stateUsers := emptyUsersResult()

		createdUsers := model.UsersResultBuilder().WithResources([]*model.User{
			model.UserBuilder().WithIPID("u1").WithSCIMID("scim-u1").WithUserName("user@test.com").WithDisplayName("User").
				WithName(model.NameBuilder().WithGivenName("U").WithFamilyName("Ser").Build()).
				WithActive(true).Build(),
		}).Build()

		mockSCIM.EXPECT().CreateUsers(gomock.Any(), gomock.Any()).Return(createdUsers, nil)
		// DeleteUsers is NOT called because there are no users to delete (state was empty)

		got, err := syncUsersFromState(ctx, mockSCIM, idpUsers, stateUsers)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 1, got.Items)
	})

	t.Run("returns error when create fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		idpUsers := model.UsersResultBuilder().WithResources([]*model.User{
			model.UserBuilder().WithIPID("u1").WithUserName("user@test.com").WithDisplayName("User").
				WithName(model.NameBuilder().WithGivenName("U").WithFamilyName("Ser").Build()).
				WithActive(true).Build(),
		}).Build()
		stateUsers := emptyUsersResult()

		mockSCIM.EXPECT().CreateUsers(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("create failed"))

		got, err := syncUsersFromState(ctx, mockSCIM, idpUsers, stateUsers)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSyncGroupsMembersFromState(t *testing.T) {
	ctx := context.Background()

	t.Run("no changes when hashes match", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		members := emptyGroupsMembersResult()
		groups := emptyGroupsResult()
		users := emptyUsersResult()

		got, err := syncGroupsMembersFromState(ctx, mockSCIM, members, members, groups, users)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, members.HashCode, got.HashCode)
	})

	t.Run("creates members when state differs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		group := model.GroupBuilder().WithIPID("g1").WithSCIMID("scim-g1").WithName("Group1").Build()
		member := model.MemberBuilder().WithIPID("u1").WithSCIMID("scim-u1").WithEmail("u@test.com").WithStatus("ACTIVE").Build()

		idpMembers := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				model.GroupMembersBuilder().WithGroup(group).WithResources([]*model.Member{member}).Build(),
			},
		}
		idpMembers.SetHashCode()

		stateMembers := emptyGroupsMembersResult()

		groups := model.GroupsResultBuilder().WithResources([]*model.Group{group}).Build()
		users := model.UsersResultBuilder().WithResources([]*model.User{
			model.UserBuilder().WithIPID("u1").WithSCIMID("scim-u1").WithUserName("u@test.com").WithDisplayName("U").
				WithName(model.NameBuilder().WithGivenName("U").WithFamilyName("T").Build()).WithActive(true).Build(),
		}).Build()

		mockSCIM.EXPECT().CreateGroupsMembers(gomock.Any(), gomock.Any()).Return(
			&model.GroupsMembersResult{Items: 1, Resources: []*model.GroupMembers{
				model.GroupMembersBuilder().WithGroup(group).WithResources([]*model.Member{member}).Build(),
			}}, nil)
		// DeleteGroupsMembers is NOT called because there are no members to delete (state was empty)

		got, err := syncGroupsMembersFromState(ctx, mockSCIM, idpMembers, stateMembers, groups, users)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestStateSync(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error on invalid last sync time", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		state := &model.State{
			LastSync:  "not-a-valid-time",
			Resources: &model.StateResources{},
		}

		gr, ur, gmr, err := stateSync(ctx, state, mockSCIM,
			emptyGroupsResult(), emptyUsersResult(), emptyGroupsMembersResult())
		assert.Error(t, err)
		assert.Nil(t, gr)
		assert.Nil(t, ur)
		assert.Nil(t, gmr)
		assert.Contains(t, err.Error(), "error parsing last sync time")
	})

	t.Run("succeeds with no changes when everything matches", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		groups := emptyGroupsResult()
		users := emptyUsersResult()
		members := emptyGroupsMembersResult()

		state := &model.State{
			LastSync: time.Now().Format(time.RFC3339),
			Resources: &model.StateResources{
				Groups:        groups,
				Users:         users,
				GroupsMembers: members,
			},
		}

		gr, ur, gmr, err := stateSync(ctx, state, mockSCIM, groups, users, members)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
		assert.NotNil(t, ur)
		assert.NotNil(t, gmr)
	})

	t.Run("propagates group sync error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSCIM := mocks.NewMockSCIMService(mockCtrl)

		idpGroups := model.GroupsResultBuilder().WithResources([]*model.Group{
			model.GroupBuilder().WithIPID("g1").WithName("NewGroup").Build(),
		}).Build()

		state := &model.State{
			LastSync: time.Now().Format(time.RFC3339),
			Resources: &model.StateResources{
				Groups:        emptyGroupsResult(),
				Users:         emptyUsersResult(),
				GroupsMembers: emptyGroupsMembersResult(),
			},
		}

		mockSCIM.EXPECT().CreateGroups(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("SCIM unavailable"))

		gr, ur, gmr, err := stateSync(ctx, state, mockSCIM, idpGroups, emptyUsersResult(), emptyGroupsMembersResult())
		assert.Error(t, err)
		assert.Nil(t, gr)
		assert.Nil(t, ur)
		assert.Nil(t, gmr)
		assert.Contains(t, err.Error(), "SCIM unavailable")
	})
}
