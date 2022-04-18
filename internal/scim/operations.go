package scim

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// patchGroupOperations assembles the operations for patch groups
// bases in the limits of operations we can execute in a single request.
func patchGroupOperations(op, path string, pvs []patchValue, gms *model.GroupMembers) []*aws.PatchGroupRequest {
	patchOperations := []*aws.PatchGroupRequest{}

	if len(pvs) > MaxPatchGroupMembersPerRequest {
		for i := 0; i < len(pvs); i += MaxPatchGroupMembersPerRequest {
			end := i + MaxPatchGroupMembersPerRequest
			if end > len(pvs) {
				end = len(pvs)
			}

			patchGroupRequest := &aws.PatchGroupRequest{
				Group: aws.Group{
					ID:          gms.Group.SCIMID,
					DisplayName: gms.Group.Name,
				},
				Patch: aws.Patch{
					Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
					Operations: []*aws.Operation{
						{
							OP:    op,
							Path:  path,
							Value: pvs[i:end],
						},
					},
				},
			}
			patchOperations = append(patchOperations, patchGroupRequest)
		}
	} else {
		patchGroupRequest := &aws.PatchGroupRequest{
			Group: aws.Group{
				ID:          gms.Group.SCIMID,
				DisplayName: gms.Group.Name,
			},
			Patch: aws.Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []*aws.Operation{
					{
						OP:    op,
						Path:  path,
						Value: pvs,
					},
				},
			},
		}
		patchOperations = append(patchOperations, patchGroupRequest)
	}

	return patchOperations
}
