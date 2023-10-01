package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestGroupsOperations(t *testing.T) {
	type args struct {
		idp   *GroupsResult
		state *GroupsResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *GroupsResult
		wantUpdate *GroupsResult
		wantEqual  *GroupsResult
		wantDelete *GroupsResult
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				idp:   GroupsResultBuilder().Build(),
				state: GroupsResultBuilder().Build(),
			},
			wantCreate: GroupsResultBuilder().Build(),
			wantUpdate: GroupsResultBuilder().Build(),
			wantEqual:  GroupsResultBuilder().Build(),
			wantDelete: GroupsResultBuilder().Build(),
			wantErr:    false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: GroupsResultBuilder().Build(),
			},
			wantCreate: nil,
			wantUpdate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "nil state",
			args: args{
				idp:   GroupsResultBuilder().Build(),
				state: nil,
			},
			wantCreate: nil,
			wantUpdate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "2 equals",
			args: args{
				idp: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("2").WithName("name2").WithEmail("2@mail.com").Build(),
					},
				).Build(),
				state: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("2").WithName("name2").WithEmail("2@mail.com").Build(),
					},
				).Build(),
			},
			wantCreate: GroupsResultBuilder().Build(),
			wantUpdate: GroupsResultBuilder().Build(),
			wantEqual: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("1").WithName("name1").WithEmail("1@mail.com").Build(),
					GroupBuilder().WithIPID("2").WithName("name2").WithEmail("2@mail.com").Build(),
				},
			).Build(),
			wantDelete: GroupsResultBuilder().Build(),
			wantErr:    false,
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("2").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
					},
				).Build(),
				state: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
					},
				).Build(),
			},
			wantCreate: GroupsResultBuilder().Build(),
			wantUpdate: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("2").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
				},
			).Build(),
			wantEqual: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
				},
			).Build(),
			wantDelete: GroupsResultBuilder().Build(),
			wantErr:    false,
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("2").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
					},
				).Build(),
				state: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
						GroupBuilder().WithIPID("3").WithSCIMID("33").WithName("name3").WithEmail("3@mail.com").Build(),
					},
				).Build(),
			},
			wantCreate: GroupsResultBuilder().Build(),
			wantUpdate: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("2").WithSCIMID("22").WithName("name2").WithEmail("2@mail.com").Build(),
				},
			).Build(),
			wantEqual: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
				},
			).Build(),
			wantDelete: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("3").WithSCIMID("33").WithName("name3").WithEmail("3@mail.com").Build(),
				},
			).Build(),
			wantErr: false,
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("2").WithSCIMID("different").WithName("name2").WithEmail("2@mail.com").Build(),
						GroupBuilder().WithIPID("4").WithSCIMID("44").WithName("name4").WithEmail("4@mail.com").Build(),
					},
				).Build(),
				state: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
						GroupBuilder().WithIPID("").WithSCIMID("different").WithName("name2").WithEmail("2@mail.com").Build(),
						GroupBuilder().WithIPID("3").WithSCIMID("33").WithName("name3").WithEmail("3@mail.com").Build(),
					},
				).Build(),
			},
			wantCreate: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("4").WithSCIMID("44").WithName("name4").WithEmail("4@mail.com").Build(),
				},
			).Build(),
			wantUpdate: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("2").WithSCIMID("different").WithName("name2").WithEmail("2@mail.com").Build(),
				},
			).Build(),
			wantEqual: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
				},
			).Build(),
			wantDelete: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("3").WithSCIMID("33").WithName("name3").WithEmail("3@mail.com").Build(),
				},
			).Build(),
			wantErr: false,
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("1").WithSCIMID("11").WithName("name1").WithEmail("1@mail.com").Build(),
					},
				).Build(),
				state: GroupsResultBuilder().WithResources(
					[]*Group{
						GroupBuilder().WithIPID("3").WithSCIMID("22").WithName("name1").WithEmail("1@mail.com").Build(),
					},
				).Build(),
			},
			wantCreate: GroupsResultBuilder().Build(),
			wantUpdate: GroupsResultBuilder().WithResources(
				[]*Group{
					GroupBuilder().WithIPID("1").WithSCIMID("22").WithName("name1").WithEmail("1@mail.com").Build(),
				},
			).Build(),
			wantEqual:  GroupsResultBuilder().Build(),
			wantDelete: GroupsResultBuilder().Build(),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.wantCreate.SetHashCode()
				tt.wantUpdate.SetHashCode()
				tt.wantEqual.SetHashCode()
				tt.wantDelete.SetHashCode()
			}

			gotCreate, gotUpdate, gotEqual, gotDelete, err := GroupsOperations(tt.args.idp, tt.args.state)

			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsOperations() error = %v", err)
				return
			}

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantCreate, gotCreate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotCreate (-tt.wantCreate +gotCreate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantUpdate, gotUpdate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotUpdate (-tt.wantUpdate +gotUpdate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantEqual, gotEqual, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotEqual (-tt.wantEqual +gotEqual):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantDelete, gotDelete, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotDelete (-tt.wantDelete +gotDelete):\n%s", diff)
			}
		})
	}
}

