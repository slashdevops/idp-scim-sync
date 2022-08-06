package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ub := UserBuilder().Build()

		u := &User{}
		u.SetHashCode()

		assert.Equal(t, "", ub.IPID)
		assert.Equal(t, "", ub.SCIMID)
		assert.Equal(t, "", ub.Name.GivenName)
		assert.Equal(t, "", ub.Name.FamilyName)
		assert.Equal(t, "", ub.DisplayName)
		assert.Equal(t, false, ub.Active)
		assert.Equal(t, "", ub.Email)
		assert.Equal(t, u.HashCode, ub.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		ub := UserBuilder()
		ub.WithIPID("ipid").
			WithSCIMID("scimid").
			WithGivenName("givenname").
			WithFamilyName("familyname").
			WithDisplayName("displayname").
			WithActive(true).
			WithEmail("email").
			Build()

		u := &User{
			IPID:        "ipid",
			SCIMID:      "scimid",
			Name:        Name{GivenName: "givenname", FamilyName: "familyname"},
			DisplayName: "displayname",
			Active:      true,
			Email:       "email",
		}
		u.SetHashCode()

		assert.Equal(t, "ipid", ub.u.IPID)
		assert.Equal(t, "scimid", ub.u.SCIMID)
		assert.Equal(t, "givenname", ub.u.Name.GivenName)
		assert.Equal(t, "familyname", ub.u.Name.FamilyName)
		assert.Equal(t, "displayname", ub.u.DisplayName)
		assert.Equal(t, true, ub.u.Active)
		assert.Equal(t, "email", ub.u.Email)
		assert.Equal(t, u.HashCode, ub.u.HashCode)
	})

	t.Run("few options", func(t *testing.T) {
		ub := UserBuilder()
		ub.WithIPID("ipid").
			WithGivenName("givenname").
			WithActive(true).
			Build()

		u := &User{
			IPID:   "ipid",
			Name:   Name{GivenName: "givenname"},
			Active: true,
		}
		u.SetHashCode()

		assert.Equal(t, "ipid", ub.u.IPID)
		assert.Equal(t, "", ub.u.SCIMID)
		assert.Equal(t, "givenname", ub.u.Name.GivenName)
		assert.Equal(t, "", ub.u.Name.FamilyName)
		assert.Equal(t, "", ub.u.DisplayName)
		assert.Equal(t, true, ub.u.Active)
		assert.Equal(t, "", ub.u.Email)
		assert.Equal(t, u.HashCode, ub.u.HashCode)
	})
}

func TestUsersResultBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		urb := UsersResultBuilder().Build()

		ur := &UsersResult{
			Resources: make([]*User, 0),
		}
		ur.SetHashCode()

		assert.Equal(t, 0, len(urb.Resources))
		assert.Equal(t, ur.HashCode, urb.HashCode)
	})

	t.Run("all options resources", func(t *testing.T) {
		urb := UsersResultBuilder()
		urb.WithResources([]*User{
			{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Email: "email", Active: true},
			{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Email: "email2", Active: true},
		}).Build()

		ur := &UsersResult{
			Items: 2,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Email: "email", Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Email: "email2", Active: true},
			},
		}
		ur.SetHashCode()

		assert.Equal(t, 2, len(urb.ur.Resources))
		assert.Equal(t, "ipid", urb.ur.Resources[0].IPID)
		assert.Equal(t, "ipid2", urb.ur.Resources[1].IPID)
		assert.Equal(t, "scimid", urb.ur.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", urb.ur.Resources[1].SCIMID)
		assert.Equal(t, "1", urb.ur.Resources[0].Name.FamilyName)
		assert.Equal(t, "user", urb.ur.Resources[0].Name.GivenName)
		assert.Equal(t, "2", urb.ur.Resources[1].Name.FamilyName)
		assert.Equal(t, "user", urb.ur.Resources[1].Name.GivenName)
		assert.Equal(t, "email", urb.ur.Resources[0].Email)
		assert.Equal(t, "email2", urb.ur.Resources[1].Email)
		assert.Equal(t, true, urb.ur.Resources[0].Active)
		assert.Equal(t, true, urb.ur.Resources[1].Active)
		assert.Equal(t, ur.HashCode, urb.ur.HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		urb := UsersResultBuilder()
		urb.WithResource(
			&User{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Email: "email", Active: true},
		).WithResource(
			&User{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Email: "email2", Active: true},
		).WithResource(
			&User{IPID: "ipid3", SCIMID: "scimid3", Name: Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Email: "email3", Active: true},
		).Build()

		ur := &UsersResult{
			Items: 3,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Email: "email", Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Email: "email2", Active: true},
				{IPID: "ipid3", SCIMID: "scimid3", Name: Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Email: "email3", Active: true},
			},
		}
		ur.SetHashCode()

		assert.Equal(t, 3, len(urb.ur.Resources))
		assert.Equal(t, "ipid", urb.ur.Resources[0].IPID)
		assert.Equal(t, "ipid2", urb.ur.Resources[1].IPID)
		assert.Equal(t, "ipid3", urb.ur.Resources[2].IPID)
		assert.Equal(t, "scimid", urb.ur.Resources[0].SCIMID)
		assert.Equal(t, "scimid2", urb.ur.Resources[1].SCIMID)
		assert.Equal(t, "scimid3", urb.ur.Resources[2].SCIMID)
		assert.Equal(t, "1", urb.ur.Resources[0].Name.FamilyName)
		assert.Equal(t, "user", urb.ur.Resources[0].Name.GivenName)
		assert.Equal(t, "2", urb.ur.Resources[1].Name.FamilyName)
		assert.Equal(t, "user", urb.ur.Resources[1].Name.GivenName)
		assert.Equal(t, "3", urb.ur.Resources[2].Name.FamilyName)
		assert.Equal(t, "user", urb.ur.Resources[2].Name.GivenName)
		assert.Equal(t, "email", urb.ur.Resources[0].Email)
		assert.Equal(t, "email2", urb.ur.Resources[1].Email)
		assert.Equal(t, "email3", urb.ur.Resources[2].Email)
		assert.Equal(t, true, urb.ur.Resources[0].Active)
		assert.Equal(t, true, urb.ur.Resources[1].Active)
		assert.Equal(t, true, urb.ur.Resources[2].Active)
		assert.Equal(t, ur.HashCode, urb.ur.HashCode)
	})
}

func TestGroupBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		gb := GroupBuilder().Build()

		g := &Group{}
		g.SetHashCode()

		assert.Equal(t, "", gb.Email)
		assert.Equal(t, "", gb.SCIMID)
		assert.Equal(t, "", gb.Name)
		assert.Equal(t, "", gb.Email)
		assert.Equal(t, g.HashCode, gb.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		gb := GroupBuilder()
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
		gb := GroupBuilder()
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
		grb := GroupsResultBuilder().Build()

		gr := &GroupsResult{
			Resources: make([]*Group, 0),
		}
		gr.SetHashCode()

		assert.Equal(t, 0, len(grb.Resources))
		assert.Equal(t, gr.HashCode, grb.HashCode)
	})

	t.Run("all options resources", func(t *testing.T) {
		grb := GroupsResultBuilder()
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
		grb := GroupsResultBuilder()
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
