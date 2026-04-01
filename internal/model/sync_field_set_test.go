package model

import "testing"

func TestNewSyncFieldSet(t *testing.T) {
	t.Run("nil creates empty set", func(t *testing.T) {
		s := NewSyncFieldSet(nil)
		if !s.IsEmpty() {
			t.Error("expected empty field set")
		}
	})

	t.Run("empty slice creates empty set", func(t *testing.T) {
		s := NewSyncFieldSet([]string{})
		if !s.IsEmpty() {
			t.Error("expected empty field set")
		}
	})

	t.Run("non-empty slice creates populated set", func(t *testing.T) {
		s := NewSyncFieldSet([]string{"phoneNumbers", "addresses"})
		if s.IsEmpty() {
			t.Error("expected non-empty field set")
		}
	})

	t.Run("slice with only empty strings creates empty set", func(t *testing.T) {
		s := NewSyncFieldSet([]string{""})
		if !s.IsEmpty() {
			t.Error("expected empty field set when input contains only empty strings")
		}
		// Should behave as 'all fields included'
		if !s.Includes(SyncUserFieldPhoneNumbers) {
			t.Error("empty-string-only set should include all fields")
		}
	})

	t.Run("empty strings are filtered from mixed input", func(t *testing.T) {
		s := NewSyncFieldSet([]string{"", "phoneNumbers", ""})
		if s.IsEmpty() {
			t.Error("expected non-empty field set")
		}
		if !s.Includes(SyncUserFieldPhoneNumbers) {
			t.Error("expected phoneNumbers to be included")
		}
		if s.Includes(SyncUserFieldAddresses) {
			t.Error("expected addresses to be excluded")
		}
	})
}

func TestSyncFieldSet_Includes(t *testing.T) {
	t.Run("nil field set includes everything", func(t *testing.T) {
		var s *SyncFieldSet
		if !s.Includes(SyncUserFieldPhoneNumbers) {
			t.Error("nil field set should include all fields")
		}
	})

	t.Run("empty field set includes everything", func(t *testing.T) {
		s := NewSyncFieldSet(nil)
		for _, field := range AllSyncUserFields {
			if !s.Includes(field) {
				t.Errorf("empty field set should include %s", field)
			}
		}
	})

	t.Run("configured set includes only specified fields", func(t *testing.T) {
		s := NewSyncFieldSet([]string{"phoneNumbers", "addresses"})

		if !s.Includes(SyncUserFieldPhoneNumbers) {
			t.Error("expected phoneNumbers to be included")
		}
		if !s.Includes(SyncUserFieldAddresses) {
			t.Error("expected addresses to be included")
		}
		if s.Includes(SyncUserFieldTitle) {
			t.Error("expected title to be excluded")
		}
		if s.Includes(SyncUserFieldEnterpriseData) {
			t.Error("expected enterpriseData to be excluded")
		}
	})
}

func TestValidateSyncUserField(t *testing.T) {
	t.Run("valid fields", func(t *testing.T) {
		for _, field := range AllSyncUserFields {
			if err := ValidateSyncUserField(string(field)); err != nil {
				t.Errorf("expected %s to be valid, got error: %v", field, err)
			}
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		if err := ValidateSyncUserField("invalidField"); err == nil {
			t.Error("expected error for invalid field")
		}
	})

	t.Run("empty string is invalid", func(t *testing.T) {
		if err := ValidateSyncUserField(""); err == nil {
			t.Error("expected error for empty field")
		}
	})
}
