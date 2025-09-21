package scim

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
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
			name: "should return an error when member has no scim id",
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
				// No expectations since we expect it to fail before calling any methods
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

func TestProvider_GetGroupsMembers(t *testing.T) {
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
		want    *model.GroupsMembersResult
		wantErr bool
	}{
		{
			name: "should get groups members",
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
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(&aws.ListGroupsResponse{
					Resources: []*aws.Group{
						{
							Members: []*aws.Member{
								{
									Value: "1",
								},
							},
						},
					},
				}, nil)
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&aws.GetUserResponse{
					Emails: []aws.Email{
						{
							Value: "user1@email.com",
						},
					},
				}, nil)
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
			wantErr: false,
		},
		{
			name: "should return an error when listing groups",
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
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "should return an error when getting user",
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
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(&aws.ListGroupsResponse{
					Resources: []*aws.Group{
						{
							Members: []*aws.Member{
								{
									Value: "1",
								},
							},
						},
					},
				}, nil)
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.GetGroupsMembers(tt.args.ctx, tt.args.gr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetGroupsMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsMembersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.GroupMembers{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Member{}, "HashCode")); diff != "" {
				t.Errorf("Provider.GetGroupsMembers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_GetGroupsMembersBruteForce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)

	type fields struct {
		scim                 AWSSCIMProvider
		maxMembersPerRequest int
	}
	type args struct {
		ctx context.Context
		gr  *model.GroupsResult
		ur  *model.UsersResult
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
			name: "should get groups members",
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
						},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "1",
							Active: true,
							Emails: []model.Email{
								{
									Value: "user1@email.com",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(&aws.ListGroupsResponse{
					ListResponse: aws.ListResponse{
						TotalResults: 1,
					},
					Resources: []*aws.Group{},
				}, nil)
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
								Status: "ACTIVE",
							},
						},
					},
				},
			},
		},
		{
			name: "should get groups members with concurrency",
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
						},
						{
							SCIMID: "2",
							Name:   "group2",
						},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "1",
							Active: true,
							Emails: []model.Email{
								{
									Value: "user1@email.com",
								},
							},
						},
						{
							SCIMID: "2",
							Active: true,
							Emails: []model.Email{
								{
									Value: "user2@email.com",
								},
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				var wg sync.WaitGroup
				var mu sync.Mutex
				var currentConcurrent int
				var maxConcurrent int

				calls := 4
				wg.Add(calls)

				for i := 0; i < calls; i++ {
					m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ string) (*aws.ListGroupsResponse, error) {
						mu.Lock()
						currentConcurrent++
						if currentConcurrent > maxConcurrent {
							maxConcurrent = currentConcurrent
						}
						mu.Unlock()

						time.Sleep(100 * time.Millisecond)

						mu.Lock()
						currentConcurrent--
						mu.Unlock()

						wg.Done()

						return &aws.ListGroupsResponse{
							ListResponse: aws.ListResponse{
								TotalResults: 1,
							},
							Resources: []*aws.Group{},
						}, nil
					})
				}

				go func() {
					wg.Wait()
					if maxConcurrent > 20 {
						t.Errorf("max concurrent calls should be less than 20, got %d", maxConcurrent)
					}
				}()
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
								Status: "ACTIVE",
							},
							{
								SCIMID: "2",
								Email:  "user2@email.com",
								Status: "ACTIVE",
							},
						},
					},
					{
						Group: &model.Group{
							SCIMID: "2",
							Name:   "group2",
						},
						Resources: []*model.Member{
							{
								SCIMID: "1",
								Email:  "user1@email.com",
								Status: "ACTIVE",
							},
							{
								SCIMID: "2",
								Email:  "user2@email.com",
								Status: "ACTIVE",
							},
						},
					},
				},
			},
		},
		{
			name: "should return an error when listing groups",
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
						},
					},
				},
				ur: &model.UsersResult{
					Resources: []*model.User{
						{
							SCIMID: "1",
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			p := &Provider{
				scim: tt.fields.scim,
			}
			got, err := p.GetGroupsMembersBruteForce(tt.args.ctx, tt.args.gr, tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetGroupsMembersBruteForce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// sort members by SCIMID to avoid order issues in the test
			opt := cmpopts.SortSlices(func(a, b *model.Member) bool {
				return a.SCIMID < b.SCIMID
			})

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.GroupsMembersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.GroupMembers{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Member{}, "HashCode"), opt); diff != "" {
				t.Errorf("Provider.GetGroupsMembersBruteForce() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProvider_PopulateSCIMIDsForGroupMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScimProvider := mock_scim.NewMockAWSSCIMProvider(ctrl)
	p := &Provider{scim: mockScimProvider}

	tests := []struct {
		name    string
		gmr     *model.GroupsMembersResult
		users   *model.UsersResult
		want    *model.GroupsMembersResult
		prepare func(m *mock_scim.MockAWSSCIMProvider)
		wantErr bool
	}{
		{
			name: "should populate SCIM IDs for group members",
			gmr: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								Email: "user1@example.com",
							},
							{
								Email: "user2@example.com",
							},
						},
					},
				},
			},
			users: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID: "user-1",
						Emails: []model.Email{
							{Value: "user1@example.com", Primary: true},
						},
					},
					{
						SCIMID: "user-2",
						Emails: []model.Email{
							{Value: "user2@example.com", Primary: true},
						},
					},
				},
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "user-1",
								Email:  "user1@example.com",
							},
							{
								SCIMID: "user-2",
								Email:  "user2@example.com",
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				// No SCIM calls expected since all users found in reconciled list
			},
		},
		{
			name: "should handle members with existing SCIM IDs",
			gmr: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "existing-id",
								Email:  "user1@example.com",
							},
							{
								Email: "user2@example.com",
							},
						},
					},
				},
			},
			users: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID: "user-1",
						Emails: []model.Email{
							{Value: "user1@example.com", Primary: true},
						},
					},
					{
						SCIMID: "user-2",
						Emails: []model.Email{
							{Value: "user2@example.com", Primary: true},
						},
					},
				},
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "existing-id", // Should keep existing SCIM ID
								Email:  "user1@example.com",
							},
							{
								SCIMID: "user-2",
								Email:  "user2@example.com",
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				// No SCIM calls expected since all users found in reconciled list
			},
		},
		{
			name: "should handle missing users with SCIM lookup fallback",
			gmr: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								Email: "user1@example.com",
							},
							{
								Email: "missing@example.com", // This user is not in the users list
							},
						},
					},
				},
			},
			users: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID: "user-1",
						Emails: []model.Email{
							{Value: "user1@example.com", Primary: true},
						},
					},
				},
			},
			want: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								SCIMID: "user-1",
								Email:  "user1@example.com",
							},
							{
								SCIMID: "scim-fallback-id", // Should get ID from SCIM lookup
								Email:  "missing@example.com",
							},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().GetUserByUserName(gomock.Any(), "missing@example.com").Return(&aws.GetUserResponse{
					ID: "scim-fallback-id",
				}, nil)
			},
		},
		{
			name: "should return error when user not found anywhere",
			gmr: &model.GroupsMembersResult{
				Resources: []*model.GroupMembers{
					{
						Group: &model.Group{
							SCIMID: "group-1",
							Name:   "group1",
						},
						Resources: []*model.Member{
							{
								Email: "user1@example.com",
							},
							{
								Email: "notfound@example.com", // This user is not in the users list and not in SCIM
							},
						},
					},
				},
			},
			users: &model.UsersResult{
				Resources: []*model.User{
					{
						SCIMID: "user-1",
						Emails: []model.Email{
							{Value: "user1@example.com", Primary: true},
						},
					},
				},
			},
			prepare: func(m *mock_scim.MockAWSSCIMProvider) {
				m.EXPECT().GetUserByUserName(gomock.Any(), "notfound@example.com").Return(&aws.GetUserResponse{
					ID: "", // Empty ID means user not found
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(mockScimProvider)
			err := p.PopulateSCIMIDsForGroupMembers(context.Background(), tt.gmr, tt.users)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.PopulateSCIMIDsForGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, tt.gmr, cmpopts.IgnoreFields(model.GroupsMembersResult{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.GroupMembers{}, "HashCode", "Items"), cmpopts.IgnoreFields(model.Member{}, "HashCode")); diff != "" {
					t.Errorf("Provider.PopulateSCIMIDsForGroupMembers() (-want +got):\n%s", diff)
				}
			}
		})
	}
}
