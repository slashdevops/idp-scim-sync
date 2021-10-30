package scim

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/scim"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewSCIMProvider(t *testing.T) {
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
		assert.ErrorIs(t, err, ErrSCIMProviderNil)
	})
}

func TestSCIMProvider_GetGroups(t *testing.T) {
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
		groups := &aws.ListGroupsResponse{}
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
		groups := &aws.ListGroupsResponse{
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
					Members:     []aws.Member{},
					Meta:        aws.Meta{ResourceType: "Group", Created: "2020-04-01T12:00:00Z", LastModified: "2020-04-01T12:00:00Z"},
				},
				{
					ID:          "2",
					DisplayName: "group 2",
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					Members:     []aws.Member{},
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

		assert.Equal(t, "1", gr.Resources[0].SCIMID)
		assert.Equal(t, "2", gr.Resources[1].SCIMID)

		assert.Equal(t, "group 1", gr.Resources[0].Name)
		assert.Equal(t, "group 2", gr.Resources[1].Name)

		assert.Equal(t, "", gr.Resources[0].Email)
		assert.Equal(t, "", gr.Resources[1].Email)

		// assert.Equal(t, "", gr.Resources[0].HashCode) //TODO: create a object to compare this
		// assert.Equal(t, "", gr.Resources[1].HashCode) //TODO: create a object to compare this
	})

	// TODO: test with a list of groups passing filter
}

func TestSCIMProvider_GetUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return a error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		mockSCIM.EXPECT().ListUsers(context.TODO(), gomock.Any()).Return(nil, errors.New("test error"))

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetUsers(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, gr)
	})

	t.Run("Should return a empty list of users and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		users := &aws.ListUsersResponse{}
		mockSCIM.EXPECT().ListUsers(context.TODO(), gomock.Any()).Return(users, nil)

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetUsers(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, gr)

		assert.Equal(t, 0, len(gr.Resources))
		assert.Equal(t, 0, gr.Items)
	})

	t.Run("Should return a list of users and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		users := &aws.ListUsersResponse{
			GeneralResponse: aws.GeneralResponse{
				TotalResults: 2,
				ItemsPerPage: 2,
				StartIndex:   0,
				Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
			},
			Resources: []*aws.User{
				{
					ID:          "1",
					ExternalId:  "1",
					Name:        aws.Name{FamilyName: "1", GivenName: "user"},
					DisplayName: "group 1",
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
					Meta:        aws.Meta{ResourceType: "User", Created: "2020-04-01T12:00:00Z", LastModified: "2020-04-01T12:00:00Z"},
					Emails:      []aws.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
				},
				{
					ID:          "2",
					ExternalId:  "2",
					Name:        aws.Name{FamilyName: "2", GivenName: "user"},
					DisplayName: "group 2",
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
					Meta:        aws.Meta{ResourceType: "User", Created: "2020-04-02T12:00:00Z", LastModified: "2020-04-02T12:00:00Z"},
					Emails:      []aws.Email{{Value: "user.2@mail.com", Type: "work", Primary: true}},
				},
			},
		}

		mockSCIM.EXPECT().ListUsers(context.TODO(), gomock.Any()).Return(users, nil)

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.GetUsers(context.TODO())

		assert.NoError(t, err)
		assert.NotNil(t, gr)
		assert.Equal(t, 2, len(gr.Resources))
		assert.Equal(t, 2, gr.Items)
		// assert.Equal(t, "", gr.HashCode) //TODO: create a object to compare this

		assert.Equal(t, "1", gr.Resources[0].SCIMID)
		assert.Equal(t, "2", gr.Resources[1].SCIMID)

		assert.Equal(t, "1", gr.Resources[0].Name.FamilyName)
		assert.Equal(t, "user", gr.Resources[0].Name.GivenName)

		assert.Equal(t, "2", gr.Resources[1].Name.FamilyName)
		assert.Equal(t, "user", gr.Resources[1].Name.GivenName)

		assert.Equal(t, "user.1@mail.com", gr.Resources[0].Email)
		assert.Equal(t, "user.2@mail.com", gr.Resources[1].Email)

		// assert.Equal(t, "", gr.Resources[0].HashCode) //TODO: create a object to compare this
		// assert.Equal(t, "", gr.Resources[1].HashCode) //TODO: create a object to compare this
	})
}

