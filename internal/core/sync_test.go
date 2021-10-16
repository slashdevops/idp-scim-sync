package core

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/core"
	"github.com/stretchr/testify/assert"
)

func Test_syncService_NewSyncService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("New Service with parameters", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockStateRepository := mocks.NewMockStateRepository(mockCtrl)

		svc, err := NewSyncService(context.TODO(), mockProviderService, mockSCIMService, mockStateRepository)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("New Service without parameters", func(t *testing.T) {
		svc, err := NewSyncService(context.TODO(), nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("New Service without IdentityProviderServoce, return specific error", func(t *testing.T) {
		svc, err := NewSyncService(context.TODO(), nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrProviderServiceNil)
	})

	t.Run("New Service without SCIMServoce, return context specific error", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		svc, err := NewSyncService(context.TODO(), mockProviderService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrSCIMServiceNil)
	})

	t.Run("New Service without Repository, return context specific error", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		svc, err := NewSyncService(context.TODO(), mockProviderService, mockSCIMService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrRepositoryNil)
	})
}

func Test_syncService_SyncGroupsAndUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Sync Groups and Users", func(t *testing.T) {
		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockStateRepository := mocks.NewMockStateRepository(mockCtrl)

		mockProviderService.EXPECT().GetGroups(ctx, gomock.Any()).Return(nil, nil).Times(1)
		mockProviderService.EXPECT().GetUsers(ctx, gomock.Any()).Return(nil, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, gomock.Any()).Return(nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, gomock.Any()).Return(nil).Times(1)

		svc, err := NewSyncService(ctx, mockProviderService, mockSCIMService, mockStateRepository)
		assert.NoError(t, err)

		err = svc.SyncGroupsAndUsers()
		assert.NoError(t, err)
	})
}

func Test_syncService_getIdentityProviderData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		filters := []string{""}
		grpr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					ID:       "1",
					Name:     "group 1",
					Email:    "group.1@mail.com",
					HashCode: "123456789",
				},
			},
		}
		usrs := &model.UsersResult{
			Items: 1,
			Resources: []model.User{
				{
					ID:    "1",
					Name:  model.Name{GivenName: "user", FamilyName: "1"},
					Email: "user.1@mail.com",
				},
			},
		}
		grpsUsrs := &model.GroupsUsersResult{
			Items: 1,
			Resources: []model.GroupUsers{
				{
					Items: 1,
					Group: model.Group{
						ID:       "1",
						Name:     "group 1",
						Email:    "group.1@mail.com",
						HashCode: "123456789",
					},
					Resources: []model.User{
						{
							ID:    "1",
							Name:  model.Name{GivenName: "user", FamilyName: "1"},
							Email: "user.1@mail.com",
						},
					},
				},
			},
		}

		mockProviderService.EXPECT().GetGroups(ctx, filters).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx, grpr).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getIdentityProviderData(ctx, mockProviderService, filters)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)

		assert.Equal(t, len(users.Resources), 1)
		assert.Equal(t, len(groups.Resources), 1)
		assert.Equal(t, len(groupsUsers.Resources), 1)

		assert.Equal(t, users.Resources[0].ID, "1")
		assert.Equal(t, users.Resources[0].Name.GivenName, "user")
		assert.Equal(t, users.Resources[0].Name.FamilyName, "1")
		assert.Equal(t, users.Resources[0].Email, "user.1@mail.com")

		assert.Equal(t, groups.Resources[0].ID, "1")
		assert.Equal(t, groups.Resources[0].Name, "group 1")
		assert.Equal(t, groups.Resources[0].Email, "group.1@mail.com")

		assert.Equal(t, groupsUsers.Resources[0].Group.ID, "1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Name, "group 1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Email, "group.1@mail.com")
	})

	t.Run("Should return error when GetGroups return error", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		filters := []string{""}

		mockProviderService.EXPECT().GetGroups(ctx, filters).Return(nil, errors.New("test error")).Times(1)

		users, groups, groupsUsers, err := getIdentityProviderData(ctx, mockProviderService, filters)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, groups)
		assert.Nil(t, groupsUsers)
	})

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		filters := []string{""}
		grpr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					ID:       "1",
					Name:     "group 1",
					Email:    "group.1@mail.com",
					HashCode: "123456789",
				},
			},
		}

		mockProviderService.EXPECT().GetGroups(ctx, filters).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx, grpr).Return(nil, nil, errors.New("test error")).Times(1)

		users, groups, groupsUsers, err := getIdentityProviderData(ctx, mockProviderService, filters)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, groups)
		assert.Nil(t, groupsUsers)
	})

	t.Run("Should return empty UsersResult, GroupsResult and GroupsUsersResult", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		filters := []string{""}
		grpr := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		usrs := &model.UsersResult{Items: 0, Resources: []model.User{}}
		grpsUsrs := &model.GroupsUsersResult{Items: 0, Resources: []model.GroupUsers{}}

		mockProviderService.EXPECT().GetGroups(ctx, filters).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx, grpr).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getIdentityProviderData(ctx, mockProviderService, filters)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)
	})
}

