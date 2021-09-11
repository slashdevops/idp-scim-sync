package core

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/internal/mocks/core"
	"github.com/stretchr/testify/assert"
)

func Test_syncService_NewSyncService(t *testing.T) {
	t.Run("New Service with parameters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockRepository := mocks.NewMockRepository(mockCtrl)

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

	t.Run("New Service without context, return specific error", func(t *testing.T) {
		svc, err := NewSyncService(nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrNilContext)
	})

	t.Run("New Service without IdentityProviderServoce, return specific error", func(t *testing.T) {
		ctx := context.TODO()
		svc, err := NewSyncService(ctx, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrProviderServiceNil)
	})

	t.Run("New Service without SCIMServoce, return context specific error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		svc, err := NewSyncService(ctx, mockProviderService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrSCIMServiceNil)
	})

	t.Run("New Service without Repository, return context specific error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		svc, err := NewSyncService(ctx, mockProviderService, mockSCIMService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrRepositoryNil)
	})
}

func Test_syncService_SyncGroupsAndUsers(t *testing.T) {
	t.Run("Sync Groups and Users", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockRepository := mocks.NewMockRepository(mockCtrl)

		mockProviderService.EXPECT().GetGroups(ctx, gomock.Any()).Return(nil, nil)
		mockProviderService.EXPECT().GetUsers(ctx, gomock.Any()).Return(nil, nil)
		mockSCIMService.EXPECT().DeleteGroups(ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().DeleteUsers(ctx, gomock.Any()).Return(nil)

		svc, err := NewSyncService(ctx, mockProviderService, mockSCIMService, mockRepository)
		assert.NoError(t, err)

		err = svc.SyncGroupsAndUsers()
		assert.NoError(t, err)
	})
}
