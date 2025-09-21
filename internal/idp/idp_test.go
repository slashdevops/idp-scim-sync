package idp

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"go.uber.org/mock/gomock"

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

func TestGetGroups(t *testing.T) {
	g1 := &model.Group{IPID: "1", Name: "group 1", Email: "group1@mail.com"}
	g1.SetHashCode()
	g2 := &model.Group{IPID: "2", Name: "group 2", Email: "group2@mail.com"}
	g2.SetHashCode()
	g4 := &model.Group{IPID: "4", Name: "group 4", Email: "group4@mail.com"}
	g4.SetHashCode()

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
				Resources: make([]*model.Group, 0),
			},
			wantErr: false,
		},
		{
			name: "should return empty GroupsResult when no groups are provided",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return([]*admin.Group{}, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantErr: false,
		},
		{
			name: "Should return GroupsResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroups := make([]*admin.Group, 0)
				googleGroups = append(googleGroups, &admin.Group{Email: "group1@mail.com", Id: "1", Name: "group 1"})
				googleGroups = append(googleGroups, &admin.Group{Email: "group2@mail.com", Id: "2", Name: "group 2"})

				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(googleGroups, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.GroupsResult{
				Items:     2,
				Resources: []*model.Group{g1, g2},
			},
			wantErr: false,
		},
		{
			name: "Should keep only one of the repeated groups name",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroups := make([]*admin.Group, 0)
				googleGroups = append(googleGroups, &admin.Group{Email: "group1@mail.com", Id: "1", Name: "group 1"})
				googleGroups = append(googleGroups, &admin.Group{Email: "group2@mail.com", Id: "2", Name: "group 2"})

				googleGroups = append(googleGroups, &admin.Group{Email: "group3@mail.com", Id: "3", Name: "group 2"}) // Repeated group name
				googleGroups = append(googleGroups, &admin.Group{Email: "group4@mail.com", Id: "4", Name: "group 4"})

				f.ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(googleGroups, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.GroupsResult{
				Items:     3,
				Resources: []*model.Group{g1, g2, g4},
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
			if !tt.wantErr {
				tt.want.SetHashCode()
			}

			got, err := g.GetGroups(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GetGroups() (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkGetGroups(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	// given
	filter := []string{""}
	ctx := context.Background()
	googleGroups := make([]*admin.Group, 0)
	googleGroups = append(googleGroups, &admin.Group{Email: "group.1@mail.com", Id: "1", Name: "group 1"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.2@mail.com", Id: "2", Name: "group 2"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.3@mail.com", Id: "3", Name: "group 3"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.4@mail.com", Id: "4", Name: "group 4"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.5@mail.com", Id: "5", Name: "group 5"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.6@mail.com", Id: "6", Name: "group 6"})
	googleGroups = append(googleGroups, &admin.Group{Email: "group.7@mail.com", Id: "7", Name: "group 7"})

	ds := mocks.NewMockGoogleProviderService(mockCtrl)
	g := &IdentityProvider{ps: ds}

	// when
	ds.EXPECT().ListGroups(ctx, gomock.Eq([]string{""})).Return(googleGroups, nil).AnyTimes()

	b.Run("benchmark GetUsers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := g.GetGroups(ctx, filter)
			if err != nil {
				b.Errorf("GoogleProvider.GetGroups() error = %v", err)
				return
			}
		}
	})
}

func TestGetUsers(t *testing.T) {
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
			name: "Should return error",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return(nil, errors.New("test error")).Times(1)
			},
			args:    args{ctx: context.Background(), filter: []string{""}},
			want:    nil,
			wantErr: true,
		},
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
				Resources: make([]*model.User, 0),
			},
			wantErr: false,
		},
		{
			name: "should return empty UsersResult when no users are provided",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return([]*admin.User{}, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantErr: false,
		},
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUsers := make([]*admin.User, 0)
				googleUsers = append(googleUsers,
					&admin.User{
						Id:           "1",
						PrimaryEmail: "user.1@mail.com",
						Name:         &admin.UserName{GivenName: "user", FamilyName: "1"},
						Suspended:    false,
					},
				)
				googleUsers = append(googleUsers,
					&admin.User{
						Id:           "2",
						PrimaryEmail: "user.2@mail.com",
						Name:         &admin.UserName{GivenName: "user", FamilyName: "2"},
						Suspended:    true,
						Emails: []admin.UserEmail{
							{
								Address: "user.2@mailcom",
								Type:    "work",
								Primary: true,
							},
						},
					},
				)

				f.ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return(googleUsers, nil).Times(1)
			},
			args: args{ctx: context.Background(), filter: []string{""}},
			want: &model.UsersResult{
				Items: 2,
				Resources: []*model.User{
					model.UserBuilder().
						WithIPID("1").
						WithName(&model.Name{GivenName: "user", FamilyName: "1"}).
						WithDisplayName("user 1").
						WithUserName("user.1@mail.com").
						WithActive(true).
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.1@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
					model.UserBuilder().
						WithIPID("2").
						WithName(&model.Name{GivenName: "user", FamilyName: "2"}).
						WithDisplayName("user 2").
						WithUserName("user.2@mail.com").
						WithActive(false).
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.2@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
				},
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

			// trigger the call to calculate the hash
			if !tt.wantErr {
				tt.want.SetHashCode()
			}

			got, err := g.GetUsers(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetUsers() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GetUsers() (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkGetUsers(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	// given
	filter := []string{""}
	ctx := context.Background()
	googleUsers := make([]*admin.User, 0)
	googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user.1@mail.com", Id: "1", Name: &admin.UserName{GivenName: "user", FamilyName: "1"}, Suspended: false})
	googleUsers = append(googleUsers, &admin.User{PrimaryEmail: "user.2@mail.com", Id: "2", Name: &admin.UserName{GivenName: "user", FamilyName: "2"}, Suspended: true})

	ds := mocks.NewMockGoogleProviderService(mockCtrl)
	g := &IdentityProvider{ps: ds}

	// when
	ds.EXPECT().ListUsers(ctx, gomock.Eq(filter)).Return(googleUsers, nil).AnyTimes()

	b.Run("benchmark GetUsers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := g.GetUsers(ctx, filter)
			if err != nil {
				b.Errorf("GoogleProvider.GetUsers() error = %v", err)
				return
			}
		}
	})
}

func TestGetGroupMembers(t *testing.T) {
	m1 := &model.Member{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"}
	m1.SetHashCode()
	m2 := &model.Member{IPID: "2", Email: "user.2@mail.com", Status: "suspended"}
	m2.SetHashCode()

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
			name:    "Should return a nil object and an error when id is empty",
			prepare: func(f *fields) {},
			args:    args{ctx: context.Background(), id: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return error when ListGroupMembers return error",
			prepare: func(f *fields) {
				ctx := context.Background()

				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("1"), gomock.Any()).Return(nil, errors.New("test error")).Times(1)
			},
			args:    args{ctx: context.Background(), id: "1"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return MembersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroupMembers := make([]*admin.Member, 0)
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.1@mail.com", Id: "1", Status: "ACTIVE"})
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.2@mail.com", Id: "2", Status: "suspended"})

				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("1"), gomock.Any()).Return(googleGroupMembers, nil).Times(1)
			},
			args: args{ctx: context.Background(), id: "1"},
			want: &model.MembersResult{
				Items:     2,
				Resources: []*model.Member{m1, m2},
			},
			wantErr: false,
		},
		{
			name: "Should return MembersResult with only one member when member type is GROUP",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleGroupMembers := make([]*admin.Member, 0)
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "group.1@mail.com", Id: "1", Type: "GROUP"})
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.2@mail.com", Id: "2", Status: "suspended"})

				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("1"), gomock.Any()).Return(googleGroupMembers, nil).Times(1)
			},
			args: args{ctx: context.Background(), id: "1"},
			want: &model.MembersResult{
				Items:     1,
				Resources: []*model.Member{m2},
			},
			wantErr: false,
		},
		{
			name: "should return empty MembersResult when no members are provided",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListGroupMembers(ctx, gomock.Eq("1"), gomock.Any()).Return([]*admin.Member{}, nil).Times(1)
			},
			args: args{ctx: context.Background(), id: "1"},
			want: &model.MembersResult{
				Items:     0,
				Resources: []*model.Member{},
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

			if !tt.wantErr {
				tt.want.SetHashCode()
			}

			got, err := g.GetGroupMembers(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GetGroupMembers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetUsersByGroupsMembers(t *testing.T) {
	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx context.Context
		gmr *model.GroupsMembersResult
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
				gmr: &model.GroupsMembersResult{},
			},
			want: &model.UsersResult{
				Items:     0,
				HashCode:  "",
				Resources: make([]*model.User, 0),
			},
			wantErr: false,
		},
		{
			name:    "nil gmr argument",
			prepare: func(f *fields) {},
			args: args{
				ctx: context.Background(),
				gmr: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return error when ListUsers fails",
			prepare: func(f *fields) {
				ctx := context.Background()
				f.ds.EXPECT().ListUsers(ctx, gomock.Any()).Return(nil, errors.New("test error")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Items: 1,
							Group: &model.Group{IPID: "1", Name: "group 1", Email: "group1@mail.com"},
							Resources: []*model.Member{
								{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return UsersResult and no error",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUser1 := &admin.User{
					Id:           "1",
					PrimaryEmail: "user.1@mail.com",
					Name:         &admin.UserName{GivenName: "user", FamilyName: "1"},
					Suspended:    false,
					// Emails: []admin.UserEmail{
					// 	{
					// 		Address: "user.1@mail.com",
					// 		Type:    "work",
					// 		Primary: true,
					// 	},
					// },
				}
				googleUser2 := &admin.User{
					Id:           "2",
					PrimaryEmail: "user.2@mail.com",
					Name:         &admin.UserName{GivenName: "user", FamilyName: "2"},
					Suspended:    true,
					// Emails: []admin.UserEmail{
					// 	{
					// 		Address: "user.2@mail.com",
					// 		Type:    "work",
					// 		Primary: true,
					// 	},
					// },
				}
				googleUser3 := &admin.User{
					Id:           "3",
					PrimaryEmail: "user.3@mail.com",
					Name:         &admin.UserName{GivenName: "user", FamilyName: "3"},
					Suspended:    true,
					// Emails: []admin.UserEmail{
					// 	{
					// 		Address: "user.3@mail.com",
					// 		Type:    "work",
					// 		Primary: true,
					// 	},
					// },
				}
				googleUser4 := &admin.User{
					Id:           "4",
					PrimaryEmail: "user.4@mail.com",
					Name:         &admin.UserName{GivenName: "user", FamilyName: "4"},
					Suspended:    true,
					// Emails: []admin.UserEmail{
					// 	{
					// 		Address: "user.4@mail.com",
					// 		Type:    "work",
					// 		Primary: true,
					// 	},
					// },
				}

				gomock.InOrder(
					f.ds.EXPECT().ListUsers(ctx, gomock.Any()).Return([]*admin.User{googleUser1, googleUser2, googleUser3, googleUser4}, nil).Times(1),
				)
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: &model.Group{IPID: "1", Name: "group 1", Email: "group1@mail.com"},
							Resources: []*model.Member{
								{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
								{IPID: "2", Email: "user.2@mail.com", Status: "ACTIVE"},
							},
						},
						{
							Items: 2,
							Group: &model.Group{IPID: "1", Name: "group 1", Email: "group1@mail.com"},
							Resources: []*model.Member{
								{IPID: "3", Email: "user.3@mail.com", Status: "ACTIVE"},
								{IPID: "4", Email: "user.4@mail.com", Status: "ACTIVE"},
							},
						},
					},
				},
			},
			want: &model.UsersResult{
				Items: 4,
				Resources: []*model.User{
					model.UserBuilder().
						WithIPID("1").
						WithName(&model.Name{GivenName: "user", FamilyName: "1"}).
						WithDisplayName("user 1").
						WithActive(true).
						WithUserName("user.1@mail.com").
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.1@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
					model.UserBuilder().
						WithIPID("2").
						WithName(&model.Name{GivenName: "user", FamilyName: "2"}).
						WithDisplayName("user 2").
						WithActive(false).
						WithUserName("user.2@mail.com").
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.2@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
					model.UserBuilder().
						WithIPID("3").
						WithName(&model.Name{GivenName: "user", FamilyName: "3"}).
						WithDisplayName("user 3").
						WithActive(false).
						WithUserName("user.3@mail.com").
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.3@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
					model.UserBuilder().
						WithIPID("4").
						WithName(&model.Name{GivenName: "user", FamilyName: "4"}).
						WithDisplayName("user 4").
						WithActive(false).
						WithUserName("user.4@mail.com").
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.4@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
				},
			},
			wantErr: false,
		},
		{
			name: "Should handle duplicate emails correctly",
			prepare: func(f *fields) {
				ctx := context.Background()
				googleUser1 := &admin.User{
					Id:           "1",
					PrimaryEmail: "user.1@mail.com",
					Name:         &admin.UserName{GivenName: "user", FamilyName: "1"},
					Suspended:    false,
				}

				f.ds.EXPECT().ListUsers(ctx, gomock.Any()).Return([]*admin.User{googleUser1}, nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 1,
							Group: &model.Group{IPID: "1", Name: "group 1", Email: "group1@mail.com"},
							Resources: []*model.Member{
								{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
							},
						},
						{
							Items: 1,
							Group: &model.Group{IPID: "2", Name: "group 2", Email: "group2@mail.com"},
							Resources: []*model.Member{
								{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"}, // Same user in multiple groups
							},
						},
					},
				},
			},
			want: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					model.UserBuilder().
						WithIPID("1").
						WithName(&model.Name{GivenName: "user", FamilyName: "1"}).
						WithDisplayName("user 1").
						WithActive(true).
						WithUserName("user.1@mail.com").
						WithEmails(
							[]model.Email{
								model.EmailBuilder().
									WithValue("user.1@mail.com").
									WithType("work").
									WithPrimary(true).
									Build(),
							},
						).
						Build(),
				},
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

			if !tt.wantErr {
				tt.want.SetHashCode()
			}

			got, err := g.GetUsersByGroupsMembers(tt.args.ctx, tt.args.gmr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetUsersFromGroupMembers() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GetUsersByGroupsMembers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetGroupsMembers(t *testing.T) {
	m1 := model.Member{IPID: "1", Email: "user.1@mail.com", Status: "ACTIVE"}
	m1.SetHashCode()

	m2 := model.Member{IPID: "2", Email: "user.2@mail.com", Status: "ACTIVE"}
	m2.SetHashCode()

	g1 := &model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"}
	g1.SetHashCode()

	gm1 := &model.GroupMembers{
		Items: 2,
		Group: g1,
		Resources: []*model.Member{
			&m1,
			&m2,
		},
	}
	gm1.SetHashCode()
	gm2 := &model.GroupMembers{
		Items:     0,
		Group:     g1,
		Resources: []*model.Member{},
	}
	gm2.SetHashCode()
	gmr1 := &model.GroupsMembersResult{
		Items: 1,
		Resources: []*model.GroupMembers{
			gm1,
		},
	}
	gmr1.SetHashCode()
	gmr2 := &model.GroupsMembersResult{
		Items: 1,
		Resources: []*model.GroupMembers{
			gm2,
		},
	}
	gmr2.SetHashCode()

	type fields struct {
		ds *mocks.MockGoogleProviderService
	}

	type args struct {
		ctx context.Context
		gr  *model.GroupsResult
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *model.GroupsMembersResult
		wantErr bool
	}{
		{
			name:    "Should return error when gr is nil",
			prepare: func(f *fields) {},
			args:    args{ctx: context.Background(), gr: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return empty GroupsMembersResult when no groups are provided",
			prepare: func(f *fields) {
			},
			args: args{
				ctx: context.Background(),
				gr:  &model.GroupsResult{},
			},
			want:    &model.GroupsMembersResult{Items: 0, Resources: []*model.GroupMembers{}},
			wantErr: false,
		},
		{
			name:    "Should return empty GroupsMembersResult when gr items is 0",
			prepare: func(f *fields) {},
			args: args{
				ctx: context.Background(),
				gr:  &model.GroupsResult{Items: 0, Resources: []*model.Group{}},
			},
			want:    &model.GroupsMembersResult{Items: 0, Resources: []*model.GroupMembers{}},
			wantErr: false,
		},
		{
			name: "Should return error when ListGroupMembers return error",
			prepare: func(f *fields) {
				f.ds.EXPECT().ListGroupMembersBatch(gomock.Any(), []string{"1"}, gomock.Any()).Return(nil, errors.New("test error")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return MembersResult and no error",
			prepare: func(f *fields) {
				googleGroupMembers := make([]*admin.Member, 0)
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.1@mail.com", Id: "1", Status: "ACTIVE"})
				googleGroupMembers = append(googleGroupMembers, &admin.Member{Email: "user.2@mail.com", Id: "2", Status: "ACTIVE"})

				membersMap := map[string][]*admin.Member{
					"1": googleGroupMembers,
				}
				f.ds.EXPECT().ListGroupMembersBatch(gomock.Any(), []string{"1"}, gomock.Any()).Return(membersMap, nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					},
				},
			},
			want:    gmr1,
			wantErr: false,
		},
		{
			name: "Should return MembersResult and no error when members is empty",
			prepare: func(f *fields) {
				googleGroupMembers := make([]*admin.Member, 0)
				membersMap := map[string][]*admin.Member{
					"1": googleGroupMembers,
				}
				f.ds.EXPECT().ListGroupMembersBatch(gomock.Any(), []string{"1"}, gomock.Any()).Return(membersMap, nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					},
				},
			},
			want:    gmr2,
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

			if !tt.wantErr {
				tt.want.SetHashCode()
			}

			got, err := g.GetGroupsMembers(tt.args.ctx, tt.args.gr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleProvider.GetGroupsMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GetGroupsMembers() (-want +got):\n%s", diff)
			}
		})
	}
}
