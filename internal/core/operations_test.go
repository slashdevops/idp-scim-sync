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
	}{
		{
			name: "empty",
			args: args{
				idp: &model.GroupsResult{
					Items:     0,
					Resources: []model.Group{},
				},
				state: &model.GroupsResult{
					Items:     0,
					Resources: []model.Group{},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantEqual: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
		},
		{
			name: "2 equals",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantEqual: &model.GroupsResult{
				Items: 2,
				Resources: []model.Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
					{IPID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.GroupsResult{
					Items: 4,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", SCIMID: "dd", Name: "name2", Email: "2@mail.com"},
						{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
						{IPID: "", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "4", SCIMID: "44", Name: "name4", Email: "4@mail.com"},
				},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "2", SCIMID: "22", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "3", SCIMID: "33", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: &model.GroupsResult{
					Items: 1,
					Resources: []model.Group{
						{IPID: "1", SCIMID: "11", Name: "name1", Email: "1@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 1,
					Resources: []model.Group{
						{IPID: "3", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", SCIMID: "22", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []model.Group{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreate, gotUpdate, gotEqual, gotDelete := groupsOperations(tt.args.idp, tt.args.state)
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
	}{
		{
			name: "empty",
			args: args{
				idp: &model.UsersResult{
					Items:     0,
					Resources: []model.User{},
				},
				state: &model.UsersResult{
					Items:     0,
					Resources: []model.User{},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantEqual: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
		},
		{
			name: "2 equals",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantEqual: &model.UsersResult{
				Items: 2,
				Resources: []model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
				},
			},
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items:     0,
				Resources: []model.User{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
						{IPID: "4", Email: "don.nadie@email.com", Name: model.Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
					},
				},
				state: &model.UsersResult{
					Items: 2,
					Resources: []model.User{
						{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
						{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foo", GivenName: "bar"}, DisplayName: "foo bar", Active: true},
						{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
					},
				},
			},
			wantCreate: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "4", Email: "don.nadie@email.com", Name: model.Name{FamilyName: "don", GivenName: "nadie"}, DisplayName: "don nadie", Active: true},
				},
			},
			wantUpdate: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "2", Email: "foo.bar@email.com", Name: model.Name{FamilyName: "foonew", GivenName: "bar"}, DisplayName: "foonew bar", Active: true},
				},
			},
			wantEqual: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "1", Email: "john.doe@email.com", Name: model.Name{FamilyName: "john", GivenName: "doe"}, DisplayName: "john doe", Active: true},
				},
			},
			wantDelete: &model.UsersResult{
				Items: 1,
				Resources: []model.User{
					{IPID: "3", Email: "donato.ricupero@email.com", Name: model.Name{FamilyName: "donato", GivenName: "ricupero"}, DisplayName: "donato ricupero", Active: true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreate, gotUpdate, gotEqual, gotDelete := usersOperations(tt.args.idp, tt.args.state)
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

func TestGroupsUsersOperations(t *testing.T) {
	type args struct {
		idp   *model.GroupsUsersResult
		state *model.GroupsUsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantCreate *model.GroupsUsersResult
		wantEqual  *model.GroupsUsersResult
		wantDelete *model.GroupsUsersResult
	}{
		{
			name: "empty",
			args: args{
				idp: &model.GroupsUsersResult{
					Items:     0,
					Resources: []model.GroupUsers{},
				},
				state: &model.GroupsUsersResult{
					Items:     0,
					Resources: []model.GroupUsers{},
				},
			},
			wantCreate: &model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
			},
			wantEqual: &model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
			},
			wantDelete: &model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
			},
		},
		{
			name: "2 equals",
			args: args{
				idp: &model.GroupsUsersResult{
					Items: 2,
					Resources: []model.GroupUsers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
								{IPID: "2", Name: model.Name{FamilyName: "user", GivenName: "2"}, Email: "u.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
							},
						},
					},
				},
				state: &model.GroupsUsersResult{
					Items: 2,
					Resources: []model.GroupUsers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
								{IPID: "2", Name: model.Name{FamilyName: "user", GivenName: "2"}, Email: "u.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
			},
			wantEqual: &model.GroupsUsersResult{
				Items: 2,
				Resources: []model.GroupUsers{
					{
						Items: 2,
						Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
						Resources: []model.User{
							{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
							{IPID: "2", Name: model.Name{FamilyName: "user", GivenName: "2"}, Email: "u.2@mail.com"},
						},
					},
					{
						Items: 1,
						Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
						Resources: []model.User{
							{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
						},
					},
				},
			},
			wantDelete: &model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
			},
		},
		{
			name: "1 equals, 1 create, 1 delete",
			args: args{
				idp: &model.GroupsUsersResult{
					Items: 2,
					Resources: []model.GroupUsers{
						{
							Items: 2,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
								{IPID: "2", Name: model.Name{FamilyName: "user", GivenName: "2"}, Email: "u.2@mail.com"},
							},
						},
						{
							Items: 1,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
							},
						},
					},
				},
				state: &model.GroupsUsersResult{
					Items: 2,
					Resources: []model.GroupUsers{
						{
							Items: 1,
							Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
							},
						},
						{
							Items: 2,
							Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
							Resources: []model.User{
								{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
								{IPID: "3", Name: model.Name{FamilyName: "user", GivenName: "3"}, Email: "u.3@mail.com"},
							},
						},
					},
				},
			},
			wantCreate: &model.GroupsUsersResult{
				Items: 1,
				Resources: []model.GroupUsers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
						Resources: []model.User{
							{IPID: "2", Name: model.Name{FamilyName: "user", GivenName: "2"}, Email: "u.2@mail.com"},
						},
					},
				},
			},
			wantEqual: &model.GroupsUsersResult{
				Items: 2,
				Resources: []model.GroupUsers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", Name: "group 1", Email: "g.1@mail.com"},
						Resources: []model.User{
							{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
						},
					},
					{
						Items: 1,
						Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
						Resources: []model.User{
							{IPID: "1", Name: model.Name{FamilyName: "user", GivenName: "1"}, Email: "u.1@mail.com"},
						},
					},
				},
			},
			wantDelete: &model.GroupsUsersResult{
				Items: 1,
				Resources: []model.GroupUsers{
					{
						Items: 1,
						Group: model.Group{IPID: "2", Name: "group 2", Email: "g.2@mail.com"},
						Resources: []model.User{
							{IPID: "3", Name: model.Name{FamilyName: "user", GivenName: "3"}, Email: "u.3@mail.com"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreate, gotEqual, gotDelete := groupsUsersOperations(tt.args.idp, tt.args.state)
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("groupsUsersOperations() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("groupsUsersOperations() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("groupsUsersOperations() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}

func Test_mergeGroupsResult(t *testing.T) {
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
				Resources: []model.Group{},
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
				Resources: []model.Group{},
				HashCode:  "",
			},
		},
		{
			name: "three groups",
			args: args{
				grs: []*model.GroupsResult{
					{
						Items: 1,
						Resources: []model.Group{
							{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []model.Group{
							{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: model.GroupsResult{
				Items: 3,
				Resources: []model.Group{
					{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
				},
				HashCode: hash.Get(
					model.GroupsResult{
						Items: 3,
						Resources: []model.Group{
							{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@gmail.com", HashCode: "1234567890"},
							{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@gmail.com", HashCode: "1234509876"},
						},
					},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMerged := mergeGroupsResult(tt.args.grs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeGroupsResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}

func Test_mergeUsersResult(t *testing.T) {
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
				Resources: []model.User{},
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
				Resources: []model.User{},
				HashCode:  "",
			},
		},
		{
			name: "three users",
			args: args{
				urs: []*model.UsersResult{
					{
						Items: 1,
						Resources: []model.User{
							{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "1234",
					},
					{
						Items: 2,
						Resources: []model.User{
							{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "1234",
					},
				},
			},
			wantMerged: model.UsersResult{
				Items: 3,
				Resources: []model.User{
					{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
					{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
					{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
				},
				HashCode: hash.Get(
					model.UsersResult{
						Items: 3,
						Resources: []model.User{
							{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
							{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
						},
					},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMerged := mergeUsersResult(tt.args.urs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeUsersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}

func Test_mergeGroupsUsersResult(t *testing.T) {
	type args struct {
		gurs []*model.GroupsUsersResult
	}
	tests := []struct {
		name       string
		args       args
		wantMerged model.GroupsUsersResult
	}{
		{
			name: "merge empty",
			args: args{
				gurs: []*model.GroupsUsersResult{},
			},
			wantMerged: model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
				HashCode:  "",
			},
		},
		{
			name: "nil args",
			args: args{
				gurs: nil,
			},
			wantMerged: model.GroupsUsersResult{
				Items:     0,
				Resources: []model.GroupUsers{},
				HashCode:  "",
			},
		},
		{
			name: "2 groups, three users",
			args: args{
				gurs: []*model.GroupsUsersResult{
					{
						Items: 1,
						Resources: []model.GroupUsers{
							{
								Items: 1,
								Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", HashCode: "1234567890"},
								Resources: []model.User{
									{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
								},
								HashCode: "123",
							},
						},
						HashCode: "456",
					},
					{
						Items: 2,
						Resources: []model.GroupUsers{
							{
								Items: 2,
								Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", HashCode: "9876543210"},
								Resources: []model.User{
									{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
									{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
								},
								HashCode: "321",
							},
						},
						HashCode: "654",
					},
				},
			},
			wantMerged: model.GroupsUsersResult{
				Items: 2,
				Resources: []model.GroupUsers{
					{
						Items: 1,
						Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", HashCode: "1234567890"},
						Resources: []model.User{
							{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
						},
						HashCode: "123",
					},
					{
						Items: 2,
						Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", HashCode: "9876543210"},
						Resources: []model.User{
							{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
							{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
						},
						HashCode: "321",
					},
				},
				HashCode: hash.Get(
					model.GroupsUsersResult{
						Items: 2,
						Resources: []model.GroupUsers{
							{
								Items: 1,
								Group: model.Group{IPID: "1", SCIMID: "1", Name: "group 1", HashCode: "1234567890"},
								Resources: []model.User{
									{IPID: "1", SCIMID: "1", Name: model.Name{GivenName: "user", FamilyName: "1"}, Email: "user.1@gmail.com", HashCode: "1234567890"},
								},
								HashCode: "123",
							},
							{
								Items: 2,
								Group: model.Group{IPID: "2", SCIMID: "2", Name: "group 2", HashCode: "9876543210"},
								Resources: []model.User{
									{IPID: "2", SCIMID: "2", Name: model.Name{GivenName: "user", FamilyName: "2"}, Email: "user.2@gmail.com", HashCode: "0987654321"},
									{IPID: "3", SCIMID: "3", Name: model.Name{GivenName: "user", FamilyName: "3"}, Email: "user.3@gmail.com", HashCode: "1234509876"},
								},
								HashCode: "321",
							},
						},
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMerged := mergeGroupsUsersResult(tt.args.gurs...); !reflect.DeepEqual(gotMerged, tt.wantMerged) {
				t.Errorf("mergeGroupsUsersResult() = %s, want %s", utils.ToJSON(gotMerged), utils.ToJSON(tt.wantMerged))
			}
		})
	}
}
