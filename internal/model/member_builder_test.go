package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemberBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mb := MemberBuilder().Build()

		m := &Member{}
		m.SetHashCode()

		assert.Equal(t, "", mb.IPID)
		assert.Equal(t, "", mb.SCIMID)
		assert.Equal(t, "", mb.Email)
		assert.Equal(t, "", mb.Status)
		assert.Equal(t, m.HashCode, mb.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		mb := MemberBuilder()
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
		mb := MemberBuilder()
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

func TestMembersResultBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mrb := MembersResultBuilder().Build()

		mr := &MembersResult{
			Resources: make([]*Member, 0),
		}
		mr.SetHashCode()

		assert.Equal(t, 0, len(mrb.Resources))
		assert.Equal(t, mr.HashCode, mrb.HashCode)
	})

	t.Run("all options resources", func(t *testing.T) {
		mrb := MembersResultBuilder()
		mrb.WithResources([]*Member{
			{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
			{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
		}).Build()

		mr := &MembersResult{
			Items: 2,
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			},
		}
		mr.SetHashCode()

		assert.Equal(t, 2, len(mrb.mr.Resources))
		assert.Equal(t, "ipid", mrb.mr.Resources[0].IPID)
		assert.Equal(t, "ipid2", mrb.mr.Resources[1].IPID)
		assert.Equal(t, "scimid", mrb.mr.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", mrb.mr.Resources[1].SCIMID)
		assert.Equal(t, "email", mrb.mr.Resources[0].Email)
		assert.Equal(t, "email2", mrb.mr.Resources[1].Email)
		assert.Equal(t, "1", mrb.mr.Resources[0].Status)
		assert.Equal(t, "2", mrb.mr.Resources[1].Status)
		assert.Equal(t, mr.HashCode, mrb.mr.HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		mrb := MembersResultBuilder()
		mrb.WithResource(
			&Member{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
		).WithResource(
			&Member{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
		).WithResource(
			&Member{IPID: "ipid3", SCIMID: "scimid3", Email: "email3", Status: "3"},
		).Build()

		mr := &MembersResult{
			Items: 3,
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
				{IPID: "ipid3", SCIMID: "scimid3", Email: "email3", Status: "3"},
			},
		}
		mr.SetHashCode()

		assert.Equal(t, 3, len(mrb.mr.Resources))
		assert.Equal(t, "ipid", mrb.mr.Resources[0].IPID)
		assert.Equal(t, "ipid2", mrb.mr.Resources[1].IPID)
		assert.Equal(t, "ipid3", mrb.mr.Resources[2].IPID)
		assert.Equal(t, "scimid", mrb.mr.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", mrb.mr.Resources[1].SCIMID)
		assert.Equal(t, "scimid3", mrb.mr.Resources[2].SCIMID)
		assert.Equal(t, "email", mrb.mr.Resources[0].Email)
		assert.Equal(t, "email2", mrb.mr.Resources[1].Email)
		assert.Equal(t, "email3", mrb.mr.Resources[2].Email)
		assert.Equal(t, "1", mrb.mr.Resources[0].Status)
		assert.Equal(t, "2", mrb.mr.Resources[1].Status)
		assert.Equal(t, "3", mrb.mr.Resources[2].Status)
		assert.Equal(t, mr.HashCode, mrb.mr.HashCode)
	})
}

func TestGroupMembersBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gmb := GroupMembersBuilder().Build()

		gm := &GroupMembers{
			Resources: make([]*Member, 0),
		}
		gm.SetHashCode()

		assert.Equal(t, 0, len(gmb.Resources))
		assert.Equal(t, gm.HashCode, gmb.HashCode)
	})

	t.Run("with group and no resources", func(t *testing.T) {
		gmb := GroupMembersBuilder()
		gmb.WithGroup(
			&Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
		).Build()

		gm := &GroupMembers{
			Items:     0,
			Group:     &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			Resources: []*Member{},
		}
		gm.SetHashCode()

		assert.Equal(t, "group", gmb.gm.Group.Name)
		assert.Equal(t, "ipid", gmb.gm.Group.IPID)
		assert.Equal(t, "scimid", gmb.gm.Group.SCIMID)
		assert.Equal(t, "email", gmb.gm.Group.Email)

		assert.Equal(t, 0, len(gmb.gm.Resources))
		assert.Equal(t, gm.HashCode, gmb.gm.HashCode)
		assert.Equal(t, gm.Items, gmb.gm.Items)
	})

	t.Run("all options resources", func(t *testing.T) {
		gmb := GroupMembersBuilder()
		gmb.WithGroup(
			&Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
		).
			WithResources([]*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			}).Build()

		gm := &GroupMembers{
			Items: 2,
			Group: &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			},
		}
		gm.SetHashCode()

		assert.Equal(t, "group", gmb.gm.Group.Name)
		assert.Equal(t, "ipid", gmb.gm.Group.IPID)
		assert.Equal(t, "scimid", gmb.gm.Group.SCIMID)
		assert.Equal(t, "email", gmb.gm.Group.Email)

		assert.Equal(t, 2, len(gmb.gm.Resources))
		assert.Equal(t, "ipid", gmb.gm.Resources[0].IPID)
		assert.Equal(t, "ipid2", gmb.gm.Resources[1].IPID)
		assert.Equal(t, "scimid", gmb.gm.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", gmb.gm.Resources[1].SCIMID)
		assert.Equal(t, "email", gmb.gm.Resources[0].Email)
		assert.Equal(t, "email2", gmb.gm.Resources[1].Email)
		assert.Equal(t, "1", gmb.gm.Resources[0].Status)
		assert.Equal(t, "2", gmb.gm.Resources[1].Status)
		assert.Equal(t, gm.HashCode, gmb.gm.HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		gmb := GroupMembersBuilder()
		gmb.WithGroup(
			&Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
		).
			WithResource(
				&Member{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
			).WithResource(
			&Member{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
		).Build()

		gm := &GroupMembers{
			Items: 2,
			Group: &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			},
		}
		gm.SetHashCode()

		assert.Equal(t, 2, len(gmb.gm.Resources))
		assert.Equal(t, "ipid", gmb.gm.Resources[0].IPID)
		assert.Equal(t, "ipid2", gmb.gm.Resources[1].IPID)
		assert.Equal(t, "scimid", gmb.gm.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", gmb.gm.Resources[1].SCIMID)
		assert.Equal(t, "email", gmb.gm.Resources[0].Email)
		assert.Equal(t, "email2", gmb.gm.Resources[1].Email)
		assert.Equal(t, "1", gmb.gm.Resources[0].Status)
		assert.Equal(t, "2", gmb.gm.Resources[1].Status)
		assert.Equal(t, gm.HashCode, gmb.gm.HashCode)
	})
}

func TestGroupsMembersResultBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gmrb := GroupsMembersResultBuilder().Build()

		gmr := &GroupsMembersResult{
			Resources: make([]*GroupMembers, 0),
		}
		gmr.SetHashCode()

		assert.Equal(t, 0, len(gmrb.Resources))
		assert.Equal(t, gmr.HashCode, gmrb.HashCode)
	})

	t.Run("all options resources", func(t *testing.T) {
		gmrb := GroupsMembersResultBuilder()
		gmrb.WithResources(
			[]*GroupMembers{
				GroupMembersBuilder().WithGroup(
					&Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
				).WithResources(
					[]*Member{
						{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
						{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
					},
				).Build(),
			},
		).Build()

		gm := &GroupMembers{
			Items: 2,
			Group: &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			},
		}
		gm.SetHashCode()

		gmr := &GroupsMembersResult{
			Items: 1,
			Resources: []*GroupMembers{
				gm,
			},
		}
		gmr.SetHashCode()

		assert.Equal(t, "group", gmrb.gmr.Resources[0].Group.Name)
		assert.Equal(t, "ipid", gmrb.gmr.Resources[0].Group.IPID)
		assert.Equal(t, "scimid", gmrb.gmr.Resources[0].Group.SCIMID)
		assert.Equal(t, "email", gmrb.gmr.Resources[0].Group.Email)

		assert.Equal(t, 1, len(gmrb.gmr.Resources))
		assert.Equal(t, gmr.HashCode, gmrb.gmr.HashCode)
		assert.Equal(t, gm.HashCode, gmrb.gmr.Resources[0].HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		gmrb := GroupsMembersResultBuilder()
		gmrb.WithResource(
			GroupMembersBuilder().WithGroup(
				&Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			).WithResource(
				&Member{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
			).WithResource(
				&Member{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			).Build(),
		).Build()

		gm := &GroupMembers{
			Items: 2,
			Group: &Group{IPID: "ipid", SCIMID: "scimid", Name: "group", Email: "email"},
			Resources: []*Member{
				{IPID: "ipid", SCIMID: "scimid", Email: "email", Status: "1"},
				{IPID: "ipid2", SCIMID: "scimid2", Email: "email2", Status: "2"},
			},
		}
		gm.SetHashCode()

		gmr := &GroupsMembersResult{
			Items: 1,
			Resources: []*GroupMembers{
				gm,
			},
		}
		gmr.SetHashCode()

		assert.Equal(t, "group", gmrb.gmr.Resources[0].Group.Name)
		assert.Equal(t, "ipid", gmrb.gmr.Resources[0].Group.IPID)
		assert.Equal(t, "scimid", gmrb.gmr.Resources[0].Group.SCIMID)
		assert.Equal(t, "email", gmrb.gmr.Resources[0].Group.Email)

		assert.Equal(t, 1, len(gmrb.gmr.Resources))
		assert.Equal(t, gmr.HashCode, gmrb.gmr.HashCode)
		assert.Equal(t, gm.HashCode, gmrb.gmr.Resources[0].HashCode)
	})
}
