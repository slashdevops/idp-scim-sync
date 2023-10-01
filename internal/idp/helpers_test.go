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
				Addresses: []interface{}{
					map[string]interface{}{"formatted": "formatted work", "type": "work"},
					map[string]interface{}{"formatted": "formatted home", "type": "home"},
				},
				Emails:        []interface{}{map[string]interface{}{"address": "primaryEmail", "type": "work", "primary": true}},
				Languages:     []interface{}{map[string]interface{}{"languageCode": "languageCode", "preference": "preferred"}},
				Organizations: []interface{}{map[string]interface{}{"costCenter": "costCenter", "department": "department", "name": "name", "title": "title", "primary": true}},
				Phones: []interface{}{
					map[string]interface{}{"value": "value work", "type": "work"},
					map[string]interface{}{"value": "value home", "type": "home"},
				},
				Id:           "id",
				Kind:         "kind",
				PrimaryEmail: "primaryEmail",
				Suspended:    false,
				Name: &admin.UserName{
					GivenName:   "givenName",
					FamilyName:  "familyName",
					DisplayName: "displayName",
					FullName:    "fullName",
				},
				IsAdmin: false,
			},
			want: model.UserBuilder().
				WithName(&model.Name{GivenName: "givenName", FamilyName: "familyName", Formatted: "fullName"}).
				WithDisplayName("fullName").
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
				WithEnterpriseData(
					*model.EnterpriseDataBuilder().
						WithCostCenter("costCenter").
						WithDepartment("department").
						WithOrganization("name").
						WithPrimary(true).
						WithTitle("title").
						Build(),
				).
				WithPreferredLanguage("languageCode").
				WithTitle("title").
				WithTimezone("").
				WithPhoneNumber(
					model.PhoneNumberBuilder().
						WithValue("value work").
						WithType("work").
						Build(),
				).
				WithAddress(
					model.AddressBuilder().
						WithFormatted("formatted work").
						WithType("work").
						WithPrimary(true).
						Build(),
				).
				Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildUser(tt.given)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
