package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
)

// Name represents a name entity and is used in other entities.
type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

// User represents a user entity.
type User struct {
	IPID        string `json:"ipid"`
	SCIMID      string `json:"scimid"`
	Name        Name   `json:"name"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Email       string `json:"email"`
	HashCode    string `json:"hashCode"`
}

// GobEncode implements the gob.GobEncoder interface for User entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (u *User) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(u.IPID); err != nil {
		panic(err)
	}
	if err := enc.Encode(u.Name); err != nil {
		panic(err)
	}
	if err := enc.Encode(u.DisplayName); err != nil {
		panic(err)
	}
	if err := enc.Encode(u.Active); err != nil {
		panic(err)
	}
	if err := enc.Encode(u.Email); err != nil {
		panic(err)
	}
	return buf.Bytes(), nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (u *User) SetHashCode() {
	u.HashCode = Hash(u)
}

// UsersResult represents a user result list entity.
type UsersResult struct {
	Items     int     `json:"items"`
	HashCode  string  `json:"hashCode"`
	Resources []*User `json:"resources"`
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
	copyResources := make([]*User, len(ur.Resources))
	copy(copyResources, ur.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &UsersResult{
		Items:     ur.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	ur.HashCode = Hash(copyStruct)
}
