package aws

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
)

func Test_NewSecretsManagerService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return SecretsManagerService and no error", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return an error if no client is provided", func(t *testing.T) {
		svc, err := NewSecretsManagerService(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrSMClientNil)
		assert.Nil(t, svc)
	})
}