func TestUsersOperations(t *testing.T) {
	type args struct {
		idp   *UsersResult
		state *UsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *UsersResult
		wantUpdate *UsersResult
		wantEqual  *UsersResult
		wantDelete *UsersResult
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				idp:   UsersResultBuilder().Build(),
				state: UsersResultBuilder().Build(),
			},
			wantCreate: UsersResultBuilder().Build(),
			wantUpdate: UsersResultBuilder().Build(),
			wantEqual:  UsersResultBuilder().Build(),
			wantDelete: UsersResultBuilder().Build(),
			wantErr:    false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: UsersResultBuilder().Build(),
			},
			wantCreate: nil,
			wantUpdate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "nil state",
			args: args{
				idp:   UsersResultBuilder().Build(),
				state: nil,
			},
			wantCreate: nil,
			wantUpdate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "2 equals",
			args: args{
				idp: UsersResultBuilder().
					WithResources(
						[]*User{
							UserBuilder().
								WithIPID("1").
								WithEmail(
									EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build(),
								).
								WithFamilyName("1").
								WithGivenName("user").
								WithDisplayName("user 1").
								WithActive(true).
								Build(),
							UserBuilder().
								WithIPID("2").
								WithEmail(
									EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build(),
								).
								WithFamilyName("2").
								WithGivenName("user").
								WithDisplayName("user 2").
								WithActive(true).
								Build(),
						},
					).Build(),
				state: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().
							WithIPID("1").
							WithEmail(
								EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build(),
							).
							WithFamilyName("1").
							WithGivenName("user").
							WithDisplayName("user 1").
							WithActive(true).
							Build(),
						UserBuilder().WithIPID("2").
							WithEmail(
								EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build(),
							).
							WithFamilyName("2").
							WithGivenName("user").
							WithDisplayName("user 2").
							WithActive(true).
							Build(),
					},
				).Build(),
			},
			wantCreate: UsersResultBuilder().Build(),
			wantUpdate: UsersResultBuilder().Build(),
			wantEqual: UsersResultBuilder().
				WithResources(
					[]*User{
						UserBuilder().
							WithIPID("1").
							WithEmail(
								EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build(),
							).
							WithFamilyName("1").
							WithGivenName("user").
							WithDisplayName("user 1").
							WithActive(true).
							Build(),
						UserBuilder().
							WithIPID("2").
							WithEmail(
								EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build(),
							).
							WithFamilyName("2").
							WithGivenName("user").
							WithDisplayName("user 2").
							WithActive(true).
							Build(),
					},
				).Build(),
			wantDelete: UsersResultBuilder().Build(),
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
						UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
					},
				).Build(),
				state: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
						UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("2").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
						UserBuilder().WithIPID("3").WithEmail(EmailBuilder().WithValue("user.3@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("3").WithGivenName("user").WithDisplayName("user 3").WithActive(true).Build(),
					},
				).Build(),
			},
			wantCreate: UsersResultBuilder().Build(),
			wantUpdate: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
				},
			).Build(),
			wantEqual: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
				},
			).Build(),
			wantDelete: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("3").WithEmail(EmailBuilder().WithValue("user.3@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("3").WithGivenName("user").WithDisplayName("user 3").WithActive(true).Build(),
				},
			).Build(),
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
						UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
					},
				).Build(),
				state: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
						UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("2").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
					},
				).Build(),
			},
			wantCreate: UsersResultBuilder().Build(),
			wantUpdate: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
				},
			).Build(),
			wantEqual: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
				},
			).Build(),
			wantDelete: UsersResultBuilder().Build(),
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: UsersResultBuilder().WithResources(
					[]*User{
						UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
						UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
						UserBuilder().WithIPID("4").WithEmail(EmailBuilder().WithValue("user.4@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("4").WithGivenName("user").WithDisplayName("user 4").WithActive(true).Build(),
					},
				).Build(),
				state: UsersResultBuilder().WithResources([]*User{
					UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
					UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("2").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
					UserBuilder().WithIPID("3").WithEmail(EmailBuilder().WithValue("user.3@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("3").WithGivenName("user").WithDisplayName("user 3").WithActive(true).Build(),
				},
				).Build(),
			},
			wantCreate: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("4").WithEmail(EmailBuilder().WithValue("user.4@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("4").WithGivenName("user").WithDisplayName("user 4").WithActive(true).Build(),
				},
			).Build(),
			wantUpdate: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("2").WithEmail(EmailBuilder().WithValue("user.2@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("different").WithGivenName("user").WithDisplayName("user 2").WithActive(true).Build(),
				},
			).Build(),
			wantEqual: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("1").WithEmail(EmailBuilder().WithValue("user.1@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("1").WithGivenName("user").WithDisplayName("user 1").WithActive(true).Build(),
				},
			).Build(),
			wantDelete: UsersResultBuilder().WithResources(
				[]*User{
					UserBuilder().WithIPID("3").WithEmail(EmailBuilder().WithValue("user.3@mail.com").WithType("Work").WithPrimary(true).Build()).WithFamilyName("3").WithGivenName("user").WithDisplayName("user 3").WithActive(true).Build(),
				},
			).Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.wantCreate.SetHashCode()
				tt.wantUpdate.SetHashCode()
				tt.wantEqual.SetHashCode()
				tt.wantDelete.SetHashCode()
			}

			gotCreate, gotUpdate, gotEqual, gotDelete, err := UsersOperations(tt.args.idp, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsOperations() error = %v", err)
				return
			}

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantCreate, gotCreate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotCreate (-tt.wantCreate +gotCreate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantUpdate, gotUpdate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotUpdate (-tt.wantUpdate +gotUpdate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantEqual, gotEqual, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotEqual (-tt.wantEqual +gotEqual):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantDelete, gotDelete, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotDelete (-tt.wantDelete +gotDelete):\n%s", diff)
			}
		})
	}
}

