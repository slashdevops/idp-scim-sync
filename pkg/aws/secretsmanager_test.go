package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
)

func Test_NewSecretsManagerService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return SecretsManagerService", func(t *testing.T) {
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

func Test_SecretsManager_GetSecretValue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return valid value when key exist", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)

		SMOut := &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("testValue"),
		}

		SMIn := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("testKey"),
			VersionStage: aws.String("AWSCURRENT"),
		}

		mockSMClientAPI.EXPECT().GetSecretValue(context.TODO(), SMIn).Return(SMOut, nil)

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		value, err := svc.GetSecretValue(context.TODO(), "testKey")
		assert.NoError(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, "testValue", value)
	})

	// t.Run("Should return a error and empty valid value when key doesn't exist", func(t *testing.T) {
	// 	mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)

	// 	SMIn := &secretsmanager.GetSecretValueInput{
	// 		SecretId:     aws.String("testKey"),
	// 		VersionStage: aws.String("AWSCURRENT"),
	// 	}

	// 	mockSMClientAPI.EXPECT().GetSecretValue(context.TODO(), SMIn).Return(nil, errors.New("test error"))

	// 	svc, err := NewSecretsManagerService(mockSMClientAPI)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, svc)

	// 	value, err := svc.GetSecretValue(context.TODO(), "testKeya")
	// 	assert.Error(t, err)
	// 	assert.Nil(t, value)
	// })
}
