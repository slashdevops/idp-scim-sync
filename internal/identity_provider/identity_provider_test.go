package provider

import (
	"context"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/mocks"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
)

func TestNewGoogleIdentityProvider(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return IdentityServiceProvider and no error", func(t *testing.T) {
		mockDS := mocks.NewMockGoogleProviderService(mockCtrl)
		svc, err := NewIdentityProvider(mockDS)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return an error if no DirectoryService is provided", func(t *testing.T) {
		svc, err := NewIdentityProvider(nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

func Test_GoogleProvider_GetGroups(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx    context.Context
		filter []string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *model.GroupsResult
		wantErr bool
	}{
		{
			name: "Should return GroupsResult and no error",
			prepare: func(f *fields) {
				googleGroups := make([]*admin.Group, 0)
				googleGroups = append(googleGroups, &admin.Group{Email: "group1@mail.com", Id: "1", Name: "group1"})
				googleGroups = append(googleGroups, &admin.Group{Email: "group2@mail.com", Id: "2", Name: "group2"})

				f.ds.EXPECT().ListGroups(gomock.Eq([]string{""})).Return(googleGroups, nil).Times(1)
			},
			args: args{ctx: context.TODO(), filter: []string{""}},
			want: &model.GroupsResult{
				Items: 2,
				Resources: []*model.Group{
					{ID: "1", Name: "group1", Email: "group1@mail.com", HashCode: hash.Get(model.Group{ID: "1", Name: "group1", Email: "group1@mail.com"})},
					{ID: "2", Name: "group2", Email: "group2@mail.com", HashCode: hash.Get(model.Group{ID: "2", Name: "group2", Email: "group2@mail.com"})},
				},
				HashCode: hash.Get(
					&model.GroupsResult{
						Items: 2,
						Resources: []*model.Group{
							{ID: "1", Name: "group1", Email: "group1@mail.com", HashCode: hash.Get(model.Group{ID: "1", Name: "group1", Email: "group1@mail.com"})},
							{ID: "2", Name: "group2", Email: "group2@mail.com", HashCode: hash.Get(model.Group{ID: "2", Name: "group2", Email: "group2@mail.com"})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				f.ds.EXPECT().ListGroups(gomock.Eq([]string{""})).Return(nil, ErrListingGroups).Times(1)
			},
			args:    args{ctx: context.TODO(), filter: []string{""}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			f := fields{ds: mocks.NewMockGoogleProviderService(mockCtrl)}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			g := &IdentityProvider{
				ps: f.ds,
			}

			got, err := g.GetGroups(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleProvider.GetGroups() = %s, want %s", utils.ToJSON(got), utils.ToJSON(tt.want))
			}
		})
	}
}

func Test_GoogleProvider_GetUsers(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx    context.Context
		filter []string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *model.UsersResult
		wantErr bool
	}{
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				googleUsers := make([]*admin.User, 0)
				googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user1@mail.com", Id: "1", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false})
				googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user2@mail.com", Id: "2", Name: &admin.UserName{GivenName: "user", FamilyName: "2"}, Suspended: true})

				f.ds.EXPECT().ListUsers(gomock.Eq([]string{""})).Return(googleUsers, nil).Times(1)
			},
			args: args{ctx: context.TODO(), filter: []string{""}},
			want: &model.UsersResult{
				Items: 2,
				Resources: []*model.User{
					{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(model.User{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true})},
					{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(model.User{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false})},
				},
				HashCode: hash.Get(
					&model.UsersResult{
						Items: 2,
						Resources: []*model.User{
							{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(model.User{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true})},
							{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(model.User{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				f.ds.EXPECT().ListUsers(gomock.Eq([]string{""})).Return(nil, ErrListingUsers).Times(1)
			},
			args:    args{ctx: context.TODO(), filter: []string{""}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			f := fields{ds: mocks.NewMockGoogleProviderService(mockCtrl)}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			g := &IdentityProvider{
				ps: f.ds,
			}

			got, err := g.GetUsers(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleProvider.GetUsers() = %s, want %s", utils.ToJSON(got), utils.ToJSON(tt.want))
			}
		})
	}
}

func Test_GoogleProvider_GetGroupMembers(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *model.MembersResult
		wantErr bool
	}{
		{
			name: "Should return MembersResult and no error",
			prepare: func(f *fields) {
				googleGroupMembers := make([]*admin.Member, 0)
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user1@mail.com", Id: "1"})
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user2@mail.com", Id: "2"})

				f.ds.EXPECT().ListGroupMembers(gomock.Eq("")).Return(googleGroupMembers, nil).Times(1)
			},
			args: args{ctx: context.TODO(), id: ""},
			want: &model.MembersResult{
				Items: 2,
				Resources: []*model.Member{
					{ID: "1", Email: "user1@mail.com", HashCode: hash.Get(model.Member{ID: "1", Email: "user1@mail.com"})},
					{ID: "2", Email: "user2@mail.com", HashCode: hash.Get(model.Member{ID: "2", Email: "user2@mail.com"})},
				},
				HashCode: hash.Get(
					&model.MembersResult{
						Items: 2,
						Resources: []*model.Member{
							{ID: "1", Email: "user1@mail.com", HashCode: hash.Get(model.Member{ID: "1", Email: "user1@mail.com"})},
							{ID: "2", Email: "user2@mail.com", HashCode: hash.Get(model.Member{ID: "2", Email: "user2@mail.com"})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				f.ds.EXPECT().ListGroupMembers(gomock.Eq("")).Return(nil, ErrListingGroupMembers).Times(1)
			},
			args:    args{ctx: context.TODO(), id: ""},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			f := fields{ds: mocks.NewMockGoogleProviderService(mockCtrl)}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			g := &IdentityProvider{
				ps: f.ds,
			}

			got, err := g.GetGroupMembers(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleProvider.GetGroupMembers() = %s, want %s", utils.ToJSON(got), utils.ToJSON(tt.want))
			}
		})
	}
}

func Test_GoogleProvider_GetUsersFromGroupMembers(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx context.Context
		mbr *model.MembersResult
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *model.UsersResult
		wantErr bool
	}{
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				googleUser1 := &admin.User{PrimaryEmail: "user1@mail.com", Id: "1", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false}
				googleUser2 := &admin.User{PrimaryEmail: "user2@mail.com", Id: "2", Name: &admin.UserName{GivenName: "user", FamilyName: "2"}, Suspended: true}

				gomock.InOrder(
					f.ds.EXPECT().GetUser(gomock.Eq("1")).Return(googleUser1, nil).Times(1),
					f.ds.EXPECT().GetUser(gomock.Eq("2")).Return(googleUser2, nil).Times(1),
				)
			},
			args: args{
				ctx: context.TODO(),
				mbr: &model.MembersResult{
					Items: 2,
					Resources: []*model.Member{
						{ID: "1", Email: "user1@mail.com"},
						{ID: "2", Email: "user2@mail.com"},
					},
				},
			},
			want: &model.UsersResult{
				Items: 2,
				Resources: []*model.User{
					{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true})},
					{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(&model.User{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false})},
				},
				HashCode: hash.Get(
					&model.UsersResult{
						Items: 2,
						Resources: []*model.User{
							{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{ID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user1@mail.com", DisplayName: "user 1", Active: true})},
							{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(&model.User{ID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user2@mail.com", DisplayName: "user 2", Active: false})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				f.ds.EXPECT().GetUser(gomock.Eq("")).Return(nil, ErrGettingUser).Times(1)
			},
			args: args{
				ctx: context.TODO(),
				mbr: &model.MembersResult{
					Items: 0,
					Resources: []*model.Member{
						{ID: "", Email: "user1@mail.com"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			f := fields{ds: mocks.NewMockGoogleProviderService(mockCtrl)}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			g := &IdentityProvider{
				ps: f.ds,
			}

			got, err := g.GetUsersFromGroupMembers(tt.args.ctx, tt.args.mbr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetUsersFromGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleProvider.GetUsersFromGroupMembers() = %s, want %s", utils.ToJSON(got), utils.ToJSON(tt.want))
			}
		})
	}
}
