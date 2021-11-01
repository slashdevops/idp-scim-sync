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

func TestSyncService_NewSyncService(t *testing.T) {
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
		assert.Equal(t, err, ErrIdentiyProviderServiceNil)
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
		assert.Equal(t, err, ErrStateRepositoryNil)
	})
}

func TestSyncService_reconcilingGroups(t *testing.T) {
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
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(nil).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
	})

	t.Run("Should return error when CreateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(nil, errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should return error when UpdateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(nil, errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should return error when DeleteGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		update := &model.GroupsResult{Items: 0, Resources: []model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
	})
}

func TestSyncService_reconcilingUsers(t *testing.T) {
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
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(nil).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
	})

	t.Run("Should return error when CreateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(nil, errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when UpdateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(nil, errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when DeleteUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com"}}}
		update := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com"}}}
		delete := &model.UsersResult{Items: 1, Resources: []model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 0, Resources: []model.User{}}
		update := &model.UsersResult{Items: 0, Resources: []model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []model.User{}}

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
	})
}
