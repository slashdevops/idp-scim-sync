package core

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func Test_groupsDifferences(t *testing.T) {
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
		},
		{
			name: "2 equals",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
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
					{ID: "1", Name: "name1", Email: "1@mail.com"},
					{ID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
		},
		{
			name: "1 equals, 1 update",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "22@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
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
					{ID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "1", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items:     0,
				Resources: []*model.Group{},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete",
			args: args{
				idp: &model.GroupsResult{
					Items: 2,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "22@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
						{ID: "3", Name: "name3", Email: "3@mail.com"},
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
					{ID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "1", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "3", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 equals, 1 update, 1 delete, 1 create",
			args: args{
				idp: &model.GroupsResult{
					Items: 4,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "22@mail.com"},
						{ID: "4", Name: "name4", Email: "4@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 3,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
						{ID: "3", Name: "name3", Email: "3@mail.com"},
					},
				},
			},
			wantCreate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "4", Name: "name4", Email: "4@mail.com"},
				},
			},
			wantUpdate: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "2", Name: "name2", Email: "22@mail.com"},
				},
			},
			wantEqual: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "1", Name: "name1", Email: "1@mail.com"},
				},
			},
			wantDelete: &model.GroupsResult{
				Items: 1,
				Resources: []*model.Group{
					{ID: "3", Name: "name3", Email: "3@mail.com"},
				},
			},
		},
		{
			name: "1 update, change the ID",
			args: args{
				idp: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{ID: "11", Name: "name1", Email: "1@mail.com"},
					},
				},
				state: &model.GroupsResult{
					Items: 1,
					Resources: []*model.Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
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
					{ID: "11", Name: "name1", Email: "1@mail.com"},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreate, gotUpdate, gotEqual, gotDelete := groupsDifferences(tt.args.idp, tt.args.state)
			if !reflect.DeepEqual(gotCreate, tt.wantCreate) {
				t.Errorf("groupsDifferences() gotCreate = %s, want %s", utils.ToJSON(gotCreate), utils.ToJSON(tt.wantCreate))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("groupsDifferences() gotUpdate = %s, want %s", utils.ToJSON(gotUpdate), utils.ToJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEqual, tt.wantEqual) {
				t.Errorf("groupsDifferences() gotEqual = %s, want %s", utils.ToJSON(gotEqual), utils.ToJSON(tt.wantEqual))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("groupsDifferences() gotDelete = %s, want %s", utils.ToJSON(gotDelete), utils.ToJSON(tt.wantDelete))
			}
		})
	}
}
