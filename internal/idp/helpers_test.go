package idp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/stretchr/testify/assert"
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
				Addresses: []any{
					&admin.UserAddress{Formatted: "formatted work", Type: "work"},
					&admin.UserAddress{Formatted: "formatted home", Type: "home"},
				},
				Languages:     []any{&admin.UserLanguage{LanguageCode: "languageCode", Preference: "preferred"}},
				Organizations: []any{&admin.UserOrganization{CostCenter: "costCenter", Department: "department", Name: "name", Title: "title", Primary: true}},
				Phones: []any{
					&admin.UserPhone{Value: "value work", Type: "work"},
					&admin.UserPhone{Value: "value home", Type: "home"},
				},
				Id:           "id",
				Kind:         "kind",
				PrimaryEmail: "user@mail.com",
				Suspended:    false,
				Name: &admin.UserName{
					GivenName:  "givenName",
					FamilyName: "familyName",
					FullName:   "fullName",
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
				Addresses: []any{
					&admin.UserAddress{Formatted: "formatted work", Type: "work"},
					&admin.UserAddress{Formatted: "formatted home", Type: "home"},
				},
				Emails:        []any{&admin.UserEmail{Address: "user@mail.com", Type: "work", Primary: true}},
				Languages:     []any{&admin.UserLanguage{LanguageCode: "languageCode", Preference: "preferred"}},
				Organizations: []any{&admin.UserOrganization{CostCenter: "costCenter", Department: "department", Name: "name", Title: "title", Primary: true}},
				Phones: []any{
					&admin.UserPhone{Value: "value work", Type: "work"},
					&admin.UserPhone{Value: "value home", Type: "home"},
				},
				Id:           "id",
				Kind:         "kind",
				PrimaryEmail: "primaryEmail",
				Suspended:    false,
				Name: &admin.UserName{
					GivenName:  "givenName",
					FamilyName: "familyName",
					FullName:   "fullName",
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

func Test_toEmails(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    []model.Email
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return a valid email",
			given: []any{
				&admin.UserEmail{Address: "user@mail.com", Type: "work", Primary: true},
			},
			want: []model.Email{
				model.EmailBuilder().
					WithValue("user@mail.com").
					WithType("work").
					WithPrimary(true).
					Build(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toEmails(tt.given)
			if (err != nil) != tt.wantErr {
				t.Errorf("toEmails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toLanguages(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    string
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    "",
			wantErr: true,
		},
		{
			name: "should return a valid language",
			given: []any{
				&admin.UserLanguage{LanguageCode: "en-US", Preference: "preferred"},
			},
			want:    "en-US",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toLanguages(tt.given)
			if (err != nil) != tt.wantErr {
				t.Errorf("toLanguages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_toAddresses(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    []model.Address
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return a valid address",
			given: []any{
				&admin.UserAddress{Formatted: "123 Main St", Type: "work"},
			},
			want: []model.Address{
				model.AddressBuilder().
					WithFormatted("123 Main St").
					Build(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toAddresses(tt.given)
			if (err != nil) != tt.wantErr {
				t.Errorf("toAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toPhones(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    []model.PhoneNumber
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return a valid phone number",
			given: []any{
				&admin.UserPhone{Value: "555-555-5555", Type: "work"},
			},
			want: []model.PhoneNumber{
				model.PhoneNumberBuilder().
					WithValue("555-555-5555").
					WithType("work").
					Build(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toPhones(tt.given)
			if (err != nil) != tt.wantErr {
				t.Errorf("toPhones() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toRelations(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    *model.Manager
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return a valid manager",
			given: []any{
				&admin.UserRelation{Value: "manager@mail.com", Type: "manager"},
			},
			want: model.ManagerBuilder().
				WithValue("manager@mail.com").
				Build(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toRelations(tt.given)
			if (err != nil) != tt.wantErr {
				t.Errorf("toRelations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_toOrganizations(t *testing.T) {
	tests := []struct {
		name    string
		given   any
		want    *model.EnterpriseData
		want1   string
		wantErr bool
	}{
		{
			name:    "should return error when given is not a []any",
			given:   "not a slice",
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		{
			name: "should return a valid organization",
			given: []any{
				&admin.UserOrganization{Name: "ACME", Title: "Developer", Primary: true},
			},
			want: model.EnterpriseDataBuilder().
				WithOrganization("ACME").
				Build(),
			want1:   "Developer",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := toOrganizations(tt.given, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("toOrganizations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
