package scim

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	mock_scim "github.com/slashdevops/idp-scim-sync/mocks/scim"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"go.uber.org/mock/gomock"
)

// patchValueGenerator helper function to generate test data for patch operations
func patchValueGenerator(from, numValues int) []patchValue {
	values := make([]patchValue, numValues)

	for i := range numValues {
		values[i] = patchValue{
			Value: strconv.Itoa(i + from),
		}
	}
	return values
}

func TestProvider_patchGroupOperations(t *testing.T) {
	type fields struct {
		scim                 AWSSCIMProvider
		maxMembersPerRequest int
	}
	type args struct {
		op   string
		path string
		pvs  []patchValue
		gms  *model.GroupMembers
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*aws.PatchGroupRequest
	}{
		{
			name: "one member",
			fields: fields{
				maxMembersPerRequest: 100,
			},
			args: args{
				op:   "add",
				path: "members",
				pvs: []patchValue{
					{
						Value: "906722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36c",
					},
				},
				gms: &model.GroupMembers{
					Group: &model.Group{
						SCIMID: "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36d",
						Name:   "group 1",
					},
				},
			},
			want: []*aws.PatchGroupRequest{
				{
					Group: aws.Group{
						ID:          "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36d",
						DisplayName: "group 1",
					},
					Patch: aws.Patch{
						Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
						Operations: []*aws.Operation{
							{
								OP:   "add",
								Path: "members",
								Value: []patchValue{
									{
										Value: "906722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36c",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "more than 100",
			fields: fields{
				maxMembersPerRequest: 100,
			},
			args: args{
				op:   "add",
				path: "members",
				pvs:  patchValueGenerator(1, 120),
				gms: &model.GroupMembers{
					Group: &model.Group{
						SCIMID: "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36e",
						Name:   "group 1",
					},
				},
			},
			want: []*aws.PatchGroupRequest{
				{
					Group: aws.Group{
						ID:          "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36e",
						DisplayName: "group 1",
					},
					Patch: aws.Patch{
						Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
						Operations: []*aws.Operation{
							{
								OP:    "add",
								Path:  "members",
								Value: patchValueGenerator(1, 100),
							},
						},
					},
				},
				{
					Group: aws.Group{
						ID:          "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36e",
						DisplayName: "group 1",
					},
					Patch: aws.Patch{
						Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
						Operations: []*aws.Operation{
							{
								OP:    "add",
								Path:  "members",
								Value: patchValueGenerator(101, 20),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				scim:                 tt.fields.scim,
				maxMembersPerRequest: tt.fields.maxMembersPerRequest,
			}
			got := p.patchGroupOperations(tt.args.op, tt.args.path, tt.args.pvs, tt.args.gms)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("patchGroupOperations() (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkProvider_patchGroupOperations(b *testing.B) {
	p := &Provider{
		maxMembersPerRequest: 100,
	}
	for b.Loop() {
		p.patchGroupOperations("add", "members", patchValueGenerator(1, 350), &model.GroupMembers{
			Group: &model.Group{
				SCIMID: "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36e",
				Name:   "group 1",
			},
		})
	}
}

func TestNewProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type args struct {
		scim AWSSCIMProvider
		opts []ProviderOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Provider
		wantErr bool
	}{
		{
			name: "should return a new provider",
			args: args{
				scim: mockScimProvider,
			},
			want: &Provider{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			wantErr: false,
		},
		{
			name: "should return an error if scim is nil",
			args: args{
				scim: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return a new provider with options",
			args: args{
				scim: mockScimProvider,
				opts: []ProviderOption{
					WithMaxMembersPerRequest(200),
				},
			},
			want: &Provider{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvider(tt.args.scim, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got.maxMembersPerRequest != tt.want.maxMembersPerRequest {
					t.Errorf("NewProvider() = %v, want %v", got.maxMembersPerRequest, tt.want.maxMembersPerRequest)
				}
			}
		})
	}
}

func TestProvider_GetGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.GroupsResult
		wantErr bool
	}{
		{
			name: "should return groups",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(&aws.ListGroupsResponse{
					Resources: []*aws.Group{
						{
							ID:          "1",
							DisplayName: "group1",
							ExternalID:  "1",
						},
					},
				}, nil)
			},
			want: &model.GroupsResult{
				Resources: []*model.Group{
					{
						SCIMID: "1",
						Name:   "group1",
						IPID:   "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return empty groups",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(&aws.ListGroupsResponse{
					Resources: []*aws.Group{},
				}, nil)
			},
			want: &model.GroupsResult{
				Resources: []*model.Group{},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.GetGroups(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Group{}, "HashCode")); diff != "" {
				t.Errorf("Provider.GetGroups() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_CreateGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		gr  *model.GroupsResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.GroupsResult
		wantErr bool
	}{
		{
			name: "should create groups",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							Name: "group1",
							IPID: "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().CreateOrGetGroup(gomock.Any(), gomock.Any()).Return(&aws.CreateGroupResponse{
					ID: "1",
				}, nil)
			},
			want: &model.GroupsResult{
				Resources: []*model.Group{
					{
						SCIMID: "1",
						Name:   "group1",
						IPID:   "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							Name: "group1",
							IPID: "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().CreateOrGetGroup(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if groups result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.CreateGroups(tt.args.ctx, tt.args.gr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.CreateGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Group{}, "HashCode")); diff != "" {
				t.Errorf("Provider.CreateGroups() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_UpdateGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		gr  *model.GroupsResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.GroupsResult
		wantErr bool
	}{
		{
			name: "should update groups",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							SCIMID: "1",
							Name:   "group1",
							IPID:   "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &model.GroupsResult{
				Resources: []*model.Group{
					{
						SCIMID: "1",
						Name:   "group1",
						IPID:   "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							SCIMID: "1",
							Name:   "group1",
							IPID:   "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if groups result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.UpdateGroups(tt.args.ctx, tt.args.gr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.UpdateGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Group{}, "HashCode")); diff != "" {
				t.Errorf("Provider.UpdateGroups() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_DeleteGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		gr  *model.GroupsResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		wantErr bool
	}{
		{
			name: "should delete groups",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							SCIMID: "1",
							Name:   "group1",
							IPID:   "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{
							SCIMID: "1",
							Name:   "group1",
							IPID:   "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "should return an error if groups result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				gr:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			err := p.DeleteGroups(tt.args.ctx, tt.args.gr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.DeleteGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProvider_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.UsersResult
		wantErr bool
	}{
		{
			name: "should return users",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(&aws.ListUsersResponse{
					Resources: []*aws.User{
						{
							ID:         "1",
							UserName:   "user1",
							ExternalID: "1",
						},
					},
				}, nil)
			},
			want: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID:   "1",
						UserName: "user1",
						IPID:     "1",
						Name:     &model.Name{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.GetUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.UsersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.User{}, "HashCode")); diff != "" {
				t.Errorf("Provider.GetUsers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_CreateUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		ur  *model.UsersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.UsersResult
		wantErr bool
	}{
		{
			name: "should create users",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							UserName: "user1",
							IPID:     "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().CreateOrGetUser(gomock.Any(), gomock.Any()).Return(&aws.CreateUserResponse{
					ID: "1",
				}, nil)
			},
			want: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID:   "1",
						UserName: "user1",
						IPID:     "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							UserName: "user1",
							IPID:     "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().CreateOrGetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if users result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.CreateUsers(tt.args.ctx, tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.CreateUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.UsersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.User{}, "HashCode")); diff != "" {
				t.Errorf("Provider.CreateUsers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_UpdateUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		ur  *model.UsersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.UsersResult
		wantErr bool
	}{
		{
			name: "should update users",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID:   "1",
							UserName: "user1",
							IPID:     "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PutUser(gomock.Any(), gomock.Any()).Return(&aws.PutUserResponse{
					ID: "1",
				}, nil)
			},
			want: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID:   "1",
						UserName: "user1",
						IPID:     "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID:   "1",
							UserName: "user1",
							IPID:     "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PutUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if users result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if user scimid is empty",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							UserName: "user1",
							IPID:     "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.UpdateUsers(tt.args.ctx, tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.UpdateUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.UsersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.User{}, "HashCode")); diff != "" {
				t.Errorf("Provider.UpdateUsers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_DeleteUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim AWSSCIMProvider
	}
	type args struct {
		ctx context.Context
		ur  *model.UsersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		wantErr bool
	}{
		{
			name: "should delete users",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should return an error",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "should return an error if users result is nil",
			fields: fields{
				scim: mockScimProvider,
			},
			args: args{
				ctx: context.Background(),
				ur:  nil,
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			err := p.DeleteUsers(tt.args.ctx, tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.DeleteUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProvider_CreateGroupsMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim                 AWSSCIMProvider
		maxMembersPerRequest int
	}
	type args struct {
		ctx context.Context
		gmr *model.GroupsMembersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.GroupsMembersResult
		wantErr bool
	}{
		{
			name: "should create groups members",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									SCIMID: "1",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should create groups members without scim id",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									Email: "user1@email.com",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().GetUserByUserName(gomock.Any(), gomock.Any()).Return(&aws.GetUserResponse{
					ID: "1",
				}, nil)
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "1",
								Email:  "user1@email.com",
							},
						},
					},
				},
			},
		},
		{
			name: "should return an error when getting user by email",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									Email: "user1@email.com",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().GetUserByUserName(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "should return an error when patching group",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									SCIMID: "1",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim:                 tt.fields.scim,
				maxMembersPerRequest: tt.fields.maxMembersPerRequest,
			}
			got, err := p.CreateGroupsMembers(tt.args.ctx, tt.args.gmr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.CreateGroupsMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsMembersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.GroupMembers{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Member{}, "HashCode")); diff != "" {
				t.Errorf("Provider.CreateGroupsMembers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_DeleteGroupsMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim                 AWSSCIMProvider
		maxMembersPerRequest int
	}
	type args struct {
		ctx context.Context
		gmr *model.GroupsMembersResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		wantErr bool
	}{
		{
			name: "should delete groups members",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									SCIMID: "1",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should return an error when patching group",
			fields: fields{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			},
			args: args{
				ctx: context.Background(),
				gmr: &model.GroupsMembersResult{
					Resources: []*model.GroupMembers{
						{
							Group: &model.Group{
								SCIMID: "1",
								Name:   "group1",
							},
							Resources: []*model.Member{
								{
									SCIMID: "1",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().PatchGroup(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim:                 tt.fields.scim,
				maxMembersPerRequest: tt.fields.maxMembersPerRequest,
			}
			err := p.DeleteGroupsMembers(tt.args.ctx, tt.args.gmr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.DeleteGroupsMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// memberSortOpt provides a stable ordering for *model.Member slices in cmp diffs.
// The implementation builds membership lists from concurrent goroutines, so the
// observed order is non-deterministic.
var memberSortOpt = cmpopts.SortSlices(func(a, b *model.Member) bool {
	return a.SCIMID < b.SCIMID
})

// groupsMembersDiffOpts collects the cmp options used by GetGroupsMembers tests.
func groupsMembersDiffOpts() []cmp.Option {
	return []cmp.Option{
		cmpopts.IgnoreFields(model.GroupsMembersResult{}, "HashCode", "Items"),
		cmpopts.IgnoreFields(model.GroupMembers{}, "HashCode", "Items"),
		cmpopts.IgnoreFields(model.Member{}, "HashCode"),
		memberSortOpt,
	}
}

func TestProvider_GetGroupsMembers(t *testing.T) {
	type args struct {
		gr *model.GroupsResult
		ur *model.UsersResult
	}
	tests := []struct {
		name    string
		args    args
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		want    *model.GroupsMembersResult
		wantErr bool
	}{
		{
			name: "single user belongs to a single in-scope group",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{SCIMID: "g1", Name: "group1"},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Active: true,
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().
					ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
					Return(&aws.ListGroupsResponse{
						Resources: []*aws.Group{{ID: "g1", DisplayName: "group1"}},
					}, nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{SCIMID: "g1", Name: "group1"},
						Resources: []*model.Member{
							{SCIMID: "u1", Email: "user1@email.com", Status: "ACTIVE"},
						},
					},
				},
			},
		},
		{
			name: "user with no group memberships yields empty members list",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{SCIMID: "g1", Name: "group1"},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Active: true,
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().
					ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
					Return(&aws.ListGroupsResponse{
						Resources: []*aws.Group{},
					}, nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group:     &model.Group{SCIMID: "g1", Name: "group1"},
						Resources: []*model.Member{},
					},
				},
			},
		},
		{
			name: "memberships in out-of-scope groups are ignored",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{SCIMID: "g1", Name: "group1"},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Active: true,
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().
					ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
					Return(&aws.ListGroupsResponse{
						Resources: []*aws.Group{
							{ID: "g1", DisplayName: "group1"},
							{ID: "g-out-of-scope", DisplayName: "aws-managed-group"},
						},
					}, nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{SCIMID: "g1", Name: "group1"},
						Resources: []*model.Member{
							{SCIMID: "u1", Email: "user1@email.com", Status: "ACTIVE"},
						},
					},
				},
			},
		},
		{
			name: "inactive users do not get ACTIVE status",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{SCIMID: "g1", Name: "group1"},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Active: false,
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().
					ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
					Return(&aws.ListGroupsResponse{
						Resources: []*aws.Group{{ID: "g1", DisplayName: "group1"}},
					}, nil)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{SCIMID: "g1", Name: "group1"},
						Resources: []*model.Member{
							{SCIMID: "u1", Email: "user1@email.com"},
						},
					},
				},
			},
		},
		{
			name: "pagination walks every cursor and aggregates groups",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{
						{SCIMID: "g1", Name: "group1"},
						{SCIMID: "g2", Name: "group2"},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Active: true,
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				gomock.InOrder(
					m.EXPECT().
						ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
						Return(&aws.ListGroupsResponse{
							ListResponse: aws.ListResponse{NextCursor: "cursor-2"},
							Resources:    []*aws.Group{{ID: "g1", DisplayName: "group1"}},
						}, nil),
					m.EXPECT().
						ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "cursor-2").
						Return(&aws.ListGroupsResponse{
							Resources: []*aws.Group{{ID: "g2", DisplayName: "group2"}},
						}, nil),
				)
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{SCIMID: "g1", Name: "group1"},
						Resources: []*model.Member{
							{SCIMID: "u1", Email: "user1@email.com", Status: "ACTIVE"},
						},
					},
					{
						Group: &model.Group{SCIMID: "g2", Name: "group2"},
						Resources: []*model.Member{
							{SCIMID: "u1", Email: "user1@email.com", Status: "ACTIVE"},
						},
					},
				},
			},
		},
		{
			name: "error from AWS is propagated",
			args: args{
				gr: &model.GroupsResult{
					Resources: []*model.Group{{SCIMID: "g1", Name: "group1"}},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "u1",
							Emails: []model.Email{{Value: "user1@email.com"}},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().
					ListGroupsWithCursor(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("aws boom"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)
			tt.prepare(mockScimProvider)

			p := &Provider{
				scim:                 mockScimProvider,
				maxMembersPerRequest: 100,
			}

			got, err := p.GetGroupsMembers(context.Background(), tt.args.gr, tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetGroupsMembers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(tt.want, got, groupsMembersDiffOpts()...); diff != "" {
				t.Errorf("GetGroupsMembers() (-want +got):\n%s", diff)
			}
		})
	}
}

// TestProvider_GetGroupsMembers_ConcurrencyLimit exercises the errgroup
// concurrency cap using testing/synctest from the Go 1.26 standard library.
// synctest.Test runs every goroutine in an isolated "bubble" with a fake
// clock, so time.Sleep below advances virtual time only after the runtime has
// proven no other goroutine in the bubble can make progress. That gives us
// deterministic observation of the maximum number of in-flight calls without
// relying on wall-clock races.
//
// Reference: https://go.dev/blog/testing-time
func TestProvider_GetGroupsMembers_ConcurrencyLimit(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

		const userCount = 20
		groups := make([]*model.Group, 0, userCount)
		users := make([]*model.User, 0, userCount)
		for i := range userCount {
			id := strconv.Itoa(i)
			groups = append(groups, &model.Group{SCIMID: "g-" + id, Name: "group-" + id})
			users = append(users, &model.User{
				SCIMID: "u-" + id,
				Active: true,
				Emails: []model.Email{{Value: "user-" + id + "@email.com"}},
			})
		}

		var (
			mu             sync.Mutex
			inFlight       int
			maxInFlight    int
			observedSerial int
		)

		mockScimProvider.EXPECT().
			ListGroupsWithCursor(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, filter, _ string) (*aws.ListGroupsResponse, error) {
				mu.Lock()
				inFlight++
				if inFlight > maxInFlight {
					maxInFlight = inFlight
				}
				observedSerial++
				mu.Unlock()

				// Virtual sleep — synctest advances the fake clock only when
				// every other goroutine in the bubble is blocked, which lets us
				// observe the true concurrency cap rather than a flaky lower bound.
				time.Sleep(50 * time.Millisecond)

				mu.Lock()
				inFlight--
				mu.Unlock()

				// The filter encodes the user ID, so we can derive a deterministic
				// "user belongs to its matching group" result without juggling extra
				// state in the test fixture.
				userSCIMID := filter[len(`members.value eq "`) : len(filter)-1]
				groupSCIMID := "g-" + userSCIMID[len("u-"):]
				return &aws.ListGroupsResponse{
					Resources: []*aws.Group{{ID: groupSCIMID}},
				}, nil
			}).Times(userCount)

		p := &Provider{
			scim:                 mockScimProvider,
			maxMembersPerRequest: 100,
		}

		got, err := p.GetGroupsMembers(context.Background(), &model.GroupsResult{Resources: groups}, &model.UsersResult{Resources: users})
		if err != nil {
			t.Fatalf("GetGroupsMembers() error = %v", err)
		}

		if observedSerial != userCount {
			t.Errorf("expected exactly one ListGroupsWithCursor per user (got %d, want %d)", observedSerial, userCount)
		}
		if maxInFlight > getGroupsMembersConcurrency {
			t.Errorf("max in-flight calls = %d, exceeds cap %d", maxInFlight, getGroupsMembersConcurrency)
		}
		if maxInFlight < 2 {
			t.Errorf("expected concurrent calls (>=2), only observed %d in flight at peak", maxInFlight)
		}
		if len(got.Resources) != userCount {
			t.Errorf("expected %d group-members entries, got %d", userCount, len(got.Resources))
		}
	})
}
