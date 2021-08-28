package core

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_syncService_NewSyncService(t *testing.T) {
	t.Run("New Service with parameters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockRepository := mocks.NewMockSyncRepository(mockCtrl)

		svc, err := NewSyncService(ctx, mockProviderService, mockSCIMService, mockRepository)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("New Service without parameters", func(t *testing.T) {
		ctx := context.TODO()

		svc, err := NewSyncService(ctx, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

func Test_syncService_SyncGroupsAndUsers(t *testing.T) {
	t.Run("Sync Groups and Users", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockRepository := mocks.NewMockSyncRepository(mockCtrl)

		mockProviderService.EXPECT().GetGroups(ctx, gomock.Any()).Return(nil, nil)
		mockProviderService.EXPECT().GetUsers(ctx, gomock.Any()).Return(nil, nil)
		mockSCIMService.EXPECT().CreateOrUpdateGroups(ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().CreateOrUpdateUsers(ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().DeleteGroups(ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().DeleteUsers(ctx, gomock.Any()).Return(nil)

		svc, err := NewSyncService(ctx, mockProviderService, mockSCIMService, mockRepository)
		assert.NoError(t, err)

		err = svc.SyncGroupsAndUsers()
		assert.NoError(t, err)
	})
}
