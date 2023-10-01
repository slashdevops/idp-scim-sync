package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
)

// Name represents a name entity and is used in other entities.
type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for Name entity.
func (n Name) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(n.Formatted); err != nil {
		return nil, err
	}
	if err := enc.Encode(n.FamilyName); err != nil {
		return nil, err
	}
	if err := enc.Encode(n.GivenName); err != nil {
		return nil, err
	}
	if err := enc.Encode(n.MiddleName); err != nil {
		return nil, err
	}
	if err := enc.Encode(n.HonorificPrefix); err != nil {
		return nil, err
	}
	if err := enc.Encode(n.HonorificSuffix); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Name entity.
func (n *Name) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&n.Formatted); err != nil {
		return err
	}
	if err := dec.Decode(&n.FamilyName); err != nil {
		return err
	}
	if err := dec.Decode(&n.GivenName); err != nil {
		return err
	}
	if err := dec.Decode(&n.MiddleName); err != nil {
		return err
	}
	if err := dec.Decode(&n.HonorificPrefix); err != nil {
		return err
	}
	if err := dec.Decode(&n.HonorificSuffix); err != nil {
		return err
	}

	return nil
}

// Email represent an email entity
type Email struct {
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for Email entity.
func (e Email) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(e.Value); err != nil {
		return nil, err
	}
	if err := enc.Encode(e.Type); err != nil {
		return nil, err
	}
	if err := enc.Encode(e.Primary); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Email entity.
func (e *Email) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&e.Value); err != nil {
		return err
	}
	if err := dec.Decode(&e.Type); err != nil {
		return err
	}
	if err := dec.Decode(&e.Primary); err != nil {
		return err
	}

	return nil
}

// Addresses represent an address entity
type Address struct {
	Type          string `json:"type,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Country       string `json:"country,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for Address entity.
func (a Address) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(a.Type); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.Formatted); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.StreetAddress); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.Locality); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.Region); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.PostalCode); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.Country); err != nil {
		return nil, err
	}
	if err := enc.Encode(a.Primary); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Address entity.
func (a *Address) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&a.Type); err != nil {
		return err
	}
	if err := dec.Decode(&a.Formatted); err != nil {
		return err
	}
	if err := dec.Decode(&a.StreetAddress); err != nil {
		return err
	}
	if err := dec.Decode(&a.Locality); err != nil {
		return err
	}
	if err := dec.Decode(&a.Region); err != nil {
		return err
	}
	if err := dec.Decode(&a.PostalCode); err != nil {
		return err
	}
	if err := dec.Decode(&a.Country); err != nil {
		return err
	}
	if err := dec.Decode(&a.Primary); err != nil {
		return err
	}

	return nil
}

type PhoneNumber struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for PhoneNumber entity.
func (pn PhoneNumber) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(pn.Value); err != nil {
		return nil, err
	}
	if err := enc.Encode(pn.Type); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for PhoneNumber entity.
func (pn *PhoneNumber) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&pn.Value); err != nil {
		return err
	}
	if err := dec.Decode(&pn.Type); err != nil {
		return err
	}

	return nil
}