func Test_syncService_getSCIMData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					ID:       "1",
					Name:     "group 1",
					Email:    "group.1@mail.com",
					HashCode: "123456789",
				},
			},
		}
		usrs := &model.UsersResult{
			Items: 1,
			Resources: []model.User{
				{
					ID:    "1",
					Name:  model.Name{GivenName: "user", FamilyName: "1"},
					Email: "user.1@mail.com",
				},
			},
		}
		grpsUsrs := &model.GroupsUsersResult{
			Items: 1,
			Resources: []model.GroupUsers{
				{
					Items: 1,
					Group: model.Group{
						ID:       "1",
						Name:     "group 1",
						Email:    "group.1@mail.com",
						HashCode: "123456789",
					},
					Resources: []model.User{
						{
							ID:    "1",
							Name:  model.Name{GivenName: "user", FamilyName: "1"},
							Email: "user.1@mail.com",
						},
					},
				},
			},
		}

		mockProviderService.EXPECT().GetGroups(ctx).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockProviderService)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)

		assert.Equal(t, len(users.Resources), 1)
		assert.Equal(t, len(groups.Resources), 1)
		assert.Equal(t, len(groupsUsers.Resources), 1)

		assert.Equal(t, users.Resources[0].ID, "1")
		assert.Equal(t, users.Resources[0].Name.GivenName, "user")
		assert.Equal(t, users.Resources[0].Name.FamilyName, "1")
		assert.Equal(t, users.Resources[0].Email, "user.1@mail.com")

		assert.Equal(t, groups.Resources[0].ID, "1")
		assert.Equal(t, groups.Resources[0].Name, "group 1")
		assert.Equal(t, groups.Resources[0].Email, "group.1@mail.com")

		assert.Equal(t, groupsUsers.Resources[0].Group.ID, "1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Name, "group 1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Email, "group.1@mail.com")
	})

	t.Run("Should return error when GetGroups return error", func(t *testing.T) {
		mockProviderService := mocks.NewMockSCIMService(mockCtrl)

		mockProviderService.EXPECT().GetGroups(ctx).Return(nil, errors.New("test error")).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockProviderService)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, groups)
		assert.Nil(t, groupsUsers)
	})

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					ID:       "1",
					Name:     "group 1",
					Email:    "group.1@mail.com",
					HashCode: "123456789",
				},
			},
		}

		mockProviderService.EXPECT().GetGroups(ctx).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx).Return(nil, nil, errors.New("test error")).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockProviderService)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, groups)
		assert.Nil(t, groupsUsers)
	})

	t.Run("Should return empty UsersResult, GroupsResult and GroupsUsersResult", func(t *testing.T) {
		mockProviderService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		usrs := &model.UsersResult{Items: 0, Resources: []model.User{}}
		grpsUsrs := &model.GroupsUsersResult{Items: 0, Resources: []model.GroupUsers{}}

		mockProviderService.EXPECT().GetGroups(ctx).Return(grpr, nil).Times(1)
		mockProviderService.EXPECT().GetUsersAndGroupsUsers(ctx).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockProviderService)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)
	})
}
