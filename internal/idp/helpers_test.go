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
				Addresses: []any{
					map[string]any{"formatted": "formatted work", "type": "work"},
					map[string]any{"formatted": "formatted home", "type": "home"},
				},
				Languages:     []any{map[string]any{"languageCode": "languageCode", "preference": "preferred"}},
				Organizations: []any{map[string]any{"costCenter": "costCenter", "department": "department", "name": "name", "title": "title", "primary": true}},
				Phones: []any{
					map[string]any{"value": "value work", "type": "work"},
					map[string]any{"value": "value home", "type": "home"},
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
				Addresses: []any{
					map[string]any{"formatted": "formatted work", "type": "work"},
					map[string]any{"formatted": "formatted home", "type": "home"},
				},
				Emails:        []any{map[string]any{"address": "user@mail.com", "type": "work", "primary": true}},
				Languages:     []any{map[string]any{"languageCode": "languageCode", "preference": "preferred"}},
				Organizations: []any{map[string]any{"costCenter": "costCenter", "department": "department", "name": "name", "title": "title", "primary": true}},
				Phones: []any{
					map[string]any{"value": "value work", "type": "work"},
					map[string]any{"value": "value home", "type": "home"},
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
			// nil field set = all fields (backward compatible)
			got := buildUser(tt.given, nil)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_buildUser_withFieldSet(t *testing.T) {
	fullUser := &admin.User{
		Addresses: []any{
			map[string]any{"formatted": "formatted work", "type": "work"},
		},
		Emails:        []any{map[string]any{"address": "user@mail.com", "type": "work", "primary": true}},
		Languages:     []any{map[string]any{"languageCode": "en", "preference": "preferred"}},
		Organizations: []any{map[string]any{"costCenter": "CC1", "department": "Eng", "name": "Acme", "title": "Engineer", "primary": true}},
		Phones: []any{
			map[string]any{"value": "+1234", "type": "work"},
		},
		Id:           "id1",
		Kind:         "admin#directory#user",
		PrimaryEmail: "user@mail.com",
		Suspended:    false,
		Name: &admin.UserName{
			GivenName:  "John",
			FamilyName: "Doe",
			FullName:   "John Doe",
		},
	}

	t.Run("empty field set syncs all fields", func(t *testing.T) {
		fields := model.NewSyncFieldSet(nil)
		got := buildUser(fullUser, fields)

		if got.Title != "Engineer" {
			t.Errorf("expected title 'Engineer', got %q", got.Title)
		}
		if len(got.PhoneNumbers) == 0 {
			t.Error("expected phone numbers to be populated")
		}
		if len(got.Addresses) == 0 {
			t.Error("expected addresses to be populated")
		}
		if got.EnterpriseData == nil {
			t.Error("expected enterprise data to be populated")
		}
		if got.PreferredLanguage != "en" {
			t.Errorf("expected preferred language 'en', got %q", got.PreferredLanguage)
		}
	})

	t.Run("only phoneNumbers field includes only phones", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"phoneNumbers"})
		got := buildUser(fullUser, fields)

		// Required fields are always present
		if got.UserName != "user@mail.com" {
			t.Errorf("expected userName 'user@mail.com', got %q", got.UserName)
		}
		if got.DisplayName != "John Doe" {
			t.Errorf("expected displayName 'John Doe', got %q", got.DisplayName)
		}

		// phoneNumbers should be included
		if len(got.PhoneNumbers) == 0 {
			t.Error("expected phone numbers to be populated")
		}

		// Everything else should be excluded
		if len(got.Addresses) != 0 {
			t.Error("expected addresses to be empty")
		}
		if got.Title != "" {
			t.Errorf("expected empty title, got %q", got.Title)
		}
		if got.EnterpriseData != nil {
			t.Error("expected enterprise data to be nil")
		}
		if got.PreferredLanguage != "" {
			t.Errorf("expected empty preferred language, got %q", got.PreferredLanguage)
		}
		if got.UserType != "" {
			t.Errorf("expected empty user type, got %q", got.UserType)
		}
	})

	t.Run("title and enterpriseData include organizations", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"title", "enterpriseData"})
		got := buildUser(fullUser, fields)

		if got.Title != "Engineer" {
			t.Errorf("expected title 'Engineer', got %q", got.Title)
		}
		if got.EnterpriseData == nil {
			t.Fatal("expected enterprise data to be populated")
		}
		if got.EnterpriseData.Department != "Eng" {
			t.Errorf("expected department 'Eng', got %q", got.EnterpriseData.Department)
		}

		// Excluded fields
		if len(got.PhoneNumbers) != 0 {
			t.Error("expected phone numbers to be empty")
		}
		if len(got.Addresses) != 0 {
			t.Error("expected addresses to be empty")
		}
	})

	t.Run("addresses only", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"addresses"})
		got := buildUser(fullUser, fields)

		if len(got.Addresses) == 0 {
			t.Error("expected addresses to be populated")
		}
		if got.Addresses[0].Formatted != "formatted work" {
			t.Errorf("expected formatted 'formatted work', got %q", got.Addresses[0].Formatted)
		}

		// Excluded
		if len(got.PhoneNumbers) != 0 {
			t.Error("expected phone numbers to be empty")
		}
		if got.Title != "" {
			t.Errorf("expected empty title, got %q", got.Title)
		}
	})
}

