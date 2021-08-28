package core

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func Test_groupsDifferences(t *testing.T) {
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
		},
		{
			name: "equals",
			args: args{
				idp: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
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
					{ID: "1", Name: "name1", Email: "1@mail.com"},
					{ID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
			},
		},
		{
			name: "1 equals, 1 update, 1 create, 1 delete",
			args: args{
				idp: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{ID: "1", Name: "name1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
					},
				},
				state: &GroupsResult{
					Items: 2,
					Resources: []*Group{
						{ID: "1", Name: "newname1", Email: "1@mail.com"},
						{ID: "2", Name: "name2", Email: "2@mail.com"},
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
					{ID: "1", Name: "name1", Email: "1@mail.com"},
					{ID: "2", Name: "name2", Email: "2@mail.com"},
				},
			},
			wantDelete: &GroupsResult{
				Items:     0,
				Resources: []*Group{},
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
