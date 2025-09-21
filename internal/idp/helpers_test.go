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
			name: "should return a valid user with no addresses, no languages, no organizations, no phones",
			given: &admin.User{
				Id:           "id",
				Kind:         "kind",
				PrimaryEmail: "user@mail.com",
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
						WithValue("user@mail.com").
						WithType("work").
						WithPrimary(true).
						Build(),
				).
				WithActive(true).
				WithIPID("id").
				WithUserName("user@mail.com").
				WithUserType("kind").
				Build(),
		},
		{
			name: "should return a valid user with all the fields and using primary email as the only email",
			given: &admin.User{
				Addresses: []interface{}{
					map[string]interface{}{"formatted": "formatted work", "type": "work"},
					map[string]interface{}{"formatted": "formatted home", "type": "home"},
				},
				Languages:     []interface{}{map[string]interface{}{"languageCode": "languageCode", "preference": "preferred"}},
				Organizations: []interface{}{map[string]interface{}{"costCenter": "costCenter", "department": "department", "name": "name", "title": "title", "primary": true}},
				Phones: []interface{}{
					map[string]interface{}{"value": "value work", "type": "work"},
					map[string]interface{}{"value": "value home", "type": "home"},
				},
				Id:           "id",
				Kind:         "kind",
				PrimaryEmail: "user@mail.com",
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
						WithValue("user@mail.com").
						WithType("work").
						WithPrimary(true).
						Build(),
				).
				WithActive(true).
				WithIPID("id").
				WithUserName("user@mail.com").
				WithUserType("kind").
				WithTitle("title").
				WithEnterpriseData(
					model.EnterpriseDataBuilder().
						WithCostCenter("costCenter").
						WithDepartment("department").
						WithOrganization("name").
						Build(),
				).
				WithPreferredLanguage("languageCode").
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
						Build(),
				).
				Build(),
		},
		{
			name: "should return a valid user with all the fields",
			given: &admin.User{
				Addresses: []interface{}{
					map[string]interface{}{"formatted": "formatted work", "type": "work"},
					map[string]interface{}{"formatted": "formatted home", "type": "home"},
				},
				Emails:        []interface{}{map[string]interface{}{"address": "user@mail.com", "type": "work", "primary": true}},
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
						WithValue("user@mail.com").
						WithType("work").
						WithPrimary(true).
						Build(),
				).
				WithActive(true).
				WithIPID("id").
				WithUserName("primaryEmail").
				WithUserType("kind").
				WithTitle("title").
				WithEnterpriseData(
					model.EnterpriseDataBuilder().
						WithCostCenter("costCenter").
						WithDepartment("department").
						WithOrganization("name").
						Build(),
				).
				WithPreferredLanguage("languageCode").
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

func Test_buildEmailQuery(t *testing.T) {
	tests := []struct {
		name   string
		emails []string
		want   string
	}{
		{
			name:   "should return empty string for empty slice",
			emails: []string{},
			want:   "",
		},
		{
			name:   "should return single email query for one email",
			emails: []string{"user@example.com"},
			want:   "email:user@example.com",
		},
		{
			name:   "should return OR query for multiple emails",
			emails: []string{"user1@example.com", "user2@example.com"},
			want:   "email:user1@example.com OR email:user2@example.com",
		},
		{
			name:   "should handle three emails correctly",
			emails: []string{"user1@example.com", "user2@example.com", "user3@example.com"},
			want:   "email:user1@example.com OR email:user2@example.com OR email:user3@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildEmailQuery(tt.emails)
			if got != tt.want {
				t.Errorf("buildEmailQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_chunkEmails(t *testing.T) {
	tests := []struct {
		name      string
		emails    []string
		chunkSize int
		want      [][]string
	}{
		{
			name:      "should return single chunk for empty slice",
			emails:    []string{},
			chunkSize: 2,
			want:      [][]string{},
		},
		{
			name:      "should return single chunk when size is larger than emails",
			emails:    []string{"email1", "email2"},
			chunkSize: 5,
			want:      [][]string{{"email1", "email2"}},
		},
		{
			name:      "should return multiple chunks when emails exceed chunk size",
			emails:    []string{"email1", "email2", "email3", "email4", "email5"},
			chunkSize: 2,
			want:      [][]string{{"email1", "email2"}, {"email3", "email4"}, {"email5"}},
		},
		{
			name:      "should handle exact chunk size division",
			emails:    []string{"email1", "email2", "email3", "email4"},
			chunkSize: 2,
			want:      [][]string{{"email1", "email2"}, {"email3", "email4"}},
		},
		{
			name:      "should return original slice for zero chunk size",
			emails:    []string{"email1", "email2", "email3"},
			chunkSize: 0,
			want:      [][]string{{"email1", "email2", "email3"}},
		},
		{
			name:      "should return original slice for negative chunk size",
			emails:    []string{"email1", "email2", "email3"},
			chunkSize: -1,
			want:      [][]string{{"email1", "email2", "email3"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chunkEmails(tt.emails, tt.chunkSize)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("chunkEmails() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