func TestSCIMProvider_CreateGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should do nothing with empty GroupsResult", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		empty := &model.GroupsResult{}

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.CreateGroups(context.TODO(), empty)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
	})

	t.Run("Should call CreateGroup 1 time and no return error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cgr := &aws.CreateGroupRequest{
			DisplayName: "group 1",
		}
		resp := &aws.CreateGroupResponse{}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateGroup(ctx, cgr).Return(resp, nil).Times(1)

		gr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					IPID:  "1",
					Name:  "group 1",
					Email: "group.1@mail.com",
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.CreateGroups(ctx, gr)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
	})

	t.Run("Should call CreateGroup 1 time and return error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cgr := &aws.CreateGroupRequest{
			DisplayName: "group 1",
		}
		resp := &aws.CreateGroupResponse{}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateGroup(ctx, cgr).Return(resp, errors.New("test error")).Times(1)

		gr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					IPID:  "1",
					Name:  "group 1",
					Email: "group.1@mail.com",
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.CreateGroups(ctx, gr)
		assert.Error(t, err)
		assert.Nil(t, gr)
	})

	t.Run("Should call CreateGroup 2 time and no error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cgr1 := &aws.CreateGroupRequest{
			DisplayName: "group 1",
		}
		cgr2 := &aws.CreateGroupRequest{
			DisplayName: "group 2",
		}
		resp1 := &aws.CreateGroupResponse{
			ID:          "11",
			DisplayName: "group 1",
			Meta: aws.Meta{
				ResourceType: "Group",
				Created:      "2020-04-02T12:00:00Z",
				LastModified: "2020-04-02T12:00:00Z",
			},
			Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
		}
		resp2 := &aws.CreateGroupResponse{
			ID:          "22",
			DisplayName: "group 2",
			Meta: aws.Meta{
				ResourceType: "Group",
				Created:      "2020-04-03T12:00:00Z",
				LastModified: "2020-04-03T12:00:00Z",
			},
			Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
		}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateGroup(ctx, cgr1).Return(resp1, nil).Times(1)
		mockSCIM.EXPECT().CreateGroup(ctx, cgr2).Return(resp2, nil).Times(1)

		gr := &model.GroupsResult{
			Items: 1,
			Resources: []model.Group{
				{
					IPID:  "1",
					Name:  "group 1",
					Email: "group.1@mail.com",
				},
				{
					IPID:  "2",
					Name:  "group 2",
					Email: "group.2@mail.com",
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		gr, err := svc.CreateGroups(ctx, gr)
		assert.NoError(t, err)
		assert.NotNil(t, gr)

		assert.Equal(t, "1", gr.Resources[0].IPID)
		assert.Equal(t, "2", gr.Resources[1].IPID)

		assert.Equal(t, "11", gr.Resources[0].SCIMID)
		assert.Equal(t, "22", gr.Resources[1].SCIMID)

		assert.Equal(t, "group 1", gr.Resources[0].Name)
		assert.Equal(t, "group 2", gr.Resources[1].Name)
	})
}

