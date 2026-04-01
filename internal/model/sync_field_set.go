package model

import "fmt"

// SyncUserField represents a configurable user field that can be included or excluded from sync.
type SyncUserField string

const (
	// SyncUserFieldPhoneNumbers controls syncing of phone number attributes.
	SyncUserFieldPhoneNumbers SyncUserField = "phoneNumbers"

	// SyncUserFieldAddresses controls syncing of address attributes.
	SyncUserFieldAddresses SyncUserField = "addresses"

	// SyncUserFieldTitle controls syncing of the user's job title.
	SyncUserFieldTitle SyncUserField = "title"

	// SyncUserFieldPreferredLanguage controls syncing of the user's preferred language.
	SyncUserFieldPreferredLanguage SyncUserField = "preferredLanguage"

	// SyncUserFieldLocale controls syncing of the user's locale.
	SyncUserFieldLocale SyncUserField = "locale"

	// SyncUserFieldTimezone controls syncing of the user's timezone.
	SyncUserFieldTimezone SyncUserField = "timezone"

	// SyncUserFieldNickName controls syncing of the user's nickname.
	SyncUserFieldNickName SyncUserField = "nickName"

	// SyncUserFieldProfileURL controls syncing of the user's profile URL.
	SyncUserFieldProfileURL SyncUserField = "profileURL"

	// SyncUserFieldUserType controls syncing of the user type attribute.
	SyncUserFieldUserType SyncUserField = "userType"

	// SyncUserFieldEnterpriseData controls syncing of enterprise extension attributes
	// (employeeNumber, costCenter, organization, division, department, manager).
	SyncUserFieldEnterpriseData SyncUserField = "enterpriseData"
)

// AllSyncUserFields contains all valid configurable sync user fields.
var AllSyncUserFields = []SyncUserField{
	SyncUserFieldPhoneNumbers,
	SyncUserFieldAddresses,
	SyncUserFieldTitle,
	SyncUserFieldPreferredLanguage,
	SyncUserFieldLocale,
	SyncUserFieldTimezone,
	SyncUserFieldNickName,
	SyncUserFieldProfileURL,
	SyncUserFieldUserType,
	SyncUserFieldEnterpriseData,
}

// validSyncUserFields is a lookup set of valid field names.
var validSyncUserFields = func() map[SyncUserField]struct{} {
	m := make(map[SyncUserField]struct{}, len(AllSyncUserFields))
	for _, f := range AllSyncUserFields {
		m[f] = struct{}{}
	}
	return m
}()

// ValidateSyncUserField checks whether a field name is a valid SyncUserField.
func ValidateSyncUserField(field string) error {
	if _, ok := validSyncUserFields[SyncUserField(field)]; !ok {
		return fmt.Errorf("unknown sync user field: %q", field)
	}
	return nil
}

// SyncFieldSet determines which optional user fields to include in the sync.
// When empty (no fields specified), all fields are included (backward compatible).
type SyncFieldSet struct {
	fields map[SyncUserField]bool
}

// NewSyncFieldSet creates a SyncFieldSet from a list of field names.
// An empty or nil slice means "include all fields".
func NewSyncFieldSet(fields []string) *SyncFieldSet {
	s := &SyncFieldSet{
		fields: make(map[SyncUserField]bool, len(fields)),
	}
	for _, f := range fields {
		s.fields[SyncUserField(f)] = true
	}
	return s
}

// IsEmpty returns true if no specific fields were configured,
// meaning all fields should be included.
func (s *SyncFieldSet) IsEmpty() bool {
	return len(s.fields) == 0
}

// Includes returns true if the given field should be included in the sync.
// When no fields are configured (empty set), all fields are included.
func (s *SyncFieldSet) Includes(field SyncUserField) bool {
	if s == nil || s.IsEmpty() {
		return true
	}
	return s.fields[field]
}