func Test_toRelations(t *testing.T) {
	t.Run("should return nil with nil input", func(t *testing.T) {
		got, err := toRelations(nil)
		if err == nil {
			t.Error("expected error for nil input")
		}
		if got != nil {
			t.Error("expected nil manager")
		}
	})

	t.Run("should return nil for invalid type", func(t *testing.T) {
		got, err := toRelations("invalid")
		if err == nil {
			t.Error("expected error for invalid type")
		}
		if got != nil {
			t.Error("expected nil manager")
		}
	})

	t.Run("should return manager from map with manager type", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "manager", "value": "boss@example.com", "customType": "direct"},
		}
		got, err := toRelations(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected manager, got nil")
		}
		if got.Value != "boss@example.com" {
			t.Errorf("expected value 'boss@example.com', got %q", got.Value)
		}
		if got.Ref != "direct" {
			t.Errorf("expected ref 'direct', got %q", got.Ref)
		}
	})

	t.Run("should skip non-manager relation types", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "assistant", "value": "asst@example.com"},
			map[string]any{"type": "manager", "value": "mgr@example.com"},
		}
		got, err := toRelations(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected manager, got nil")
		}
		if got.Value != "mgr@example.com" {
			t.Errorf("expected value 'mgr@example.com', got %q", got.Value)
		}
	})

	t.Run("should return nil when no manager relation exists", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "assistant", "value": "asst@example.com"},
		}
		got, err := toRelations(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != nil {
			t.Error("expected nil manager for non-manager relations")
		}
	})

	t.Run("should handle invalid element type", func(t *testing.T) {
		input := []any{42}
		got, err := toRelations(input)
		if err == nil {
			t.Error("expected error for invalid element type")
		}
		if got != nil {
			t.Error("expected nil manager")
		}
	})
}

