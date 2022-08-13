package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		sb := StateBuilder().Build()

		s := &State{}
		s.SetHashCode()

		assert.Equal(t, StateSchemaVersion, sb.SchemaVersion)
		assert.Equal(t, "", sb.CodeVersion)
		assert.Equal(t, "", sb.LastSync)
		assert.Equal(t, s.HashCode, sb.HashCode)
		assert.Equal(t, 0, sb.Resources.Groups.Items)
		assert.Equal(t, 0, len(sb.Resources.Groups.Resources))
		assert.Equal(t, 0, sb.Resources.Users.Items)
		assert.Equal(t, 0, len(sb.Resources.Users.Resources))
		assert.Equal(t, 0, sb.Resources.GroupsMembers.Items)
		assert.Equal(t, 0, len(sb.Resources.GroupsMembers.Resources))
	})

	t.Run("empty resources", func(t *testing.T) {
		sb := StateBuilder().
			WithSchemaVersion("1.0").
			WithCodeVersion("codeVersion").
			WithLastSync("lastSync").
			WithGroups(
				&GroupsResult{},
			).
			WithUsers(
				&UsersResult{},
			).
			WithGroupsMembers(
				&GroupsMembersResult{},
			).Build()

		s := &State{
			SchemaVersion: "1.0",
			CodeVersion:   "codeVersion",
			LastSync:      "lastSync",
			Resources: &StateResources{
				Groups:        &GroupsResult{},
				Users:         &UsersResult{},
				GroupsMembers: &GroupsMembersResult{},
			},
		}
		s.SetHashCode()

		assert.Equal(t, "1.0", sb.SchemaVersion)
		assert.Equal(t, "codeVersion", sb.CodeVersion)
		assert.Equal(t, "lastSync", sb.LastSync)
		assert.Equal(t, s.HashCode, sb.HashCode)
		assert.Equal(t, 0, sb.Resources.Groups.Items)
		assert.Equal(t, 0, len(sb.Resources.Groups.Resources))
		assert.Equal(t, 0, sb.Resources.Users.Items)
		assert.Equal(t, 0, len(sb.Resources.Users.Resources))
		assert.Equal(t, 0, sb.Resources.GroupsMembers.Items)
		assert.Equal(t, 0, len(sb.Resources.GroupsMembers.Resources))
	})

	t.Run("all resources", func(t *testing.T) {
		sb := StateBuilder().
			WithSchemaVersion("1.0").
			WithCodeVersion("codeVersion").
			WithLastSync("lastSync").
			WithGroups(
				&GroupsResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*Group{
						{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
					},
				},
			).
			WithUsers(
				&UsersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*User{
						{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, Email: "email"},
					},
				},
			).
			WithGroupsMembers(
				&GroupsMembersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*GroupMembers{
						{
							Items:    1,
							HashCode: "hashCode",
							Group:    &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
							Resources: []*Member{
								{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
							},
						},
					},
				},
			).Build()

		s := &State{
			SchemaVersion: "1.0",
			CodeVersion:   "codeVersion",
			LastSync:      "lastSync",
			Resources: &StateResources{
				Groups: &GroupsResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*Group{
						{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
					},
				},
				Users: &UsersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*User{
						{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, Email: "email"},
					},
				},
				GroupsMembers: &GroupsMembersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*GroupMembers{
						{
							Items:    1,
							HashCode: "hashCode",
							Group:    &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
							Resources: []*Member{
								{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
							},
						},
					},
				},
			},
		}
		s.SetHashCode()

		assert.Equal(t, "1.0", sb.SchemaVersion)
		assert.Equal(t, "codeVersion", sb.CodeVersion)
		assert.Equal(t, "lastSync", sb.LastSync)
		assert.Equal(t, s.HashCode, sb.HashCode)
		assert.Equal(t, 1, sb.Resources.Groups.Items)
		assert.Equal(t, 1, len(sb.Resources.Groups.Resources))
		assert.Equal(t, 1, sb.Resources.Users.Items)
		assert.Equal(t, 1, len(sb.Resources.Users.Resources))
		assert.Equal(t, 1, sb.Resources.GroupsMembers.Items)
		assert.Equal(t, 1, len(sb.Resources.GroupsMembers.Resources))
	})
}
