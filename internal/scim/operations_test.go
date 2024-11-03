package scim

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// patchValueGenerator helper function to generate test data for patch operations
func patchValueGenerator(from, numValues int) []patchValue {
	values := make([]patchValue, numValues)

	for i := 0; i < numValues; i++ {
		values[i] = patchValue{
			Value: strconv.Itoa(i + from),
		}
	}
	return values
}

func Test_patchGroupOperations(t *testing.T) {
	type args struct {
		op   string
		path string
		pvs  []patchValue
		gms  *model.GroupMembers
	}
	tests := []struct {
		name string
		args args
		want []*aws.PatchGroupRequest
	}{
		{
			name: "one member",
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
			got := patchGroupOperations(tt.args.op, tt.args.path, tt.args.pvs, tt.args.gms)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("patchGroupOperations() (-want +got):\n%s", diff)
			}
		})
	}
}

func Benchmark_patchGroupOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		patchGroupOperations("add", "members", patchValueGenerator(1, 350), &model.GroupMembers{
			Group: &model.Group{
				SCIMID: "016722b2be-ee23ed58-6e4e-4b2f-a94a-3ace8456a36e",
				Name:   "group 1",
			},
		})
	}
}
