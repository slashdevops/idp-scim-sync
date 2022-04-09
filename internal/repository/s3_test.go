package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewS3Repository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return S3Repository and no error", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return an error if no client is provided", func(t *testing.T) {
		svc, err := NewS3Repository(nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrS3ClientNil)
		assert.Nil(t, svc)
	})

	t.Run("Should return an error if no opts WithBucket is provided", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		svc, err := NewS3Repository(mockS3Repository)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOptionWithBucketNil)
		assert.Nil(t, svc)
	})

	t.Run("Should return an error if no opts WithKey is provided", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"))
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOptionWithKeyNil)
		assert.Nil(t, svc)
	})
}

func TestGetState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return valid state and no error", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		sObj := model.State{
			SchemaVersion: "1.0.0",
			CodeVersion:   "0.0.1",
			LastSync:      "2020-01-01T00:00:00Z",
			HashCode:      "123456789",
			Resources: model.StateResources{
				Groups:        model.GroupsResult{},
				Users:         model.UsersResult{},
				GroupsMembers: model.GroupsMembersResult{},
			},
		}

		sObjBytes, err := sObj.MarshalJSON()
		assert.NoError(t, err)
		assert.NotNil(t, sObjBytes)

		S3Out := &s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewBuffer(sObjBytes)),
		}

		mockS3Repository.EXPECT().GetObject(context.TODO(), gomock.Any()).Return(S3Out, nil)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		state, err := svc.GetState(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, state)

		assert.Equal(t, "1.0.0", state.SchemaVersion)
		assert.Equal(t, "0.0.1", state.CodeVersion)
		assert.Equal(t, "2020-01-01T00:00:00Z", state.LastSync)
		assert.Equal(t, "123456789", state.HashCode)
	})

	t.Run("Should return error decoding state", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		sObjBytes := []byte("invalid json")

		S3Out := &s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewBuffer(sObjBytes)),
		}

		mockS3Repository.EXPECT().GetObject(context.TODO(), gomock.Any()).Return(S3Out, nil)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		state, err := svc.GetState(context.TODO())
		assert.Error(t, err)
		assert.Nil(t, state)
	})

	t.Run("Should not return valid state and error instead", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		mockS3Repository.EXPECT().GetObject(context.TODO(), gomock.Any()).Return(nil, errors.New("error"))

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		state, err := svc.GetState(context.TODO())
		assert.Error(t, err)
		assert.Nil(t, state)
	})
}

func TestS3SetState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return error with nil state", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		err = svc.SetState(context.TODO(), nil)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrStateNil)
	})

	t.Run("Should work fine with valid state", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		sObj := &model.State{
			SchemaVersion: "1.0.0",
			CodeVersion:   "0.0.1",
			LastSync:      "2020-01-01T00:00:00Z",
			HashCode:      "123456789",
			Resources: model.StateResources{
				Groups:        model.GroupsResult{},
				Users:         model.UsersResult{},
				GroupsMembers: model.GroupsMembersResult{},
			},
		}

		mockS3Repository.EXPECT().PutObject(context.TODO(), gomock.Any()).Return(nil, nil)

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		err = svc.SetState(context.TODO(), sObj)
		assert.NoError(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {
		mockS3Repository := mocks.NewMockS3ClientAPI(mockCtrl)

		sObj := &model.State{}

		mockS3Repository.EXPECT().PutObject(context.TODO(), gomock.Any()).Return(nil, errors.New("error"))

		svc, err := NewS3Repository(mockS3Repository, WithBucket("MyBucket"), WithKey("MyKey"))
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		err = svc.SetState(context.TODO(), sObj)
		assert.Error(t, err)
	})
}
