package idp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

func Test_buildUser(t *testing.T) {
	tests := []struct {
		name  string
		given *admin.User
		want  *model.User
	}{
		{
			name:  "should return nil when given is nil",
			given: nil,
			want:  nil,
		},
		{
			name:  "should return nil when given.Name is nil",
			given: &admin.User{Name: nil},
			want:  nil,
		},
		{
			name:  "should return nil when given.PrimaryEmail is empty",
			given: &admin.User{PrimaryEmail: ""},
			want:  nil,
		},
		{
			name:  "should return nil when given.Name.GivenName is empty",
			given: &admin.User{Name: &admin.UserName{GivenName: ""}},
			want:  nil,
		},
		{
			name:  "should return nil when given.Name.FamilyName is empty",
			given: &admin.User{Name: &admin.UserName{GivenName: "givenName", FamilyName: ""}},
			want:  nil,
		},
		{
			name: "should return a valid user with all the fields",
			given: &admin.User{
				Addresses:     []byte(`[{"country":"country","formatted":"formatted","locality":"locality","postalCode":"postalCode","primary":true,"region":"region","streetAddress":"streetAddress","type":"type"}]`),
				Emails:        []byte(`[{"address":"address","customType":"customType","primary":true,"type":"type"}]`),
				Languages:     []byte(`[{"customLanguage":"customLanguage","languageCode":"languageCode"}]`),
				Organizations: []byte(`[{"costCenter":"costCenter","customType":"customType","department":"department","description":"description","domain":"domain","fullTimeEquivalent":0,"location":"location","name":"name","primary":true,"symbol":"symbol","title":"title","type":"type"}]`),
				Phones:        []byte(`[{"customType":"customType","primary":true,"type":"type","value":"value"}]`),
				Id:            "id",
				Kind:          "kind",
				PrimaryEmail:  "primaryEmail",
				Suspended:     false,
				OrgUnitPath:   "orgUnitPath",
				Name: &admin.UserName{
					GivenName:   "givenName",
					FamilyName:  "familyName",
					DisplayName: "displayName",
					FullName:    "fullName",
				},
				IsAdmin: false,
			},
			want: model.UserBuilder().
				WithGivenName("givenName").
				WithFamilyName("familyName").
				WithDisplayName("displayName").
				WithEmail(
					model.EmailBuilder().
						WithValue("primaryEmail").
						WithType("work").
						WithPrimary(true).
						Build(),
				).
				WithActive(true).
				WithIPID("id").
				WithUserName("primaryEmail").
				WithUserType("kind").
				WithProfileURL("orgUnitPath").
				WithEnterpriseData(
					*model.EnterpriseDataBuilder().
						WithCostCenter("costCenter").
						WithDepartment("department").
						WithOrganization("name").
						Build(),
				).
				WithPreferredLanguage("languageCode").
				WithTitle("title").
				WithTimezone("").
				WithPhoneNumber(
					model.PhoneNumberBuilder().
						WithValue("value").
						WithType("type").
						Build(),
				).
				WithAddress(
					model.AddressBuilder().
						WithCountry("country").
						WithFormatted("formatted").
						WithLocality("locality").
						WithPostalCode("postalCode").
						WithPrimary(true).
						WithType("type").
						WithRegion("region").
						WithStreetAddress("streetAddress").
						Build(),
				).
				Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildUser(tt.given)

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("buildUser() got = %v, want %v", string(utils.ToJSON(got)), string(utils.ToJSON(tt.want)))
			// }

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
