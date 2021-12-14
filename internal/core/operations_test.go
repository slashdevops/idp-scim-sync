package core

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func TestGroupsOperations(t *testing.T) {
	type args struct {
		idp   *model.GroupsResult
		state *model.GroupsResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *model.GroupsResult
		wantUpdate *model.GroupsResult
		wantEqual  *model.GroupsResult
		wantDelete *model.GroupsResult
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				idp: &model.GroupsResult{
					Items:     0,
					Resources: []*model.Group{},
				},
				state: &model.GroupsResult{
					Items:     0,
					Resources: []*model.Group{},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantEqual: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantErr: false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: &model.GroupsResult{},
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
				idp:   &model.GroupsResult{},
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
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantEqual: &model.GroupsResult{
				Items: 2,
				Resources: []*model.Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
					{IPID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.GroupsResult{
					Items: 4,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "dd", Name: "name2", Email: "2@mail.com"},
						{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
				},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{IPID: "3", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{IPID: "1", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
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

			gotCreate, gotUpdate, gotEqual, gotDelete, err := groupsOperations(tt.args.idp, tt.args.state)

			if (err != nil) != tt.wantErr {
				t.Errorf("groupsOperations() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("groupsOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("groupsOperations() gotUpdate = %s, want %s", utils.ToJSON(gotUpdate), utils.ToJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("groupsOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("groupsOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}

func TestUsersOperations(t *testing.T) {
	type args struct {
		idp   *model.UsersResult
		state *model.UsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *model.UsersResult
		wantUpdate *model.UsersResult
		wantEqual  *model.UsersResult
		wantDelete *model.UsersResult
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				idp: &model.UsersResult{
					Items:     0,
					Resources: []*model.User{},
				},
				state: &model.UsersResult{
					Items:     0,
					Resources: []*model.User{},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantEqual: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantErr: false,
		},
		{
			name: "nil idp",
			args: args{
				idp:   nil,
				state: &model.UsersResult{},
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
				idp:   &model.UsersResult{},
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
				idp: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantEqual: &model.UsersResult{
				Items: 2,
				Resources: []*model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
				},
			},
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
						{IPID: "4", Email: "don.nadie@email.com", Name: model.Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []*model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "4", Email: "don.nadie@email.com", Name: model.Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
				},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items: 1,
				Resources: []*model.User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
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

			gotCreate, gotUpdate, gotEqual, gotDelete, err := usersOperations(tt.args.idp, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupsOperations() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("usersOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("usersOperations() gotUpdate = %s, want %s", utils.ToJSON(gotUpdate), utils.ToJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("usersOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("usersOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}

func TestMembersOperations(t *testing.T) {
	type args struct {
		idp  *model.GroupsMembersResult
		scim *model.GroupsMembersResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *model.GroupsMembersResult
		wantEqual  *model.GroupsMembersResult
		wantDelete *model.GroupsMembersResult
		wantErr    bool
	}{
		{
			name: "empty, return empty",
			args: args{
				idp:  &model.GroupsMembersResult{},
				scim: &model.GroupsMembersResult{},
			},
			wantCreate: &model.GroupsMembersResult{
				Items:     0,
				Resources: []*model.GroupMembers{},
			},
			wantEqual: &model.GroupsMembersResult{
				Items:     0,
				Resources: []*model.GroupMembers{},
			},
			wantDelete: &model.GroupsMembersResult{
				Items:     0,
				Resources: []*model.GroupMembers{},
			},
			wantErr: false,
		},
		{
			name: "nil idp, return error",
			args: args{
				idp:  nil,
				scim: &model.GroupsMembersResult{},
			},
			wantCreate: nil,
			wantEqual:  nil,
			wantDelete: nil,
			wantErr:    true,
		},
		{
			name: "nil scim, return error",
			args: args{
				idp:  &model.GroupsMembersResult{},
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
				idp: &model.GroupsMembersResult{
					Items: 1,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
					},
				},
				scim: &model.GroupsMembersResult{
					Items: 1,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &model.GroupsMembersResult{
				Items: 1,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "2", Email: "user.2@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &model.GroupsMembersResult{
				Items: 1,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							},
						}),
					},
				},
			},
			wantDelete: &model.GroupsMembersResult{
				Items: 1,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
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
				idp: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []model.Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
				scim: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
						{
							Items: 2,
							Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []model.Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &model.GroupsMembersResult{
				Items: 1,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "1", Email: "user.1@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "1", Email: "user.1@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &model.GroupsMembersResult{
				Items: 2,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
					{
						Items: 1,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []model.Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
							Resources: []model.Member{
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
				},
			},
			wantDelete: &model.GroupsMembersResult{
				Items: 2,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "3", SCIMID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
					{
						Items: 1,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []model.Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
							Resources: []model.Member{
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
				idp: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "1", Email: "user.1@mail.com"},
								{IPID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []model.Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						},
						{
							Items:     0,
							Group:     model.Group{IPID: "3", Name: "group 3", Email: "group.3@mail.com"},
							Resources: []model.Member{},
						},
					},
				},
				scim: &model.GroupsMembersResult{
					Items: 2,
					Resources: []*model.GroupMembers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"},
							Resources: []model.Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						},
						{
							Items:     0,
							Group:     model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"},
							Resources: []model.Member{},
						},
						{
							Items:     0,
							Group:     model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"},
							Resources: []model.Member{},
						},
					},
				},
			},
			wantCreate: &model.GroupsMembersResult{
				Items:    1,
				HashCode: "",
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})},
						Resources: []model.Member{
							{IPID: "3", Email: "user.3@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 1,
							Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: hash.Get(&model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"})}, Resources: []model.Member{
								{IPID: "3", Email: "user.3@mail.com"},
							},
						}),
					},
				},
			},
			wantEqual: &model.GroupsMembersResult{
				Items: 2,
				Resources: []*model.GroupMembers{
					{
						Items: 2,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
						Resources: []model.Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
						},
						HashCode: hash.Get(&model.GroupMembers{
							Items: 2,
							Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: hash.Get(&model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"})},
							Resources: []model.Member{
								{IPID: "1", SCIMID: "1", Email: "user.1@mail.com"},
								{IPID: "2", SCIMID: "2", Email: "user.2@mail.com"},
							},
						}),
					},
					{
						Items:     0,
						Group:     model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: hash.Get(&model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
						Resources: []model.Member{},
						HashCode: hash.Get(&model.GroupMembers{
							Items:     0,
							Group:     model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: hash.Get(&model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"})},
							Resources: []model.Member{},
						}),
					},
				},
			},
			wantDelete: &model.GroupsMembersResult{
				Items:     0,
				HashCode:  "",
				Resources: []*model.GroupMembers{},
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

			gotCreate, gotEqual, gotDelete, err := membersOperations(tt.args.idp, tt.args.scim)

			if (err != nil) != tt.wantErr {
				t.Errorf("membersOperations() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("membersOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("membersOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("membersOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}

func TestMergeGroupsResult(t *testing.T) {
	type args struct {
		grs []*model.GroupsResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged model.GroupsResult
	}{
		{
			name: "merge empty",
			args: args{
				grs: []*model.GroupsResult{},
			},
			wantMerged: model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
				HashCode:  "",
			},
		},
		{
			name: "nil arg",
			args: args{
				grs: nil,
			},
			wantMerged: model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
				HashCode:  "",
			},
		},
		{
			name: "three groups",
			args: args{
				grs: []*model.GroupsResult{
					{
						Items: 1,
						Resources: []*model.Group{
							{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []*model.Group{
							{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: model.GroupsResult{
				Items: 3,
				Resources: []*model.Group{
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

			if gotMerged := mergeGroupsResult(tt.args.grs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeGroupsResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}

func TestMergeUsersResult(t *testing.T) {
	type args struct {
		urs []*model.UsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged model.UsersResult
	}{
		{
			name: "merge empty",
			args: args{
				urs: []*model.UsersResult{},
			},
			wantMerged: model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
				HashCode:  "",
			},
		},
		{
			name: "nil arg",
			args: args{
				urs: nil,
			},
			wantMerged: model.UsersResult{
				Items:     0,
				Resources: []*model.User{},
				HashCode:  "",
			},
		},
		{
			name: "three users",
			args: args{
				urs: []*model.UsersResult{
					{
						Items: 1,
						Resources: []*model.User{
							{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []*model.User{
							{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: model.UsersResult{
				Items: 3,
				Resources: []*model.User{
					{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantMerged.SetHashCode()

			if gotMerged := mergeUsersResult(tt.args.urs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeUsersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}

func TestMergeGroupsMembersResult(t *testing.T) {
	type args struct {
		gms []*model.GroupsMembersResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged model.GroupsMembersResult
	}{
		{
			name: "empty",
			args: args{
				gms: []*model.GroupsMembersResult{},
			},
			wantMerged: model.GroupsMembersResult{
				Items:     0,
				Resources: []*model.GroupMembers{},
				HashCode:  "",
			},
		},
		{
			name: "nil arg",
			args: args{
				gms: nil,
			},
			wantMerged: model.GroupsMembersResult{
				Items:     0,
				Resources: []*model.GroupMembers{},
				HashCode:  "",
			},
		},
		{
			name: "two groups, two members each",
			args: args{
				gms: []*model.GroupsMembersResult{
					{
						Items:    1,
						HashCode: "123",
						Resources: []*model.GroupMembers{
							{
								Items: 2,
								Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []model.Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*model.GroupMembers{
							{
								Items: 2,
								Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []model.Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
								},
							},
						},
					},
				},
			},
			wantMerged: model.GroupsMembersResult{
				Items: 2,
				Resources: []*model.GroupMembers{
					{
						Items: 2,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []model.Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
						},
					},
					{
						Items: 2,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []model.Member{
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
				gms: []*model.GroupsMembersResult{
					{
						Items:    1,
						HashCode: "123",
						Resources: []*model.GroupMembers{
							{
								Items: 1,
								Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
								Resources: []model.Member{
									{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*model.GroupMembers{
							{
								Items: 1,
								Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
								Resources: []model.Member{
									{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
								},
							},
						},
					},
					{
						Items:    1,
						HashCode: "321",
						Resources: []*model.GroupMembers{
							{
								Items: 1,
								Group: model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
								Resources: []model.Member{
									{IPID: "3", SCIMID: "3", Email: "user.3@gmail.com", HashCode: "5612309870"},
								},
							},
						},
					},
				},
			},
			wantMerged: model.GroupsMembersResult{
				Items: 3,
				Resources: []*model.GroupMembers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com", HashCode: "1234567890"},
						Resources: []model.Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
					},
					{
						Items: 1,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com", HashCode: "0987654321"},
						Resources: []model.Member{
							{IPID: "2", SCIMID: "2", Email: "user.2@gmail.com", HashCode: "0987654321"},
						},
					},
					{
						Items: 1,
						Group: model.Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com", HashCode: "6543219870"},
						Resources: []model.Member{
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

			if gotMerged := mergeGroupsMembersResult(tt.args.gms...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeGroupsMembersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}
