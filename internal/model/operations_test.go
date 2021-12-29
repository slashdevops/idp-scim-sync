package model

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
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
				idp: &GroupsResult{
					Items:     0,
					Resources: []*Group{},
				},
				state: &GroupsResult{
					Items:     0,
					Resources: []*Group{},
				},
			},
			wantCreate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantUpdate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantEqual: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantErr: false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: &GroupsResult{},
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
				idp:   &GroupsResult{},
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
				idp: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantUpdate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantEqual: &GroupsResult{
				Items: 2,
				Resources: []*Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
					{IPID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantUpdate: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 3,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantUpdate: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &GroupsResult{
					Items: 4,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "dd", Name: "name2", Email: "2@mail.com"},
						{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 3,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
				},
			},
			wantUpdate: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: &GroupsResult{
					Items: 1,
					Resources: []*Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 1,
					Resources: []*Group{
						{IPID: "3", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
					},
				},
			},
			wantCreate: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantUpdate: &GroupsResult{
				Items: 1,
				Resources: []*Group{
					{IPID: "1", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantEqual: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
			wantErr: false,
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
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("GroupsOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("GroupsOperations() gotUpdate = %s, want %s", utils.ToJSON(gotUpdate), utils.ToJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("GroupsOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("GroupsOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
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
				idp: &UsersResult{
					Items:     0,
					Resources: []*User{},
				},
				state: &UsersResult{
					Items:     0,
					Resources: []*User{},
				},
			},
			wantCreate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantUpdate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantEqual: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantDelete: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantErr: false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: &UsersResult{},
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
				idp:   &UsersResult{},
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
				idp: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
				state: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantUpdate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantEqual: &UsersResult{
				Items: 2,
				Resources: []*User{
					{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
					{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
				},
			},
			wantDelete: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantUpdate: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
				},
			},
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
			wantUpdate: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &UsersResult{
				Items:     0,
				Resources: []*User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
						{IPID: "4", Email: "don.nadie@email.com", Name: Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
					},
				},
				state: &UsersResult{
					Items: 2,
					Resources: []*User{
						{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "4", Email: "don.nadie@email.com", Name: Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
				},
			},
			wantUpdate: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "2", Email: "foo.bar@email.com", Name: Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "1", Email: "john.doe@email.com", Name: Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &UsersResult{
				Items: 1,
				Resources: []*User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
				},
			},
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
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("UsersOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("UsersOperations() gotUpdate = %s, want %s", utils.ToJSON(gotUpdate), utils.ToJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("UsersOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("UsersOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
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
				idp:  &GroupsMembersResult{},
				scim: &GroupsMembersResult{},
			},
			wantCreate: &GroupsMembersResult{
				Items:     0,
				Resources: []*GroupMembers{},
			},
			wantEqual: &GroupsMembersResult{
				Items:     0,
				Resources: []*GroupMembers{},
			},
			wantDelete: &GroupsMembersResult{
				Items:     0,
				Resources: []*GroupMembers{},
			},
			wantErr: false,
		},
		{
			name: "nil idp, return error",
			args: args{
				idp:  nil,
				scim: &GroupsMembersResult{},
			},
			wantCreate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "nil scim, return error",
			args: args{
				idp:  &GroupsMembersResult{},
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
				idp: &GroupsMembersResult{
					Items: 1,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
					},
				},
				scim: &GroupsMembersResult{
					Items: 1,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &GroupsMembersResult{
				Items: 1,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "2", Email: "user.2@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &GroupsMembersResult{
				Items: 1,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							},
						}),
					},
				},
			},
			wantDelete: &GroupsMembersResult{
				Items: 1,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "two groups: g1 -> add 1, g1 -> equal 1, g2 -> equal 1, g1 -> delete 1, g2 -> delete 1",
			args: args{
				idp: &GroupsMembersResult{
					Items: 2,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []*Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
				scim: &GroupsMembersResult{
					Items: 2,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
						{
							Items: 2,
							Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &GroupsMembersResult{
				Items: 1,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", Email: "user.1@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "1", Email: "user.1@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &GroupsMembersResult{
				Items: 2,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
					{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
							Resources: []*Member{
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
				},
			},
			wantDelete: &GroupsMembersResult{
				Items: 2,
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
					{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							},
						}),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "two groups: 2 equals, 1 add",
			args: args{
				idp: &GroupsMembersResult{
					Items: 2,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []*Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						},
						{
							Items:     0,
							Group:     Group{IPID: "3", Name: "group 3", Email: "group.3@mail.com"},
							Resources: []*Member{},
						},
					},
				},
				scim: &GroupsMembersResult{
					Items: 2,
					Resources: []*GroupMembers{
						{
							Items: 2,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items:     0,
							Group:     Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []*Member{},
						},
						{
							Items:     0,
							Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"},
							Resources: []*Member{},
						},
					},
				},
			},
			wantCreate: &GroupsMembersResult{
				Items:    1,
				HashCode: "",
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							{IPID: "3", Email: "user.3@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 1,
							Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})}, Resources: []*Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &GroupsMembersResult{
				Items: 2,
				Resources: []*GroupMembers{
					{
						Items: 2,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
						HashCode: Hash(&GroupMembers{
							Items: 2,
							Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []*Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
					{
						Items:     0,
						Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
						Resources: []*Member{},
						HashCode: Hash(&GroupMembers{
							Items:     0,
							Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
							Resources: []*Member{},
						}),
					},
				},
			},
			wantDelete: &GroupsMembersResult{
				Items:     0,
				HashCode:  "",
				Resources: []*GroupMembers{},
			},
			wantErr: false,
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
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("MembersOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("MembersOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("MembersOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
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

			if gotMerged := MergeGroupsResult(tt.args.grs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("MergeGroupsResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
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
							{IPID: "1", SCIMID: "1", Name: Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []*User{
							{IPID: "2", SCIMID: "2", Name: Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: &UsersResult{
				Items: 3,
				Resources: []*User{
					{IPID: "1", SCIMID: "1", Name: Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			if gotMerged := MergeUsersResult(tt.args.urs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("MergeUsersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
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
								Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
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
								Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
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
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
						},
					},
					{
						Items: 2,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
							{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
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
								Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []*Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
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
								Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []*Member{
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
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
								Group: Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
								Resources: []*Member{
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
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
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
					},
					{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []*Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
						},
					},
					{
						Items: 1,
						Group: Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			if gotMerged := MergeGroupsMembersResult(tt.args.gms...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("MergeGroupsMembersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
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
						Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "1", Email: "user.1@mail.com"},
							{IPID: "2", Email: "user.2@mail.com"},
						},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "2", Email: "user.2@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "2", Email: "user.2@mail.com"},
						},
					}),
				},
			},

			wantEqual: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						},
					}),
				},
			},

			wantDelete: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
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
						Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "1", Email: "user.1@mail.com"},
							{IPID: "2", Email: "user.2@mail.com"},
						},
					},
					{
						Items: 1,
						Group: Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							{IPID: "3", Email: "user.3@mail.com"},
						},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
					},
					{
						Items: 2,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "1", Email: "user.1@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", Email: "user.1@mail.com"},
						},
					}),
				},
			},
			wantEqual: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
					}),
				},
				{
					Items: 1,
					Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
					}),
				},
			},
			wantDelete: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
					}),
				},
				{
					Items: 1,
					Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
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
						Group: Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "1", Email: "user.1@mail.com"},
							{IPID: "2", Email: "user.2@mail.com"},
						},
					},
					{
						Items: 1,
						Group: Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{
							{IPID: "3", Email: "user.3@mail.com"},
						},
					},
					{
						Items:     0,
						Group:     Group{IPID: "3", Name: "group 3", Email: "group.3@mail.com"},
						Resources: []*Member{},
					},
				},
				scim: []*GroupMembers{
					{
						Items: 2,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
					},
					{
						Items:     0,
						Group:     Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
						Resources: []*Member{},
					},
					{
						Items:     0,
						Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"},
						Resources: []*Member{},
					},
				},
			},
			wantCreate: []*GroupMembers{
				{
					Items: 1,
					Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
					Resources: []*Member{
						{IPID: "3", Email: "user.3@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 1,
						Group: Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: Hash(&Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})}, Resources: []*Member{
							{IPID: "3", Email: "user.3@mail.com"},
						},
					}),
				},
			},

			wantEqual: []*GroupMembers{
				{
					Items: 2,
					Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
					Resources: []*Member{
						{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
					},
					HashCode: Hash(&GroupMembers{
						Items: 2,
						Group: Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: Hash(&Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
					}),
				},
				{
					Items:     0,
					Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
					Resources: []*Member{},
					HashCode: Hash(&GroupMembers{
						Items:     0,
						Group:     Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: Hash(&Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
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

			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("membersDataSets() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("membersDataSets() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("membersDataSets() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}