func TestMembersOperations(t *testing.T) {
	type args struct {
		idp  *GroupsMembersResult
		scim *GroupsMembersResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *GroupsMembersResult
		wantEqual  *GroupsMembersResult
		wantDelete *GroupsMembersResult
		wantErr    bool
	}{
		{
			name: "empty, return empty",
			args: args{
				idp:  GroupsMembersResultBuilder().Build(),
				scim: GroupsMembersResultBuilder().Build(),
			},
			wantCreate: GroupsMembersResultBuilder().Build(),
			wantEqual:  GroupsMembersResultBuilder().Build(),
			wantDelete: GroupsMembersResultBuilder().Build(),
			wantErr:    false,
		},
		{
			name: "nil idp, return error",
			args: args{
				idp:  nil,
				scim: GroupsMembersResultBuilder().Build(),
			},
			wantCreate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "nil scim, return error",
			args: args{
				idp:  GroupsMembersResultBuilder().Build(),
				scim: nil,
			},
			wantCreate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "one group: 1 add, 1 equal, 1 delete",
			args: args{
				idp: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						{
							Items: 2,
							Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
								MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
							},
						},
					},
				).Build(),
				scim: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						{
							Items: 2,
							Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
								MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
							},
						},
					},
				).Build(),
			},
			wantCreate: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantEqual: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantDelete: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantErr: false,
		},
		{
			name: "two groups: g1 -> add 1, g1 -> equal 1, g2 -> equal 1, g1 -> delete 1, g2 -> delete 1",
			args: args{
				idp: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
									MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
					},
				).Build(),
				scim: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
									MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
									MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
					},
				).Build(),
			},
			wantCreate: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantEqual: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantDelete: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantErr: false,
		},
		{
			name: "two groups: 2 equals, 1 add",
			args: args{
				idp: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
									MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "3", Name: "group 3", Email: "group.3@mail.com"}).
							WithResources(
								[]*Member{},
							).Build(),
					},
				).Build(),
				scim: GroupsMembersResultBuilder().WithResources(
					[]*GroupMembers{
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							).
							WithResources(
								[]*Member{
									MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
									MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
								},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							).
							WithResources(
								[]*Member{},
							).Build(),
						GroupMembersBuilder().
							WithGroup(
								&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"},
							).
							WithResources(
								[]*Member{},
							).Build(),
					},
				).Build(),
			},
			wantCreate: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
				},
			).Build(),
			wantEqual: GroupsMembersResultBuilder().WithResources(
				[]*GroupMembers{
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						).
						WithResources(
							[]*Member{
								MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
								MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
							},
						).Build(),
					GroupMembersBuilder().
						WithGroup(
							&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})}).
						WithResources(
							[]*Member{},
						).Build(),
				},
			).Build(),
			wantDelete: GroupsMembersResultBuilder().WithResources([]*GroupMembers{}).Build(),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.wantCreate.SetHashCode()
				tt.wantEqual.SetHashCode()
				tt.wantDelete.SetHashCode()
			}

			gotCreate, gotEqual, gotDelete, err := MembersOperations(tt.args.idp, tt.args.scim)

			if (err != nil) != tt.wantErr {
				t.Errorf("MembersOperations() error = %v", err)
				return
			}

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantCreate, gotCreate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotCreate (-tt.wantCreate +gotCreate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantEqual, gotEqual, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotEqual (-tt.wantEqual +gotEqual):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantDelete, gotDelete, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsOperations() gotDelete (-tt.wantDelete +gotDelete):\n%s", diff)
			}
		})
	}
}

