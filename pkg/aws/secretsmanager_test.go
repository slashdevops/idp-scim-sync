package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewSecretsManagerService(t *testing.T) {
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
		assert.ErrorIs(t, err, ErrSecretManagerClientNil)
		assert.Nil(t, svc)
	})
}

func TestSecretsManager_GetSecretValue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return valid value when key exist", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)
		ctx := context.TODO()

		SMOut := &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("testValue"),
		}

		SMIn := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String("testKey"),
			VersionStage: aws.String("AWSCURRENT"),
		}

		mockSMClientAPI.EXPECT().GetSecretValue(ctx, SMIn).Times(1).Return(SMOut, nil)

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		value, err := svc.GetSecretValue(ctx, "testKey")
		assert.NoError(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, "testValue", value)
	})

	t.Run("Should return a error and empty valid value when key doesn't exist", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)
		ctx := context.TODO()

		mockSMClientAPI.EXPECT().GetSecretValue(ctx, gomock.Any()).Times(1).Return(nil, errors.New("test error"))

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		value, err := svc.GetSecretValue(ctx, "testKeya")
		assert.Error(t, err)
		assert.Empty(t, value)
	})

	t.Run("Should return error decoding bin value", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)
		ctx := context.TODO()

		SMOut := &secretsmanager.GetSecretValueOutput{
			SecretBinary: []byte("testValue"),
		}

		mockSMClientAPI.EXPECT().GetSecretValue(ctx, gomock.Any()).Times(1).Return(SMOut, nil)

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		value, err := svc.GetSecretValue(ctx, "testValue")
		assert.Error(t, err)
		assert.Empty(t, value)
	})

	t.Run("Should return valid value when key is bin", func(t *testing.T) {
		mockSMClientAPI := mocks.NewMockSecretsManagerClientAPI(mockCtrl)
		ctx := context.TODO()
		key := "MyTestKey"

		SMIn := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(key),
			VersionStage: aws.String("AWSCURRENT"),
		}

		message := base64.StdEncoding.EncodeToString([]byte(key))

		SMOut := &secretsmanager.GetSecretValueOutput{
			SecretBinary: []byte(message),
		}

		mockSMClientAPI.EXPECT().GetSecretValue(ctx, SMIn).Times(1).Return(SMOut, nil)

		svc, err := NewSecretsManagerService(mockSMClientAPI)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		value, err := svc.GetSecretValue(ctx, key)
		assert.NoError(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, key, value)
	})
}
