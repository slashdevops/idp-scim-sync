package sync

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_syncService_NewSyncService(t *testing.T) {

	t.Run("New Service with parameters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := NewMockProviderService(mockCtrl)
		mockSCIMService := NewMockSCIMService(mockCtrl)

		svc, err := NewSyncService(&ctx, mockProviderService, mockSCIMService)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("New Service without parameters", func(t *testing.T) {
		ctx := context.TODO()

		svc, err := NewSyncService(&ctx, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

func Test_syncService_SyncGroupsAndUsers(t *testing.T) {

	t.Run("Sync Groups and Users", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx := context.TODO()
		mockProviderService := NewMockProviderService(mockCtrl)
		mockSCIMService := NewMockSCIMService(mockCtrl)

		// urFilter := "name=aaa*"
		// grFilter := "name=aaa*"
		grResult := &GroupResult{
			Total:    10,
			Items:    1,
			NextItem: 2,
			Resources: []*Group{
				{
					Name:  "group1",
					Email: "group1@gmail.com",
				},
			},
		}
		urResult := &UserResult{
			Total:    10,
			Items:    1,
			NextItem: 2,
			Resources: []*User{
				{
					Name:  "user1",
					Email: "user1@gmail.com",
				},
			},
		}

		mockProviderService.EXPECT().GetGroups(&ctx, gomock.Any()).Return(grResult, nil)
		mockProviderService.EXPECT().GetUsers(&ctx, gomock.Any()).Return(urResult, nil)
		mockSCIMService.EXPECT().CreateOrUpdateGroups(&ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().CreateOrUpdateUsers(&ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().DeleteGroups(&ctx, gomock.Any()).Return(nil)
		mockSCIMService.EXPECT().DeleteUsers(&ctx, gomock.Any()).Return(nil)

		svc, err := NewSyncService(&ctx, mockProviderService, mockSCIMService)
		assert.NoError(t, err)

		err = svc.SyncGroupsAndUsers()
		assert.NoError(t, err)

	})

}
