package google

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/aws-sso-gws-sync/internal/sync"
	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
)

// toJSON return a json pretty of the stc
func toJSON(stc interface{}) []byte {
	JSON, err := json.MarshalIndent(stc, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return JSON
}

func TestNewGoogleIdentityProvider(t *testing.T) {

	t.Run("Should return IdentityServiceProvider and no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDS := NewMockDirectoryService(mockCtrl)

		svc, err := NewGoogleIdentityProvider(mockDS)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return an error if no DirectoryService is provided", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		svc, err := NewGoogleIdentityProvider(nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

func Test_googleProvider_GetGroups(t *testing.T) {
	type fields struct {
		ds *MockDirectoryService
	}

	type args struct {
		ctx    context.Context
		filter []string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *sync.GroupsResult
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
			want: &sync.GroupsResult{
				Items: 2,
				Resources: []*sync.Group{
					{Id: sync.Id{IdentityProvider: "1", SCIM: ""}, Name: "group1", Email: "group1@mail.com"},
					{Id: sync.Id{IdentityProvider: "2", SCIM: ""}, Name: "group2", Email: "group2@mail.com"},
				},
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

			f := fields{ds: NewMockDirectoryService(mockCtrl)}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			g := &googleProvider{
				ds: f.ds,
			}

			got, err := g.GetGroups(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleProvider.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleProvider.GetGroups() = %s, want %s", toJSON(got), toJSON(tt.want))
			}
		})
	}
}

func Test_googleProvider_GetUsers(t *testing.T) {
	type fields struct {
		ds DirectoryService
	}
	type args struct {
		ctx    context.Context
		filter []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sync.UsersResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &googleProvider{
				ds: tt.fields.ds,
			}
			got, err := g.GetUsers(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleProvider.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleProvider.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_googleProvider_GetGroupMembers(t *testing.T) {
	type fields struct {
		ds DirectoryService
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sync.MembersResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &googleProvider{
				ds: tt.fields.ds,
			}
			got, err := g.GetGroupMembers(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleProvider.GetGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleProvider.GetGroupMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_googleProvider_GetUsersFromGroupMembers(t *testing.T) {
	type fields struct {
		ds DirectoryService
	}
	type args struct {
		ctx context.Context
		mbr *sync.MembersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sync.UsersResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &googleProvider{
				ds: tt.fields.ds,
			}
			got, err := g.GetUsersFromGroupMembers(tt.args.ctx, tt.args.mbr)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleProvider.GetUsersFromGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleProvider.GetUsersFromGroupMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}
