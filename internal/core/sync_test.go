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

func Test_syncService_getIdentityProviderData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		filters := []string{""}
		grpr := &model.GroupsResult{
			Items:     1,
			HashCode:  "",
			Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}},
		}
		usrs := &model.UsersResult{
			Items:     1,
			HashCode:  "",
			Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}},
		}
		grpsUsrs := &model.GroupsUsersResult{
			Items:     1,
			HashCode:  "",
			Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}},
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

		assert.Equal(t, users.Resources[0].IPID, "1")
		assert.Equal(t, users.Resources[0].Name.GivenName, "user")
		assert.Equal(t, users.Resources[0].Name.FamilyName, "1")
		assert.Equal(t, users.Resources[0].Email, "user.1@mail.com")

		assert.Equal(t, groups.Resources[0].IPID, "1")
		assert.Equal(t, groups.Resources[0].Name, "group 1")
		assert.Equal(t, groups.Resources[0].Email, "group.1@mail.com")

		assert.Equal(t, groupsUsers.Resources[0].Group.IPID, "1")
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
			Items:     1,
			HashCode:  "",
			Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}},
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
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{
			Items:     1,
			HashCode:  "",
			Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}},
		}
		usrs := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		grpsUsrs := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}

		mockSCIMService.EXPECT().GetGroups(ctx).Return(grpr, nil).Times(1)
		mockSCIMService.EXPECT().GetUsersAndGroupsUsers(ctx).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockSCIMService)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)

		assert.Equal(t, len(users.Resources), 1)
		assert.Equal(t, len(groups.Resources), 1)
		assert.Equal(t, len(groupsUsers.Resources), 1)

		assert.Equal(t, users.Resources[0].IPID, "1")
		assert.Equal(t, users.Resources[0].Name.GivenName, "user")
		assert.Equal(t, users.Resources[0].Name.FamilyName, "1")
		assert.Equal(t, users.Resources[0].Email, "user.1@mail.com")

		assert.Equal(t, groups.Resources[0].IPID, "1")
		assert.Equal(t, groups.Resources[0].Name, "group 1")
		assert.Equal(t, groups.Resources[0].Email, "group.1@mail.com")

		assert.Equal(t, groupsUsers.Resources[0].Group.IPID, "1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Name, "group 1")
		assert.Equal(t, groupsUsers.Resources[0].Group.Email, "group.1@mail.com")
	})

	t.Run("Should return error when GetGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		mockSCIMService.EXPECT().GetGroups(ctx).Return(nil, errors.New("test error")).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockSCIMService)
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, groups)
		assert.Nil(t, groupsUsers)
	})

	t.Run("Should return valid values", func(t *testing.T) {
		mockProviderService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{
			Items:     1,
			HashCode:  "",
			Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}},
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
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		grpr := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		usrs := &model.UsersResult{Items: 0, Resources: []model.User{}}
		grpsUsrs := &model.GroupsUsersResult{Items: 0, Resources: []model.GroupUsers{}}

		mockSCIMService.EXPECT().GetGroups(ctx).Return(grpr, nil).Times(1)
		mockSCIMService.EXPECT().GetUsersAndGroupsUsers(ctx).Return(usrs, grpsUsrs, nil).Times(1)

		users, groups, groupsUsers, err := getSCIMData(ctx, mockSCIMService)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotNil(t, groups)
		assert.NotNil(t, groupsUsers)
	})
}

func Test_syncService_reconcilingSCIMGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(delete, nil).Times(1)

		grc, gru, grd, err := reconcilingSCIMGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
		assert.NotNil(t, grd)
	})

	t.Run("Should return error when CreateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(nil, errors.New("test error")).Times(1)

		grc, gru, grd, err := reconcilingSCIMGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
		assert.Nil(t, grd)
	})

	t.Run("Should return error when UpdateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(nil, errors.New("test error")).Times(1)

		grc, gru, grd, err := reconcilingSCIMGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
		assert.Nil(t, grd)
	})

	t.Run("Should return error when DeleteGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(nil, errors.New("test error")).Times(1)

		grc, gru, grd, err := reconcilingSCIMGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
		assert.Nil(t, grd)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		update := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []model.Group{}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(delete, nil).Times(1)

		grc, gru, grd, err := reconcilingSCIMGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
		assert.NotNil(t, grd)
	})
}

func Test_syncService_reconcilingSCIMUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(delete, nil).Times(1)

		urc, uru, urd, err := reconcilingSCIMUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
		assert.NotNil(t, urd)
	})

	t.Run("Should return error when CreateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(nil, errors.New("test error")).Times(1)

		urc, uru, urd, err := reconcilingSCIMUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
		assert.Nil(t, urd)
	})

	t.Run("Should return error when UpdateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(nil, errors.New("test error")).Times(1)

		urc, uru, urd, err := reconcilingSCIMUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
		assert.Nil(t, urd)
	})

	t.Run("Should return error when DeleteUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(nil, errors.New("test error")).Times(1)

		urc, uru, urd, err := reconcilingSCIMUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
		assert.Nil(t, urd)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 0, Resources: []model.User{}}
		update := &model.UsersResult{Items: 0, Resources: []model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []model.User{}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(delete, nil).Times(1)

		urc, uru, urd, err := reconcilingSCIMUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
		assert.NotNil(t, urd)
	})
}

func Test_syncService_reconcilingSCIMGroupsUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}
		delete := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}

		mockSCIMService.EXPECT().CreateMembers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().DeleteMembers(ctx, delete).Return(delete, nil).Times(1)

		gurc, gurd, err := reconcilingSCIMGroupsUsers(ctx, mockSCIMService, create, delete)
		assert.NoError(t, err)
		assert.NotNil(t, gurc)
		assert.NotNil(t, gurd)
	})

	t.Run("Should return error when CreateMembers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}
		delete := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}

		mockSCIMService.EXPECT().CreateMembers(ctx, create).Return(nil, errors.New("test error")).Times(1)

		gurc, gurd, err := reconcilingSCIMGroupsUsers(ctx, mockSCIMService, create, delete)
		assert.Error(t, err)
		assert.Nil(t, gurc)
		assert.Nil(t, gurd)
	})

	t.Run("Should return error when DeleteMembers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}
		delete := &model.GroupsUsersResult{Items: 1, Resources: []model.GroupUsers{{Items: 1, Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "123456789"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}}}

		mockSCIMService.EXPECT().CreateMembers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().DeleteMembers(ctx, delete).Return(delete, errors.New("test error")).Times(1)

		gurc, gurd, err := reconcilingSCIMGroupsUsers(ctx, mockSCIMService, create, delete)
		assert.Error(t, err)
		assert.Nil(t, gurc)
		assert.Nil(t, gurd)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsUsersResult{Items: 0, Resources: []model.GroupUsers{}}
		delete := &model.GroupsUsersResult{Items: 0, Resources: []model.GroupUsers{}}

		mockSCIMService.EXPECT().CreateMembers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().DeleteMembers(ctx, delete).Return(delete, nil).Times(1)

		gurc, gurd, err := reconcilingSCIMGroupsUsers(ctx, mockSCIMService, create, delete)
		assert.NoError(t, err)
		assert.NotNil(t, gurc)
		assert.NotNil(t, gurd)
	})
}