func TestMergeGroupsResult(t *testing.T) {
	type args struct {
		grs []*GroupsResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged *GroupsResult
	}{
		{
			name: "merge empty",
			args: args{
				grs: []*GroupsResult{},
			},
			wantMerged: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
				HashCode:  "",
			},
		},
		{
			name: "nil arg",
			args: args{
				grs: nil,
			},
			wantMerged: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
				HashCode:  "",
			},
		},
		{
			name: "three groups",
			args: args{
				grs: []*GroupsResult{
					{
						Items: 1,
						Resources: []*Group{
							{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []*Group{
							{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: &GroupsResult{
				Items: 3,
				Resources: []*Group{
					{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			got := MergeGroupsResult(tt.args.grs...)

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantMerged, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("MergeGroupsResult() got (-tt.wantMerged +got):\n%s", diff)
			}
		})
	}
}

func TestMergeUsersResult(t *testing.T) {
	type args struct {
		urs []*UsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged *UsersResult
	}{
		{
			name: "merge empty",
			args: args{
				urs: []*UsersResult{},
			},
			wantMerged: &UsersResult{
				Items:     0,
				Resources: []*User{},
				HashCode:  "",
			},
		},
		{
			name: "nil arg",
			args: args{
				urs: nil,
			},
			wantMerged: &UsersResult{
				Items:     0,
				Resources: []*User{},
				HashCode:  "",
			},
		},
		{
			name: "three users",
			args: args{
				urs: []*UsersResult{
					{
						Items: 1,
						Resources: []*User{
							{IPID: "1", SCIMID: "1", Name: &Name{GivenName: "user", FamilyName: "1"}, Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}, HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []*User{
							{IPID: "2", SCIMID: "2", Name: &Name{GivenName: "user", FamilyName: "2"}, Emails: []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}, HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: &Name{GivenName: "user", FamilyName: "3"}, Emails: []Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}, HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: &UsersResult{
				Items: 3,
				Resources: []*User{
					{IPID: "1", SCIMID: "1", Name: &Name{GivenName: "user", FamilyName: "1"}, Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}, HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: &Name{GivenName: "user", FamilyName: "2"}, Emails: []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}, HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: &Name{GivenName: "user", FamilyName: "3"}, Emails: []Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}, HashCode: "1234509876"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			got := MergeUsersResult(tt.args.urs...)

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantMerged, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("MergeUsersResult() got (-tt.wantMerged +got):\n%s", diff)
			}
		})
	}
}

func TestMergeGroupsMembersResult(t *testing.T) {
	type args struct {
		gms []*GroupsMembersResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged *GroupsMembersResult
	}{
		{
			name: "empty",
			args: args{
				gms: []*GroupsMembersResult{},
			},
			wantMerged: &GroupsMembersResult{
				Items:     0,
				Resources: make([]*GroupMembers, 0),
			},
		},
		{
			name: "nil arg",
			args: args{
				gms: nil,
			},
			wantMerged: &GroupsMembersResult{
				Items:     0,
				Resources: make([]*GroupMembers, 0),
			},
		},
		{
			name: "two groups, two members each",
			args: args{
				gms: []*GroupsMembersResult{
					{
						Items:    1,
						HashCode: "123",
						Resources: []*GroupMembers{
							{
								Items: 2,
								Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321", Status: "ACTIVE"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*GroupMembers{
							{
								Items: 2,
								Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870", Status: "ACTIVE"},
								},
							},
						},
					},
				},
			},
			wantMerged: &GroupsMembersResult{
				Items:    2,
				HashCode: "123",
				Resources: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321", Status: "ACTIVE"},
						},
					},
					{
						Items: 2,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
							{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870", Status: "ACTIVE"},
						},
					},
				},
			},
		},
		{
			name: "three groups, one members each",
			args: args{
				gms: []*GroupsMembersResult{
					{
						Items:    1,
						HashCode: "123",
						Resources: []*GroupMembers{
							{
								Items: 1,
								Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*GroupMembers{
							{
								Items: 1,
								Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []*Member{
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321", Status: "ACTIVE"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*GroupMembers{
							{
								Items: 1,
								Group: &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
								Resources: []*Member{
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870", Status: "ACTIVE"},
								},
							},
						},
					},
				},
			},
			wantMerged: &GroupsMembersResult{
				Items:    3,
				HashCode: "123",
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890", Status: "ACTIVE"},
						},
					},
					{
						Items: 1,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []*Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321", Status: "ACTIVE"},
						},
					},
					{
						Items: 1,
						Group: &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870", Status: "ACTIVE"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			got := MergeGroupsMembersResult(tt.args.gms...)

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantMerged, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("MergeGroupsMembersResult() got (-tt.wantMerged +got):\n%s", diff)
			}
		})
	}
}

func TestMembersDataSets(t *testing.T) {
	type args struct {
		idp  []*GroupMembers
		scim []*GroupMembers
	}
	tests := []struct {
		name       string
		args       args
		wantCreate []*GroupMembers
		wantEqual  []*GroupMembers
		wantDelete []*GroupMembers
	}{
		{
			name: "empty return empty",
			args: args{
				idp:  make([]*GroupMembers, 0),
				scim: make([]*GroupMembers, 0),
			},
			wantCreate: make([]*GroupMembers, 0),
			wantEqual:  make([]*GroupMembers, 0),
			wantDelete: make([]*GroupMembers, 0),
		},
		{
			name: "one group: 1 add, 1 equal, 1 delete",
			args: args{
				idp: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},

			wantEqual: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},

			wantDelete: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},
		},
		{
			name: "two groups: g1 -> add 1, g1 -> equal 1, g2 -> equal 1, g1 -> delete 1, g2 -> delete 1",
			args: args{
				idp: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
					{
						Items: 1,
						Group: &Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
					{
						Items: 2,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},
			wantEqual: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
				{
					Items: 1,
					Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},
			wantDelete: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("3").WithSCIMID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
				{
					Items: 1,
					Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},
		},
		{
			name: "two groups: 2 equals, 1 add",
			args: args{
				idp: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
					{
						Items: 1,
						Group: &Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
					{
						Items:     0,
						Group:     &Group{IPID: "3", Name: "group 3", Email: "group.3@mail.com"},
						Resources: []*Member{},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					},
					{
						Items:     0,
						Group:     &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{},
					},
					{
						Items:     0,
						Group:     &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"},
						Resources: []*Member{},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})}, Resources: []*Member{
							MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
			},

			wantEqual: []*GroupMembers{
				{
					Items: 2,
					Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
						MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
					},
					HashCode: Hash(&GroupMembers{
						Items: 2,
						Group: &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							MemberBuilder().WithIPID("1").WithSCIMID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
							MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
						},
					}),
				},
				{
					Items:     0,
					Group:     &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
					Resources: []*Member{},
					HashCode: Hash(&GroupMembers{
						Items:     0,
						Group:     &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
						Resources: []*Member{},
					}),
				},
			},
			wantDelete: []*GroupMembers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, item := range tt.wantCreate {
				item.SetHashCode()
			}
			for _, item := range tt.wantEqual {
				item.SetHashCode()
			}
			for _, item := range tt.wantDelete {
				item.SetHashCode()
			}

			gotCreate, gotEqual, gotDelete := membersDataSets(tt.args.idp, tt.args.scim)

			sort := func(x, y string) bool { return x > y }

			if diff := cmp.Diff(tt.wantCreate, gotCreate, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("membersDataSets() gotCreate (-tt.wantCreate +gotCreate):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantEqual, gotEqual, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("membersDataSets() gotEqual (-tt.wantEqual +gotEqual):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantDelete, gotDelete, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("membersDataSets() gotDelete (-tt.wantDelete +gotDelete):\n%s", diff)
			}
		})
	}
}

func TestUpdateGroupsMembersSCIMID(t *testing.T) {
	t.Run("empty arguments", func(t *testing.T) {
		idp := &GroupsMembersResult{}
		scim := &GroupsResult{}
		scimUser := &UsersResult{}

		got := UpdateGroupsMembersSCIMID(idp, scim, scimUser)

		assert.Equal(t, 0, got.Items)
		assert.Equal(t, 0, len(got.Resources))
		assert.NotEqual(t, "", got.HashCode)
	})

	t.Run("update idp SCIM groups and users ids", func(t *testing.T) {
		idp := &GroupsMembersResult{
			Items: 2,
			Resources: []*GroupMembers{
				{
					Items: 2, Group: &Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
					Resources: []*Member{
						MemberBuilder().WithIPID("1").WithEmail("user.1@mail.com").WithStatus("ACTIVE").Build(),
						MemberBuilder().WithIPID("2").WithSCIMID("2").WithEmail("user.2@mail.com").WithStatus("ACTIVE").Build(),
					},
				},
				{
					Items: 3, Group: &Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
					Resources: []*Member{
						MemberBuilder().WithIPID("3").WithEmail("user.3@mail.com").WithStatus("ACTIVE").Build(),
						{IPID: "4", Email: "user.4@mail.com", Status: "ACTIVE"},
						{IPID: "5", Email: "user.5@mail.com", Status: "ACTIVE"},
					},
				},
			},
		}
		scim := &GroupsResult{
			Items: 2,
			Resources: []*Group{
				{SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
				{SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
			},
		}
		scimUser := &UsersResult{
			Items: 5,
			Resources: []*User{
				{SCIMID: "1", Name: &Name{GivenName: "user", FamilyName: "1"}, DisplayName: "user 1", Active: true, Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}},
				{SCIMID: "2", Name: &Name{GivenName: "user", FamilyName: "2"}, DisplayName: "user 2", Active: true, Emails: []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}},
				{SCIMID: "3", Name: &Name{GivenName: "user", FamilyName: "3"}, DisplayName: "user 3", Active: true, Emails: []Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}},
				{SCIMID: "4", Name: &Name{GivenName: "user", FamilyName: "4"}, DisplayName: "user 4", Active: true, Emails: []Email{{Value: "user.4@mail.com", Type: "work", Primary: true}}},
				{SCIMID: "5", Name: &Name{GivenName: "user", FamilyName: "5"}, DisplayName: "user 5", Active: true, Emails: []Email{{Value: "user.5@mail.com", Type: "work", Primary: true}}},
			},
		}

		got := UpdateGroupsMembersSCIMID(idp, scim, scimUser)

		assert.Equal(t, 2, got.Items)
		assert.Equal(t, 2, len(got.Resources))
		assert.NotEqual(t, "", got.HashCode)

		assert.Equal(t, 2, got.Resources[0].Items)
		assert.Equal(t, 3, got.Resources[1].Items)

		assert.Equal(t, 2, len(got.Resources[0].Resources))
		assert.Equal(t, 3, len(got.Resources[1].Resources))

		assert.Equal(t, "1", got.Resources[0].Group.SCIMID)
		assert.Equal(t, "2", got.Resources[1].Group.SCIMID)

		assert.Equal(t, "1", got.Resources[0].Group.IPID)
		assert.Equal(t, "2", got.Resources[1].Group.IPID)

		assert.Equal(t, "group.1@mail.com", got.Resources[0].Group.Email)
		assert.Equal(t, "group.2@mail.com", got.Resources[1].Group.Email)

		assert.Equal(t, "1", got.Resources[0].Resources[0].SCIMID)
		assert.Equal(t, "2", got.Resources[0].Resources[1].SCIMID)
		assert.Equal(t, "3", got.Resources[1].Resources[0].SCIMID)
		assert.Equal(t, "4", got.Resources[1].Resources[1].SCIMID)
		assert.Equal(t, "5", got.Resources[1].Resources[2].SCIMID)

		assert.Equal(t, "1", got.Resources[0].Resources[0].IPID)
		assert.Equal(t, "2", got.Resources[0].Resources[1].IPID)
		assert.Equal(t, "3", got.Resources[1].Resources[0].IPID)
		assert.Equal(t, "4", got.Resources[1].Resources[1].IPID)
		assert.Equal(t, "5", got.Resources[1].Resources[2].IPID)

		assert.Equal(t, "user.1@mail.com", got.Resources[0].Resources[0].Email)
		assert.Equal(t, "user.2@mail.com", got.Resources[0].Resources[1].Email)
		assert.Equal(t, "user.3@mail.com", got.Resources[1].Resources[0].Email)
		assert.Equal(t, "user.4@mail.com", got.Resources[1].Resources[1].Email)
		assert.Equal(t, "user.5@mail.com", got.Resources[1].Resources[2].Email)

		assert.Equal(t, "ACTIVE", got.Resources[0].Resources[0].Status)
		assert.Equal(t, "ACTIVE", got.Resources[0].Resources[1].Status)
		assert.Equal(t, "ACTIVE", got.Resources[1].Resources[0].Status)
		assert.Equal(t, "ACTIVE", got.Resources[1].Resources[1].Status)
		assert.Equal(t, "ACTIVE", got.Resources[1].Resources[2].Status)
	})
}
