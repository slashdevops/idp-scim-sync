package core

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/core"
	"github.com/stretchr/testify/assert"
)

func TestReconcilingGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(nil).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
	})

	t.Run("Should return error when CreateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(nil, errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should return error when UpdateGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(nil, errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should return error when DeleteGroups return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}}}
		update := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}}}
		delete := &model.GroupsResult{Items: 1, Resources: []*model.Group{{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}}}

		mockSCIMService.EXPECT().CreateGroups(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateGroups(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroups(ctx, delete).Return(errors.New("test error")).Times(1)

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		update := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, grc)
		assert.NotNil(t, gru)
	})

	t.Run("Should return error when SCIM service in nil", func(t *testing.T) {
		create := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		update := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, nil, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should error when create groups is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		update := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, nil, update, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should error when update groups is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		delete := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, nil, delete)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})

	t.Run("Should error when delete groups is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
		update := &model.GroupsResult{Items: 0, Resources: []*model.Group{}}

		grc, gru, err := reconcilingGroups(ctx, mockSCIMService, create, update, nil)
		assert.Error(t, err)
		assert.Nil(t, grc)
		assert.Nil(t, gru)
	})
}

func TestReconcilingUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Emails: []model.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}}}
		update := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Emails: []model.Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}}}
		delete := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Emails: []model.Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(nil).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
	})

	t.Run("Should return error when CreateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Emails: []model.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}}}
		update := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Emails: []model.Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}}}
		delete := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Emails: []model.Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(nil, errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when UpdateUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Emails: []model.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}}}
		update := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Emails: []model.Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}}}
		delete := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Emails: []model.Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(nil, errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when DeleteUsers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Emails: []model.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}}}
		update := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Emails: []model.Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}}}
		delete := &model.UsersResult{Items: 1, Resources: []*model.User{{IPID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Emails: []model.Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}}}

		mockSCIMService.EXPECT().CreateUsers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().UpdateUsers(ctx, update).Return(update, nil).Times(1)
		mockSCIMService.EXPECT().DeleteUsers(ctx, delete).Return(errors.New("test error")).Times(1)

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		update := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []*model.User{}}

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, delete)
		assert.NoError(t, err)
		assert.NotNil(t, urc)
		assert.NotNil(t, uru)
	})

	t.Run("Should return error when SCIM service is nil", func(t *testing.T) {
		create := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		update := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []*model.User{}}

		urc, uru, err := reconcilingUsers(ctx, nil, create, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when create users is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		update := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []*model.User{}}

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, nil, update, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when update users is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		delete := &model.UsersResult{Items: 0, Resources: []*model.User{}}

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, nil, delete)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})

	t.Run("Should return error when delete users is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.UsersResult{Items: 0, Resources: []*model.User{}}
		update := &model.UsersResult{Items: 0, Resources: []*model.User{}}

		urc, uru, err := reconcilingUsers(ctx, mockSCIMService, create, update, nil)
		assert.Error(t, err)
		assert.Nil(t, urc)
		assert.Nil(t, uru)
	})
}

func TestReconcilingGroupsMembers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	t.Run("Should call all the methods one time each and no error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*model.Member{{IPID: "1", Email: "user.1@mail.com"}},
				},
			},
		}
		delete := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*model.Member{{IPID: "2", Email: "user.2@mail.com"}},
				},
			},
		}

		mockSCIMService.EXPECT().CreateGroupsMembers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroupsMembers(ctx, delete).Return(nil).Times(1)

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, create, delete)
		assert.NoError(t, err)
		assert.NotNil(t, gmrc)
	})

	t.Run("Should return error when CreateGroupsMembers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*model.Member{{IPID: "1", Email: "user.1@mail.com"}},
				},
			},
		}
		delete := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*model.Member{{IPID: "2", Email: "user.2@mail.com"}},
				},
			},
		}

		mockSCIMService.EXPECT().CreateGroupsMembers(ctx, create).Return(nil, errors.New("test error")).Times(1)

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, create, delete)
		assert.Error(t, err)
		assert.Nil(t, gmrc)
	})

	t.Run("Should return error when DeleteGroupsMembers return error", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*model.Member{{IPID: "1", Email: "user.1@mail.com"}},
				},
			},
		}
		delete := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*model.Member{{IPID: "2", Email: "user.2@mail.com"}},
				},
			},
		}

		mockSCIMService.EXPECT().CreateGroupsMembers(ctx, create).Return(create, nil).Times(1)
		mockSCIMService.EXPECT().DeleteGroupsMembers(ctx, delete).Return(errors.New("test error")).Times(1)

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, create, delete)
		assert.Error(t, err)
		assert.Nil(t, gmrc)
	})

	t.Run("Should call all the methods one time each and no error when resources are empty", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsMembersResult{
			Items:     0,
			Resources: []*model.GroupMembers{},
		}
		delete := &model.GroupsMembersResult{
			Items:     0,
			Resources: []*model.GroupMembers{},
		}

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, create, delete)
		assert.NoError(t, err)
		assert.NotNil(t, gmrc)

		assert.Equal(t, 0, gmrc.Items)
	})

	t.Run("Should return error when SCIM service in nil", func(t *testing.T) {
		create := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*model.Member{{IPID: "1", Email: "user.1@mail.com"}},
				},
			},
		}
		delete := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*model.Member{{IPID: "2", Email: "user.2@mail.com"}},
				},
			},
		}

		gmrc, err := reconcilingGroupsMembers(ctx, nil, create, delete)
		assert.Error(t, err)
		assert.Nil(t, gmrc)
	})

	t.Run("Should error when create groups is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		delete := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*model.Member{{IPID: "2", Email: "user.2@mail.com"}},
				},
			},
		}

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, nil, delete)
		assert.Error(t, err)
		assert.Nil(t, gmrc)
	})

	t.Run("Should error when delete groups is nil", func(t *testing.T) {
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		create := &model.GroupsMembersResult{
			Items: 1,
			Resources: []*model.GroupMembers{
				{
					Items:     1,
					Group:     &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*model.Member{{IPID: "1", Email: "user.1@mail.com"}},
				},
			},
		}

		gmrc, err := reconcilingGroupsMembers(ctx, mockSCIMService, create, nil)
		assert.Error(t, err)
		assert.Nil(t, gmrc)
	})
}