func TestSCIMProvider_CreateUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should do nothing with empty GroupsResult", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		empty := &model.UsersResult{}

		svc, _ := NewSCIMProvider(mockSCIM)
		cur, err := svc.CreateUsers(context.TODO(), empty)

		assert.NoError(t, err)
		assert.NotNil(t, cur)
	})

	t.Run("Should call CreateUser 1 time and no return error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cur := &aws.CreateUserRequest{
			DisplayName: "user 1",
			ExternalId:  "1",
			Name:        aws.Name{FamilyName: "1", GivenName: "user"},
			Emails: []aws.Email{
				{Value: "user.1@mail.com", Type: "work"},
			},
			Active: true,
		}
		resp := &aws.CreateUserResponse{}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateUser(ctx, cur).Return(resp, nil).Times(1)

		usr := &model.UsersResult{
			Items: 1,
			Resources: []model.User{
				{
					IPID:        "1",
					Name:        model.Name{FamilyName: "1", GivenName: "user"},
					DisplayName: "user 1",
					Email:       "user.1@mail.com",
					Active:      true,
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		ur, err := svc.CreateUsers(ctx, usr)

		assert.NoError(t, err)
		assert.NotNil(t, ur)
	})

	t.Run("Should call CreateUser 1 time and return error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cur := &aws.CreateUserRequest{
			DisplayName: "user 1",
			ExternalId:  "1",
			Name:        aws.Name{FamilyName: "1", GivenName: "user"},
			Emails: []aws.Email{
				{Value: "user.1@mail.com", Type: "work"},
			},
			Active: false,
		}
		resp := &aws.CreateUserResponse{}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateUser(ctx, cur).Return(resp, errors.New("test error")).Times(1)

		usr := &model.UsersResult{
			Items: 1,
			Resources: []model.User{
				{
					IPID:        "1",
					Name:        model.Name{FamilyName: "1", GivenName: "user"},
					DisplayName: "user 1",
					Email:       "user.1@mail.com",
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		ur, err := svc.CreateUsers(ctx, usr)

		assert.Error(t, err)
		assert.Nil(t, ur)
	})

	t.Run("Should call CreateUser 2 time and no return error", func(t *testing.T) {
		mockSCIM := mocks.NewMockAWSSCIMProvider(mockCtrl)
		cur1 := &aws.CreateUserRequest{
			DisplayName: "user 1",
			ExternalId:  "1",
			Name:        aws.Name{FamilyName: "1", GivenName: "user"},
			Emails: []aws.Email{
				{Value: "user.1@mail.com", Type: "work"},
			},
			Active: true,
		}
		cur2 := &aws.CreateUserRequest{
			DisplayName: "user 2",
			ExternalId:  "2",
			Name:        aws.Name{FamilyName: "2", GivenName: "user"},
			Emails: []aws.Email{
				{Value: "user.2@mail.com", Type: "work"},
			},
			Active: true,
		}
		resp1 := &aws.CreateUserResponse{
			ID:         "11",
			ExternalId: "1",
			Name: aws.Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Active:      true,
			Emails:      []aws.Email{{Value: "user.1@mail.com", Type: "work"}},
		}
		resp2 := &aws.CreateUserResponse{
			ID:         "22",
			ExternalId: "2",
			Name: aws.Name{
				FamilyName: "2",
				GivenName:  "user",
			},
			DisplayName: "user 2",
			Active:      true,
			Emails:      []aws.Email{{Value: "user.2@mail.com", Type: "work"}},
		}
		ctx := context.TODO()

		mockSCIM.EXPECT().CreateUser(ctx, cur1).Return(resp1, nil).Times(1)
		mockSCIM.EXPECT().CreateUser(ctx, cur2).Return(resp2, nil).Times(1)

		usr := &model.UsersResult{
			Items: 1,
			Resources: []model.User{
				{
					IPID:        "1",
					Name:        model.Name{FamilyName: "1", GivenName: "user"},
					DisplayName: "user 1",
					Email:       "user.1@mail.com",
					Active:      true,
				},
				{
					IPID:        "2",
					Name:        model.Name{FamilyName: "2", GivenName: "user"},
					DisplayName: "user 2",
					Email:       "user.2@mail.com",
					Active:      true,
				},
			},
		}

		svc, _ := NewSCIMProvider(mockSCIM)
		ur, err := svc.CreateUsers(ctx, usr)

		assert.NoError(t, err)
		assert.NotNil(t, ur)

		assert.Equal(t, "1", ur.Resources[0].IPID)
		assert.Equal(t, "2", ur.Resources[1].IPID)

		assert.Equal(t, "11", ur.Resources[0].SCIMID)
		assert.Equal(t, "22", ur.Resources[1].SCIMID)

		assert.Equal(t, "1", ur.Resources[0].Name.FamilyName)
		assert.Equal(t, "user", ur.Resources[0].Name.GivenName)

		assert.Equal(t, "2", ur.Resources[1].Name.FamilyName)
		assert.Equal(t, "user", ur.Resources[1].Name.GivenName)
	})
}