func Test_toLanguages(t *testing.T) {
	t.Run("should return error for invalid type", func(t *testing.T) {
		got, err := toLanguages("invalid")
		if err == nil {
			t.Error("expected error")
		}
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("should return empty for no preferred language", func(t *testing.T) {
		input := []any{
			map[string]any{"languageCode": "en", "preference": "not-preferred"},
		}
		got, err := toLanguages(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("should handle invalid element type", func(t *testing.T) {
		input := []any{42}
		got, err := toLanguages(input)
		if err == nil {
			t.Error("expected error for invalid element type")
		}
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func Test_toAddresses(t *testing.T) {
	t.Run("should return error for invalid type", func(t *testing.T) {
		got, err := toAddresses("invalid")
		if err == nil {
			t.Error("expected error")
		}
		if got != nil {
			t.Error("expected nil")
		}
	})

	t.Run("should skip non-work/home types", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "other", "formatted": "other addr"},
		}
		got, err := toAddresses(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Error("expected empty addresses for non-work/home type")
		}
	})

	t.Run("should return home address", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "home", "formatted": "123 Home St"},
		}
		got, err := toAddresses(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatal("expected 1 address")
		}
		if got[0].Formatted != "123 Home St" {
			t.Errorf("expected '123 Home St', got %q", got[0].Formatted)
		}
	})

	t.Run("should handle invalid element type", func(t *testing.T) {
		input := []any{42}
		got, err := toAddresses(input)
		if err == nil {
			t.Error("expected error for invalid element type")
		}
		if got != nil {
			t.Error("expected nil")
		}
	})
}

func Test_toPhones(t *testing.T) {
	t.Run("should return error for invalid type", func(t *testing.T) {
		got, err := toPhones("invalid")
		if err == nil {
			t.Error("expected error")
		}
		if got != nil {
			t.Error("expected nil")
		}
	})

	t.Run("should skip non-work/home types", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "mobile", "value": "+1234"},
		}
		got, err := toPhones(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Error("expected empty phones for non-work/home type")
		}
	})

	t.Run("should return home phone", func(t *testing.T) {
		input := []any{
			map[string]any{"type": "home", "value": "+5678"},
		}
		got, err := toPhones(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatal("expected 1 phone")
		}
		if got[0].Value != "+5678" {
			t.Errorf("expected '+5678', got %q", got[0].Value)
		}
		if got[0].Type != "home" {
			t.Errorf("expected 'home', got %q", got[0].Type)
		}
	})

	t.Run("should handle invalid element type", func(t *testing.T) {
		input := []any{42}
		got, err := toPhones(input)
		if err == nil {
			t.Error("expected error for invalid element type")
		}
		if got != nil {
			t.Error("expected nil")
		}
	})
}

func Test_toOrganizations(t *testing.T) {
	t.Run("should return error for invalid type", func(t *testing.T) {
		org, title, err := toOrganizations("invalid", nil)
		if err == nil {
			t.Error("expected error")
		}
		if org != nil || title != "" {
			t.Error("expected nil org and empty title")
		}
	})

	t.Run("should skip non-primary organizations", func(t *testing.T) {
		input := []any{
			map[string]any{"primary": false, "name": "Acme", "title": "Dev"},
		}
		org, title, err := toOrganizations(input, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if org != nil {
			t.Error("expected nil org for non-primary")
		}
		if title != "" {
			t.Errorf("expected empty title, got %q", title)
		}
	})

	t.Run("should return primary organization with all fields", func(t *testing.T) {
		input := []any{
			map[string]any{
				"primary":    true,
				"name":       "Acme Corp",
				"title":      "Engineer",
				"costCenter": "CC100",
				"department": "Engineering",
				"domain":     "acme.com",
			},
		}
		mgr := &model.Manager{Value: "boss@acme.com"}
		org, title, err := toOrganizations(input, mgr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if org == nil {
			t.Fatal("expected org, got nil")
		}
		if title != "Engineer" {
			t.Errorf("expected 'Engineer', got %q", title)
		}
		if org.Organization != "Acme Corp" {
			t.Errorf("expected 'Acme Corp', got %q", org.Organization)
		}
		if org.CostCenter != "CC100" {
			t.Errorf("expected 'CC100', got %q", org.CostCenter)
		}
		if org.Department != "Engineering" {
			t.Errorf("expected 'Engineering', got %q", org.Department)
		}
		if org.Division != "acme.com" {
			t.Errorf("expected 'acme.com', got %q", org.Division)
		}
		if org.Manager == nil || org.Manager.Value != "boss@acme.com" {
			t.Error("expected manager to be set")
		}
	})

	t.Run("should handle invalid element type", func(t *testing.T) {
		input := []any{42}
		org, title, err := toOrganizations(input, nil)
		if err == nil {
			t.Error("expected error for invalid element type")
		}
		if org != nil || title != "" {
			t.Error("expected nil org and empty title")
		}
	})
}
