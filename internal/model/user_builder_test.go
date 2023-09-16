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
		assert.Equal(t, nil, ub.Emails)
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
			WithEmail(Email{Value: "email"}).
			Build()

		u := &User{
			IPID:        "ipid",
			SCIMID:      "scimid",
			Name:        Name{GivenName: "givenname", FamilyName: "familyname"},
			DisplayName: "displayname",
			Active:      true,
			Emails:      []Email{{Value: "email"}},
		}
		u.SetHashCode()

		assert.Equal(t, "ipid", ub.u.IPID)
		assert.Equal(t, "scimid", ub.u.SCIMID)
		assert.Equal(t, "givenname", ub.u.Name.GivenName)
		assert.Equal(t, "familyname", ub.u.Name.FamilyName)
		assert.Equal(t, "displayname", ub.u.DisplayName)
		assert.Equal(t, true, ub.u.Active)
		assert.Equal(t, "email", ub.u.Emails[0].Value)
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
		assert.Equal(t, nil, ub.u.Emails)
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
			{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
			{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
		}).Build()

		ur := &UsersResult{
			Items: 2,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
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
		assert.Equal(t, "email", urb.ur.Resources[0].Emails[0].Value)
		assert.Equal(t, "email2", urb.ur.Resources[1].Emails[0].Value)
		assert.Equal(t, true, urb.ur.Resources[0].Active)
		assert.Equal(t, true, urb.ur.Resources[1].Active)
		assert.Equal(t, ur.HashCode, urb.ur.HashCode)
	})

	t.Run("all options resource", func(t *testing.T) {
		urb := UsersResultBuilder()
		urb.WithResource(
			&User{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
		).WithResource(
			&User{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
		).WithResource(
			&User{IPID: "ipid3", SCIMID: "scimid3", Name: Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Emails: []Email{{Value: "email3"}}, Active: true},
		).Build()

		ur := &UsersResult{
			Items: 3,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
				{IPID: "ipid3", SCIMID: "scimid3", Name: Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Emails: []Email{{Value: "email3"}}, Active: true},
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
		assert.Equal(t, "email", urb.ur.Resources[0].Emails[0].Value)
		assert.Equal(t, "email2", urb.ur.Resources[1].Emails[0].Value)
		assert.Equal(t, "email3", urb.ur.Resources[2].Emails[0].Value)
		assert.Equal(t, true, urb.ur.Resources[0].Active)
		assert.Equal(t, true, urb.ur.Resources[1].Active)
		assert.Equal(t, true, urb.ur.Resources[2].Active)
		assert.Equal(t, ur.HashCode, urb.ur.HashCode)
	})
}
