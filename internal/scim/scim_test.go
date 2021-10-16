package scim

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/scim"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/stretchr/testify/assert"
)

func Test_NewSCIMProvider(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return SCIMProvider and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		svc, err := NewSCIMProvider(mockSCIM)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return an error if no AWSSCIMProvider is provided", func(t *testing.T) {
		svc, err := NewSCIMProvider(nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.ErrorIs(t, err, ErrAWSSCIMProviderNil)
	})
}

func Test_SCIMProvider_GetGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return a error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		mockSCIM.EXPECT().ListGroups(context.TODO(), gomock.Any()).Return(nil, errors.New("test error"))

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetGroups(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, gr)
	})

	t.Run("Should return a empty list of groups and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		groups := &aws.ListsGroupsResponse{}
		mockSCIM.EXPECT().ListGroups(context.TODO(), gomock.Any()).Return(groups, nil)

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetGroups(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, gr)

		assert.Equal(t, 0, len(gr.Resources))
		assert.Equal(t, 0, gr.Items)
	})

	t.Run("Should return a list of groups and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		groups := &aws.ListsGroupsResponse{
			GeneralResponse: aws.GeneralResponse{
				TotalResults: 2,
				ItemsPerPage: 2,
				StartIndex:   0,
				Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
			},
			Resources: []*aws.Group{
				{
					ID:          "1",
					DisplayName: "group 1",
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					Members:     []aws.Memeber{},
					Meta:        aws.Meta{ResourceType: "Group", Created: "2020-04-01T12:00:00Z", LastModified: "2020-04-01T12:00:00Z"},
				},
				{
					ID:          "2",
					DisplayName: "group 2",
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					Members:     []aws.Memeber{},
					Meta:        aws.Meta{ResourceType: "Group", Created: "2020-04-02T12:00:00Z", LastModified: "2020-04-02T12:00:00Z"},
				},
			},
		}

		mockSCIM.EXPECT().ListGroups(context.TODO(), gomock.Any()).Return(groups, nil)

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetGroups(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, gr)
		assert.Equal(t, 2, len(gr.Resources))
		assert.Equal(t, 2, gr.Items)
		// assert.Equal(t, "", gr.HashCode) //TODO: create a object to compare this

		assert.Equal(t, "1", gr.Resources[0].ID)
		assert.Equal(t, "2", gr.Resources[1].ID)

		assert.Equal(t, "group 1", gr.Resources[0].Name)
		assert.Equal(t, "group 1", gr.Resources[0].Name)

		assert.Equal(t, "", gr.Resources[0].Email)
		assert.Equal(t, "", gr.Resources[1].Email)

		// assert.Equal(t, "", gr.Resources[0].HashCode) //TODO: create a object to compare this
		// assert.Equal(t, "", gr.Resources[1].HashCode) //TODO: create a object to compare this
	})
}
