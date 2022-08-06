package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gb := NewGroupBuilder().Build()

		g := &Group{}
		g.SetHashCode()

		assert.Equal(t, "", gb.Email)
		assert.Equal(t, "", gb.SCIMID)
		assert.Equal(t, "", gb.Name)
		assert.Equal(t, "", gb.Email)
		assert.Equal(t, g.HashCode, gb.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		gb := NewGroupBuilder()
		gb.WithIPID("ipid").WithSCIMID("scimid").WithName("name").WithEmail("email").Build()

		g := &Group{
			IPID:   "ipid",
			SCIMID: "scimid",
			Name:   "name",
			Email:  "email",
		}
		g.SetHashCode()

		assert.Equal(t, "ipid", gb.g.IPID)
		assert.Equal(t, "scimid", gb.g.SCIMID)
		assert.Equal(t, "name", gb.g.Name)
		assert.Equal(t, "email", gb.g.Email)
		assert.Equal(t, g.HashCode, gb.g.HashCode)
	})

	t.Run("few options", func(t *testing.T) {
		gb := NewGroupBuilder()
		gb.WithName("name").WithEmail("email").Build()

		g := &Group{
			Name:  "name",
			Email: "email",
		}
		g.SetHashCode()

		assert.Equal(t, "", gb.g.IPID)
		assert.Equal(t, "", gb.g.SCIMID)
		assert.Equal(t, "name", gb.g.Name)
		assert.Equal(t, "email", gb.g.Email)
		assert.Equal(t, g.HashCode, gb.g.HashCode)
	})
}

func TestGroupsResultBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		grb := NewGroupsResultBuilder().Build()

		gr := &GroupsResult{
			Resources: make([]*Group, 0),
		}
		gr.SetHashCode()

		assert.Equal(t, 0, len(grb.Resources))
		assert.Equal(t, gr.HashCode, grb.HashCode)
	})

	t.Run("all options resources", func(t *testing.T) {
		grb := NewGroupsResultBuilder()
		grb.WithResources([]*Group{
			{IPID: "ipid", SCIMID: "scimid", Name: "name", Email: "email"},
			{IPID: "ipid2", SCIMID: "scimid2", Name: "name2", Email: "email2"},
		}).Build()

		gr := &GroupsResult{
			Items: 2,
			Resources: []*Group{
				{IPID: "ipid", SCIMID: "scimid", Name: "name", Email: "email"},
				{IPID: "ipid2", SCIMID: "scimid2", Name: "name2", Email: "email2"},
			},
		}
		gr.SetHashCode()

		assert.Equal(t, 2, len(grb.gr.Resources))
		assert.Equal(t, "ipid", grb.gr.Resources[0].IPID)
		assert.Equal(t, "ipid2", grb.gr.Resources[1].IPID)
		assert.Equal(t, "scimid", grb.gr.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", grb.gr.Resources[1].SCIMID)
		assert.Equal(t, "name", grb.gr.Resources[0].Name)
		assert.Equal(t, "name2", grb.gr.Resources[1].Name)
		assert.Equal(t, "email", grb.gr.Resources[0].Email)
		assert.Equal(t, "email2", grb.gr.Resources[1].Email)
		assert.Equal(t, gr.HashCode, grb.gr.HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		grb := NewGroupsResultBuilder()
		grb.WithResource(
			&Group{IPID: "ipid", SCIMID: "scimid", Name: "name", Email: "email"},
		).WithResource(
			&Group{IPID: "ipid2", SCIMID: "scimid2", Name: "name2", Email: "email2"},
		).WithResource(
			&Group{IPID: "ipid3", SCIMID: "scimid3", Name: "name3", Email: "email3"},
		).Build()

		gr := &GroupsResult{
			Items: 3,
			Resources: []*Group{
				{IPID: "ipid", SCIMID: "scimid", Name: "name", Email: "email"},
				{IPID: "ipid2", SCIMID: "scimid2", Name: "name2", Email: "email2"},
				{IPID: "ipid3", SCIMID: "scimid3", Name: "name3", Email: "email3"},
			},
		}
		gr.SetHashCode()

		assert.Equal(t, 3, len(grb.gr.Resources))
		assert.Equal(t, "ipid", grb.gr.Resources[0].IPID)
		assert.Equal(t, "ipid2", grb.gr.Resources[1].IPID)
		assert.Equal(t, "ipid3", grb.gr.Resources[2].IPID)
		assert.Equal(t, "scimid", grb.gr.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", grb.gr.Resources[1].SCIMID)
		assert.Equal(t, "scimid3", grb.gr.Resources[2].SCIMID)
		assert.Equal(t, "name", grb.gr.Resources[0].Name)
		assert.Equal(t, "name2", grb.gr.Resources[1].Name)
		assert.Equal(t, "name3", grb.gr.Resources[2].Name)
		assert.Equal(t, "email", grb.gr.Resources[0].Email)
		assert.Equal(t, "email2", grb.gr.Resources[1].Email)
		assert.Equal(t, "email3", grb.gr.Resources[2].Email)
		assert.Equal(t, gr.HashCode, grb.gr.HashCode)
	})
}

func TestMemberBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mb := NewMemberBuilder().Build()

		m := &Member{}
		m.SetHashCode()

		assert.Equal(t, "", mb.IPID)
		assert.Equal(t, "", mb.SCIMID)
		assert.Equal(t, "", mb.Email)
		assert.Equal(t, "", mb.Status)
		assert.Equal(t, m.HashCode, mb.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		mb := NewMemberBuilder()
		mb.WithIPID("ipid").WithSCIMID("scimid").WithEmail("email").WithStatus("status").Build()

		m := &Member{
			IPID:   "ipid",
			SCIMID: "scimid",
			Email:  "email",
			Status: "status",
		}
		m.SetHashCode()

		assert.Equal(t, "ipid", mb.m.IPID)
		assert.Equal(t, "scimid", mb.m.SCIMID)
		assert.Equal(t, "email", mb.m.Email)
		assert.Equal(t, "status", mb.m.Status)
		assert.Equal(t, m.HashCode, mb.m.HashCode)
	})

	t.Run("few options", func(t *testing.T) {
		mb := NewMemberBuilder()
		mb.WithIPID("ipid").WithStatus("status").Build()

		m := &Member{
			IPID:   "ipid",
			Status: "status",
		}
		m.SetHashCode()

		assert.Equal(t, "ipid", mb.m.IPID)
		assert.Equal(t, "", mb.m.SCIMID)
		assert.Equal(t, "", mb.m.Email)
		assert.Equal(t, "status", mb.m.Status)
		assert.Equal(t, m.HashCode, mb.m.HashCode)
	})
}
