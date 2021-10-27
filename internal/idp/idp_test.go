package idp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/idp"
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

func TestGoogleProvider_GetGroups(t *testing.T) {
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
			name: "Should return empty GroupsResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroups := make([]*admin.Group, 0)

				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(googleGroups, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.GroupsResult{
				Items:     0,
				Resources: make([]model.Group, 0),
			},
			wantErr: false,
		},
		{
			name: "Should return GroupsResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroups := make([]*admin.Group, 0)
				googleGroups = append(googleGroups, &admin.Group{Email: "group1@mail.com", Id: "1", Name: "group1"})
				googleGroups = append(googleGroups, &admin.Group{Email: "group2@mail.com", Id: "2", Name: "group2"})

				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(googleGroups, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.GroupsResult{
				Items: 2,
				Resources: []model.Group{
					{IPID: "1", Name: "group1", Email: "group1@mail.com", HashCode: hash.Get(model.Group{IPID: "1", Name: "group1", Email: "group1@mail.com"})},
					{IPID: "2", Name: "group2", Email: "group2@mail.com", HashCode: hash.Get(model.Group{IPID: "2", Name: "group2", Email: "group2@mail.com"})},
				},
				HashCode: hash.Get(
					&model.GroupsResult{
						Items: 2,
						Resources: []model.Group{
							{IPID: "1", Name: "group1", Email: "group1@mail.com", HashCode: hash.Get(model.Group{IPID: "1", Name: "group1", Email: "group1@mail.com"})},
							{IPID: "2", Name: "group2", Email: "group2@mail.com", HashCode: hash.Get(model.Group{IPID: "2", Name: "group2", Email: "group2@mail.com"})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(nil, errors.New("test error")).Times(1)
			},
			args:    args{ctx: context.Background(), filter: []string{""}},
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

func TestGoogleProvider_GetUsers(t *testing.T) {
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
			name: "Should return empty UsersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUsers := make([]*admin.User, 0)

				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return(googleUsers, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.UsersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]model.User, 0),
			},
			wantErr: false,
		},
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUsers := make([]*admin.User, 0)
				googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user.1@mail.com", Id: "1", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false})
				googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user.2@mail.com", Id: "2", Name: &admin.UserName{GivenName: "user", FamilyName: "2"}, Suspended: true})

				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return(googleUsers, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.UsersResult{
				Items: 2,
				Resources: []model.User{
					{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
					{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(model.User{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false})},
				},
				HashCode: hash.Get(
					&model.UsersResult{
						Items: 2,
						Resources: []model.User{
							{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
							{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(model.User{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return(nil, errors.New("test error")).Times(1)
			},
			args:    args{ctx: context.Background(), filter: []string{""}},
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

func TestGoogleProvider_GetGroupMembers(t *testing.T) {
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
			name: "Should return empty MembersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroupMembers := make([]*admin.Member, 0)

				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("")).Return(googleGroupMembers, nil).Times(1)
			},
			args: args{ctx: context.Background(), id: ""},
			want: &model.MembersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]model.Member, 0),
			},
			wantErr: false,
		},
		{
			name: "Should return MembersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroupMembers := make([]*admin.Member, 0)
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.1@mail.com", Id: "1"})
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.2@mail.com", Id: "2"})

				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("")).Return(googleGroupMembers, nil).Times(1)
			},
			args: args{ctx: context.Background(), id: ""},
			want: &model.MembersResult{
				Items: 2,
				Resources: []model.Member{
					{IPID: "1", Email: "user.1@mail.com", HashCode: hash.Get(model.Member{IPID: "1", Email: "user.1@mail.com"})},
					{IPID: "2", Email: "user.2@mail.com", HashCode: hash.Get(model.Member{IPID: "2", Email: "user.2@mail.com"})},
				},
				HashCode: hash.Get(
					&model.MembersResult{
						Items: 2,
						Resources: []model.Member{
							{IPID: "1", Email: "user.1@mail.com", HashCode: hash.Get(model.Member{IPID: "1", Email: "user.1@mail.com"})},
							{IPID: "2", Email: "user.2@mail.com", HashCode: hash.Get(model.Member{IPID: "2", Email: "user.2@mail.com"})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("")).Return(nil, errors.New("test error")).Times(1)
			},
			args:    args{ctx: context.Background(), id: ""},
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

func TestGoogleProvider_GetUsersFromGroupMembers(t *testing.T) {
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
			name:    "empty mbr argument",
			prepare: func(f *fields) {},
			args: args{
				ctx: context.Background(),
				mbr: &model.MembersResult{},
			},
			want: &model.UsersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]model.User, 0),
			},
			wantErr: false,
		},
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUser1 := &admin.User{PrimaryEmail: "user.1@mail.com", Id: "1", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false}
				googleUser2 := &admin.User{PrimaryEmail: "user.2@mail.com", Id: "2", Name: &admin.UserName{GivenName: "user", FamilyName: "2"}, Suspended: true}

				gomock.InOrder(
					f.ds.EXPECT().GetUser(ctx, gomock.Eq("1")).Return(googleUser1, nil).Times(1),
					f.ds.EXPECT().GetUser(ctx, gomock.Eq("2")).Return(googleUser2, nil).Times(1),
				)
			},
			args: args{
				ctx: context.Background(),
				mbr: &model.MembersResult{
					Items: 2,
					Resources: []model.Member{
						{IPID: "1", Email: "user.1@mail.com"},
						{IPID: "2", Email: "user.2@mail.com"},
					},
				},
			},
			want: &model.UsersResult{
				Items: 2,
				Resources: []model.User{
					{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
					{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(&model.User{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false})},
				},
				HashCode: hash.Get(
					&model.UsersResult{
						Items: 2,
						Resources: []model.User{
							{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
							{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false, HashCode: hash.Get(&model.User{IPID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@mail.com", DisplayName: "user 2", Active: false})},
						},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().GetUser(ctx, gomock.Eq("")).Return(nil, errors.New("test error")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				mbr: &model.MembersResult{
					Items: 0,
					Resources: []model.Member{
						{IPID: "", Email: "user.1@mail.com"},
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

func TestGoogleProvider_GetUsersAndGroupsUsers(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx    context.Context
		groups *model.GroupsResult
	}

	tests := []struct {
		name            string
		prepare         func(f *fields)
		args            args
		wantUsers       *model.UsersResult
		wantGroupsUsers *model.GroupsUsersResult
		wantErr         bool
	}{
		{
			name:    "empty gr argument",
			prepare: func(f *fields) {},
			args: args{
				ctx:    context.Background(),
				groups: &model.GroupsResult{},
			},
			wantUsers: &model.UsersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]model.User, 0),
			},
			wantGroupsUsers: &model.GroupsUsersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]model.GroupUsers, 0),
			},
			wantErr: false,
		},
		{
			name: "same user in two groups",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUser1 := &admin.User{Id: "1", PrimaryEmail: "user.1@mail.com", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false}
				googleMembers := []*admin.Member{{Id: "1", Email: "user.1@mail.com"}}

				gomock.InOrder(
					f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("1")).Return(googleMembers, nil).Times(1),
					f.ds.EXPECT().GetUser(ctx, gomock.Eq("1")).Return(googleUser1, nil).Times(1),

					f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("2")).Return(googleMembers, nil).Times(1),
					f.ds.EXPECT().GetUser(ctx, gomock.Eq("1")).Return(googleUser1, nil).Times(1),
				)
			},
			args: args{
				ctx: context.Background(),
				groups: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					},
				},
			},
			wantUsers: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
				},
				HashCode: hash.Get(
					&model.UsersResult{
						Items: 1,
						Resources: []model.User{
							{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})},
						},
					},
				),
			},
			wantGroupsUsers: &model.GroupsUsersResult{
				Items: 2,
				Resources: []model.GroupUsers{
					{Items: 1, HashCode: hash.Get(model.GroupUsers{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}}), Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}},
					{Items: 1, HashCode: hash.Get(model.GroupUsers{Items: 1, Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}}), Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}},
				},
				HashCode: hash.Get(
					&model.GroupsUsersResult{
						Items: 2,
						Resources: []model.GroupUsers{
							{Items: 1, HashCode: hash.Get(model.GroupUsers{Items: 1, Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}}), Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}},
							{Items: 1, HashCode: hash.Get(model.GroupUsers{Items: 1, Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}}), Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"}, Resources: []model.User{{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true, HashCode: hash.Get(&model.User{IPID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@mail.com", DisplayName: "user 1", Active: true})}}},
						},
					},
				),
			},
			wantErr: false,
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

			gotUsers, gotGroupsUsers, err := g.GetUsersAndGroupsUsers(tt.args.ctx, tt.args.groups)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetUsersAndGroupsUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotUsers, tt.wantUsers) {
				t.Errorf("GoogleProvider.GetUsersAndGroupsUsers() -> got Users = %s, want %s", utils.ToJSON(gotUsers), utils.ToJSON(tt.wantUsers))
			}

			if !reflect.DeepEqual(gotGroupsUsers, tt.wantGroupsUsers) {
				t.Errorf("GoogleProvider.GetUsersAndGroupsUsers() -> got GroupsUsers = %s, want %s", utils.ToJSON(gotGroupsUsers), utils.ToJSON(tt.wantGroupsUsers))
			}
		})
	}
}
