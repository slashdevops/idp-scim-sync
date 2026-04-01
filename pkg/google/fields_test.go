package google

import (
	"strings"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

func Test_buildUserFields(t *testing.T) {
	t.Run("nil field set returns all fields including locations", func(t *testing.T) {
		result := buildUserFields(nil)

		requiredFields := []string{"id", "etag", "primaryEmail", "name", "suspended", "kind", "emails"}
		for _, f := range requiredFields {
			if !strings.Contains(result, f) {
				t.Errorf("expected field %q in result %q", f, result)
			}
		}

		optionalFields := []string{"addresses", "phones", "languages", "organizations", "relations", "locations"}
		for _, f := range optionalFields {
			if !strings.Contains(result, f) {
				t.Errorf("expected field %q in result %q", f, result)
			}
		}
	})

	t.Run("empty field set returns all fields", func(t *testing.T) {
		fields := model.NewSyncFieldSet(nil)
		result := buildUserFields(fields)

		if !strings.Contains(result, "addresses") {
			t.Errorf("expected 'addresses' in result %q", result)
		}
		if !strings.Contains(result, "phones") {
			t.Errorf("expected 'phones' in result %q", result)
		}
		if !strings.Contains(result, "locations") {
			t.Errorf("expected 'locations' in result %q", result)
		}
	})

	t.Run("only phoneNumbers includes phones but not addresses", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"phoneNumbers"})
		result := buildUserFields(fields)

		if !strings.Contains(result, "phones") {
			t.Errorf("expected 'phones' in result %q", result)
		}
		if strings.Contains(result, "addresses") {
			t.Errorf("unexpected 'addresses' in result %q", result)
		}
		if strings.Contains(result, "organizations") {
			t.Errorf("unexpected 'organizations' in result %q", result)
		}
		if strings.Contains(result, "locations") {
			t.Errorf("unexpected 'locations' in result %q", result)
		}
	})

	t.Run("enterpriseData includes organizations and relations", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"enterpriseData"})
		result := buildUserFields(fields)

		if !strings.Contains(result, "organizations") {
			t.Errorf("expected 'organizations' in result %q", result)
		}
		if !strings.Contains(result, "relations") {
			t.Errorf("expected 'relations' in result %q", result)
		}
	})

	t.Run("title includes organizations but not relations", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"title"})
		result := buildUserFields(fields)

		if !strings.Contains(result, "organizations") {
			t.Errorf("expected 'organizations' in result %q", result)
		}
		if strings.Contains(result, "relations") {
			t.Errorf("unexpected 'relations' in result %q", result)
		}
	})

	t.Run("required fields always present", func(t *testing.T) {
		fields := model.NewSyncFieldSet([]string{"phoneNumbers"})
		result := buildUserFields(fields)

		for _, required := range []string{"primaryEmail", "name", "suspended", "kind", "emails"} {
			if !strings.Contains(result, required) {
				t.Errorf("expected required field %q in result %q", required, result)
			}
		}
	})
}