type Manager struct {
	Value string `json:"value,omitempty"`
	Ref   string `json:"$ref,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for Manager entity.
func (m Manager) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(m.Value); err != nil {
		return nil, err
	}
	if err := enc.Encode(m.Ref); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Manager entity.
func (m *Manager) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&m.Value); err != nil {
		return err
	}
	if err := dec.Decode(&m.Ref); err != nil {
		return err
	}

	return nil
}

type EnterpriseData struct {
	EmployeeNumber string   `json:"employeeNumber,omitempty"`
	CostCenter     string   `json:"costCenter,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Division       string   `json:"division,omitempty"`
	Department     string   `json:"department,omitempty"`
	Title          string   `json:"title,omitempty"`
	Manager        *Manager `json:"manager,omitempty"`
	Primary        bool     `json:"primary,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for EnterpriseData entity.
func (ed EnterpriseData) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(ed.EmployeeNumber); err != nil {
		return nil, err
	}
	if err := enc.Encode(ed.CostCenter); err != nil {
		return nil, err
	}
	if err := enc.Encode(ed.Organization); err != nil {
		return nil, err
	}
	if err := enc.Encode(ed.Division); err != nil {
		return nil, err
	}
	if err := enc.Encode(ed.Department); err != nil {
		return nil, err
	}
	if err := enc.Encode(ed.Title); err != nil {
		return nil, err
	}

	if ed.Manager != nil {
		if err := enc.Encode(ed.Manager); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for EnterpriseData entity.
func (ed *EnterpriseData) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&ed.EmployeeNumber); err != nil {
		return err
	}
	if err := dec.Decode(&ed.CostCenter); err != nil {
		return err
	}
	if err := dec.Decode(&ed.Organization); err != nil {
		return err
	}
	if err := dec.Decode(&ed.Division); err != nil {
		return err
	}
	if err := dec.Decode(&ed.Department); err != nil {
		return err
	}
	if err := dec.Decode(&ed.Title); err != nil {
		return err
	}

	if err := dec.Decode(&ed.Manager); err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}

	return nil
}

// User represents a user entity.
type User struct {
	IPID              string          `json:"ipid,omitempty"`
	SCIMID            string          `json:"scimid,omitempty"`
	UserName          string          `json:"userName,omitempty"`
	DisplayName       string          `json:"displayName,omitempty"`
	NickName          string          `json:"nickName,omitempty"`
	ProfileURL        string          `json:"profileURL,omitempty"`
	Title             string          `json:"title,omitempty"`
	UserType          string          `json:"userType,omitempty"`
	PreferredLanguage string          `json:"preferredLanguage,omitempty"`
	Locale            string          `json:"locale,omitempty"`
	Timezone          string          `json:"timezone,omitempty"`
	HashCode          string          `json:"hashCode,omitempty"`
	Emails            []Email         `json:"emails,omitempty"`
	Addresses         []Address       `json:"addresses,omitempty"`
	PhoneNumbers      []PhoneNumber   `json:"phoneNumbers,omitempty"`
	Name              *Name           `json:"name,omitempty"`
	EnterpriseData    *EnterpriseData `json:"enterpriseData,omitempty"`
	Active            bool            `json:"active,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for User entity.
// This is necessary to avoid include the value in the field SCIMID and hashcode until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (u User) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(u.IPID); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.UserName); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.DisplayName); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.NickName); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.ProfileURL); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Title); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.UserType); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.PreferredLanguage); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Locale); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Timezone); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Active); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Emails); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.Addresses); err != nil {
		return nil, err
	}
	if err := enc.Encode(u.PhoneNumbers); err != nil {
		return nil, err
	}

	if u.Name != nil {
		if err := enc.Encode(u.Name); err != nil {
			return nil, err
		}
	}

	if u.EnterpriseData != nil {
		if err := enc.Encode(u.EnterpriseData); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for User entity.
func (u *User) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&u.IPID); err != nil {
		return err
	}
	if err := dec.Decode(&u.UserName); err != nil {
		return err
	}
	if err := dec.Decode(&u.DisplayName); err != nil {
		return err
	}
	if err := dec.Decode(&u.NickName); err != nil {
		return err
	}
	if err := dec.Decode(&u.ProfileURL); err != nil {
		return err
	}
	if err := dec.Decode(&u.Title); err != nil {
		return err
	}
	if err := dec.Decode(&u.UserType); err != nil {
		return err
	}
	if err := dec.Decode(&u.PreferredLanguage); err != nil {
		return err
	}
	if err := dec.Decode(&u.Locale); err != nil {
		return err
	}
	if err := dec.Decode(&u.Timezone); err != nil {
		return err
	}
	if err := dec.Decode(&u.Active); err != nil {
		return err
	}
	if err := dec.Decode(&u.Emails); err != nil {
		return err
	}
	if err := dec.Decode(&u.Addresses); err != nil {
		return err
	}
	if err := dec.Decode(&u.PhoneNumbers); err != nil {
		return err
	}

	// when the user has pointer to Name, but the Name is nil, the gob decoder returns an error
	if err := dec.Decode(&u.Name); err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}

	// when the user has pointer to EnterpriseData, but the Name is nil, the gob decoder returns an error
	if err := dec.Decode(&u.EnterpriseData); err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}

	return nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (u *User) SetHashCode() {
	u.HashCode = Hash(u)
}

// GetPrimaryEmailAddress returns the primary email address of the user.
func (u *User) GetPrimaryEmailAddress() string {
	for _, email := range u.Emails {
		if email.Primary {
			return email.Value
		}
	}
	return ""
}

// UsersResult represents a user result list entity.
type UsersResult struct {
	Items     int     `json:"items"`
	HashCode  string  `json:"hashCode"`
	Resources []*User `json:"resources"`
}

// MarshalBinary implements the gob.GobEncoder interface for UsersResult entity.
func (ur UsersResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(ur.Items); err != nil {
		return nil, err
	}

	if ur.Resources != nil {
		if err := enc.Encode(ur.Resources); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for UsersResult entity.
func (ur *UsersResult) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&ur.Items); err != nil {
		return err
	}

	if err := dec.Decode(&ur.Resources); err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for UsersResult entity.
func (ur *UsersResult) MarshalJSON() ([]byte, error) {
	if ur.Resources == nil {
		ur.Resources = make([]*User, 0)
	}

	return json.MarshalIndent(*ur, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (ur *UsersResult) SetHashCode() {
	// this copy is necessary to avoid changing the original data
	// with the sort.Slice function and always be consistent
	// when calculating the hash code
	c := make([]*User, len(ur.Resources))
	copy(c, ur.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &UsersResult{
		Items:     ur.Items,
		Resources: c,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	ur.HashCode = Hash(copyStruct)
}
