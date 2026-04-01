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

		assert.Equal(t, u.IPID, ub.IPID)
		assert.Equal(t, u.SCIMID, ub.SCIMID)
		assert.Equal(t, u.Name, ub.Name)
		assert.Equal(t, u.Active, ub.Active)
		assert.Equal(t, u.Emails, ub.Emails)
		assert.Equal(t, u.HashCode, ub.HashCode)
	})

	t.Run("all options", func(t *testing.T) {
		ub := UserBuilder()
		ub.WithIPID("ipid").
			WithSCIMID("scimid").
			WithName(&Name{GivenName: "givenname", FamilyName: "familyname"}).
			WithDisplayName("displayname").
			WithActive(true).
			WithEmail(Email{Value: "email"}).
			Build()

		u := &User{
			IPID:        "ipid",
			SCIMID:      "scimid",
			Name:        &Name{GivenName: "givenname", FamilyName: "familyname"},
			DisplayName: "displayname",
			Active:      true,
			Emails:      []Email{{Value: "email"}},
		}
		u.SetHashCode()

		assert.Equal(t, u.IPID, ub.u.IPID)
		assert.Equal(t, u.SCIMID, ub.u.SCIMID)
		assert.Equal(t, u.Name, ub.u.Name)
		assert.Equal(t, u.Active, ub.u.Active)
		assert.Equal(t, u.Emails, ub.u.Emails)
		assert.Equal(t, u.HashCode, ub.u.HashCode)
	})

	t.Run("few options", func(t *testing.T) {
		ub := UserBuilder()
		ub.WithIPID("ipid").
			WithName(&Name{GivenName: "givenname"}).
			WithActive(true).
			Build()

		u := &User{
			IPID:   "ipid",
			Name:   &Name{GivenName: "givenname"},
			Active: true,
		}
		u.SetHashCode()

		assert.Equal(t, u.IPID, ub.u.IPID)
		assert.Equal(t, u.SCIMID, ub.u.SCIMID)
		assert.Equal(t, u.Name, ub.u.Name)
		assert.Equal(t, u.Active, ub.u.Active)
		assert.Equal(t, u.Emails, ub.u.Emails)
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
			{IPID: "ipid", SCIMID: "scimid", Name: &Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
			{IPID: "ipid2", SCIMID: "scimid2", Name: &Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
		}).Build()

		ur := &UsersResult{
			Items: 2,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: &Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: &Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
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

	t.Run("single resource", func(t *testing.T) {
		urb := UsersResultBuilder()
		urb.WithResource(
			&User{IPID: "ipid", SCIMID: "scimid", Name: &Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
		).WithResource(
			&User{IPID: "ipid2", SCIMID: "scimid2", Name: &Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
		).WithResource(
			&User{IPID: "ipid3", SCIMID: "scimid3", Name: &Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Emails: []Email{{Value: "email3"}}, Active: true},
		).Build()

		ur := &UsersResult{
			Items: 3,
			Resources: []*User{
				{IPID: "ipid", SCIMID: "scimid", Name: &Name{FamilyName: "1", GivenName: "user"}, DisplayName: "user 1", Emails: []Email{{Value: "email"}}, Active: true},
				{IPID: "ipid2", SCIMID: "scimid2", Name: &Name{FamilyName: "2", GivenName: "user"}, DisplayName: "user 2", Emails: []Email{{Value: "email2"}}, Active: true},
				{IPID: "ipid3", SCIMID: "scimid3", Name: &Name{FamilyName: "3", GivenName: "user"}, DisplayName: "user 3", Emails: []Email{{Value: "email3"}}, Active: true},
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

func TestNameBuilder(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		n := NameBuilder().
			WithFormatted("John M Doe Jr").
			WithFamilyName("Doe").
			WithGivenName("John").
			WithMiddleName("M").
			WithHonorificPrefix("Mr").
			WithHonorificSuffix("Jr").
			Build()

		assert.Equal(t, "John M Doe Jr", n.Formatted)
		assert.Equal(t, "Doe", n.FamilyName)
		assert.Equal(t, "John", n.GivenName)
		assert.Equal(t, "M", n.MiddleName)
		assert.Equal(t, "Mr", n.HonorificPrefix)
		assert.Equal(t, "Jr", n.HonorificSuffix)
	})

	t.Run("empty", func(t *testing.T) {
		n := NameBuilder().Build()
		assert.NotNil(t, n)
		assert.Equal(t, "", n.GivenName)
		assert.Equal(t, "", n.FamilyName)
	})
}

func TestEnterpriseDataBuilder(t *testing.T) {
	t.Run("all fields with manager", func(t *testing.T) {
		mgr := ManagerBuilder().
			WithValue("boss-id").
			WithRef("Users/boss-id").
			Build()

		ed := EnterpriseDataBuilder().
			WithEmployeeNumber("E001").
			WithCostCenter("CC100").
			WithOrganization("Acme Corp").
			WithDivision("Technology").
			WithDepartment("Engineering").
			WithManager(mgr).
			Build()

		assert.Equal(t, "E001", ed.EmployeeNumber)
		assert.Equal(t, "CC100", ed.CostCenter)
		assert.Equal(t, "Acme Corp", ed.Organization)
		assert.Equal(t, "Technology", ed.Division)
		assert.Equal(t, "Engineering", ed.Department)
		assert.NotNil(t, ed.Manager)
		assert.Equal(t, "boss-id", ed.Manager.Value)
		assert.Equal(t, "Users/boss-id", ed.Manager.Ref)
	})

	t.Run("empty", func(t *testing.T) {
		ed := EnterpriseDataBuilder().Build()
		assert.NotNil(t, ed)
		assert.Equal(t, "", ed.EmployeeNumber)
		assert.Nil(t, ed.Manager)
	})
}

func TestManagerBuilder(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		m := ManagerBuilder().
			WithValue("mgr-id").
			WithRef("Users/mgr-id").
			Build()

		assert.Equal(t, "mgr-id", m.Value)
		assert.Equal(t, "Users/mgr-id", m.Ref)
	})

	t.Run("empty", func(t *testing.T) {
		m := ManagerBuilder().Build()
		assert.NotNil(t, m)
		assert.Equal(t, "", m.Value)
	})
}

func TestAddressBuilder(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		a := AddressBuilder().
			WithFormatted("123 Main St, Springfield, IL 62701, US").
			WithStreetAddress("123 Main St").
			WithLocality("Springfield").
			WithRegion("IL").
			WithPostalCode("62701").
			WithCountry("US").
			Build()

		assert.Equal(t, "123 Main St, Springfield, IL 62701, US", a.Formatted)
		assert.Equal(t, "123 Main St", a.StreetAddress)
		assert.Equal(t, "Springfield", a.Locality)
		assert.Equal(t, "IL", a.Region)
		assert.Equal(t, "62701", a.PostalCode)
		assert.Equal(t, "US", a.Country)
	})

	t.Run("empty", func(t *testing.T) {
		a := AddressBuilder().Build()
		assert.Equal(t, "", a.Formatted)
		assert.Equal(t, "", a.Country)
	})
}

func TestPhoneNumberBuilder(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		pn := PhoneNumberBuilder().
			WithValue("+1-555-1234").
			WithType("work").
			Build()

		assert.Equal(t, "+1-555-1234", pn.Value)
		assert.Equal(t, "work", pn.Type)
	})

	t.Run("empty", func(t *testing.T) {
		pn := PhoneNumberBuilder().Build()
		assert.Equal(t, "", pn.Value)
		assert.Equal(t, "", pn.Type)
	})
}

func TestUserBuilder_WithAddress(t *testing.T) {
	t.Run("add address to empty user", func(t *testing.T) {
		u := UserBuilder().
			WithIPID("ipid").
			WithAddress(AddressBuilder().WithFormatted("addr1").Build()).
			Build()

		assert.Len(t, u.Addresses, 1)
		assert.Equal(t, "addr1", u.Addresses[0].Formatted)
	})

	t.Run("replace existing address", func(t *testing.T) {
		u := UserBuilder().
			WithIPID("ipid").
			WithAddress(AddressBuilder().WithFormatted("first").Build()).
			WithAddress(AddressBuilder().WithFormatted("replaced").Build()).
			Build()

		assert.Len(t, u.Addresses, 1)
		assert.Equal(t, "replaced", u.Addresses[0].Formatted)
	})
}

func TestUserBuilder_WithPhoneNumber(t *testing.T) {
	t.Run("add phone to empty user", func(t *testing.T) {
		u := UserBuilder().
			WithIPID("ipid").
			WithPhoneNumber(PhoneNumberBuilder().WithValue("+1234").WithType("work").Build()).
			Build()

		assert.Len(t, u.PhoneNumbers, 1)
		assert.Equal(t, "+1234", u.PhoneNumbers[0].Value)
	})

	t.Run("replace existing phone", func(t *testing.T) {
		u := UserBuilder().
			WithIPID("ipid").
			WithPhoneNumber(PhoneNumberBuilder().WithValue("first").WithType("work").Build()).
			WithPhoneNumber(PhoneNumberBuilder().WithValue("replaced").WithType("home").Build()).
			Build()

		assert.Len(t, u.PhoneNumbers, 1)
		assert.Equal(t, "replaced", u.PhoneNumbers[0].Value)
		assert.Equal(t, "home", u.PhoneNumbers[0].Type)
	})
}

func TestUserBuilder_AllOptionalFields(t *testing.T) {
	t.Run("nickname, profileURL, locale, timezone", func(t *testing.T) {
		u := UserBuilder().
			WithIPID("ipid").
			WithNickName("johnny").
			WithProfileURL("https://example.com/john").
			WithLocale("en_US").
			WithTimezone("America/Chicago").
			WithUserType("Employee").
			WithTitle("Engineer").
			WithPreferredLanguage("en").
			Build()

		assert.Equal(t, "johnny", u.NickName)
		assert.Equal(t, "https://example.com/john", u.ProfileURL)
		assert.Equal(t, "en_US", u.Locale)
		assert.Equal(t, "America/Chicago", u.Timezone)
		assert.Equal(t, "Employee", u.UserType)
		assert.Equal(t, "Engineer", u.Title)
		assert.Equal(t, "en", u.PreferredLanguage)
	})
}
