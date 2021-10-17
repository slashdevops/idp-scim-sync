package core

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func Test_groupsOperations(t *testing.T) {
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
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "22@mail.com"},
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
				Items: 1,
				Resources: []model.Group{
					{IPID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
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
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "22@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", Name: "name3", Email: "3@mail.com"},
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
					{IPID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "3", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.GroupsResult{
					Items: 4,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "22@mail.com"},
						{IPID: "4", Name: "name4", Email: "4@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
						{IPID: "2", Name: "name2", Email: "2@mail.com"},
						{IPID: "3", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "4", Name: "name4", Email: "4@mail.com"},
				},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "1", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []model.Group{
					{IPID: "3", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: &model.GroupsResult{
					Items: 1,
					Resources: []model.Group{
						{IPID: "11", Name: "name1", Email: "1@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 1,
					Resources: []model.Group{
						{IPID: "1", Name: "name1", Email: "1@mail.com"},
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
					{IPID: "11", Name: "name1", Email: "1@mail.com"},
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

func Test_usersOperations(t *testing.T) {
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

func Test_groupsUsersOperations(t *testing.T) {
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
