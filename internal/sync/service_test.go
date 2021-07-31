package sync

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSyncService(t *testing.T) {

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
