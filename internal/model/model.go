package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
)

// Name represents a name entity.
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
// the hash.Get function use gob to calculate the hash code.
func (u User) GobEncode() ([]byte, error) {
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
// only fields comming from the Identity Provider are used.
func (u *User) SetHashCode() {
	u.HashCode = hash.Get(u)
}

// UsersResult represents a user result list entity.
type UsersResult struct {
	Items     int    `json:"items"`
	HashCode  string `json:"hashCode"`
	Resources []User `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for UsersResult entity.
func (ur *UsersResult) MarshalJSON() ([]byte, error) {
	if ur.Resources == nil {
		ur.Resources = make([]User, 0)
	}
	return json.MarshalIndent(*ur, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (ur *UsersResult) SetHashCode() {
	if len(ur.Resources) == 0 {
		return
	}

	copyResources := make([]User, len(ur.Resources))
	copy(copyResources, ur.Resources)

	// only these fields are used in the hash calculation
	copyOfStruct := UsersResult{
		Items:     ur.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyOfStruct.Resources, func(i, j int) bool {
		return copyOfStruct.Resources[i].HashCode < copyOfStruct.Resources[j].HashCode
	})

	ur.HashCode = hash.Get(copyOfStruct)
}

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// GobEncode implements the gob.GobEncoder interface for User entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the hash.Get function use gob to calculate the hash code.
func (g Group) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(g.IPID); err != nil {
		panic(err)
	}
	if err := enc.Encode(g.Name); err != nil {
		panic(err)
	}
	if err := enc.Encode(g.Email); err != nil {
		panic(err)
	}
	return buf.Bytes(), nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (g *Group) SetHashCode() {
	g.HashCode = hash.Get(g)
}

// GroupsResult represents a group result list entity.
type GroupsResult struct {
	Items     int     `json:"items"`
	HashCode  string  `json:"hashCode"`
	Resources []Group `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsResult entity.
func (gr *GroupsResult) MarshalJSON() ([]byte, error) {
	if gr.Resources == nil {
		gr.Resources = make([]Group, 0)
	}
	return json.MarshalIndent(*gr, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (gr *GroupsResult) SetHashCode() {
	if len(gr.Resources) == 0 {
		return
	}

	copyResources := make([]Group, len(gr.Resources))
	copy(copyResources, gr.Resources)

	// only these fields are used in the hash calculation
	copyOfStruct := GroupsResult{
		Items:     gr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyOfStruct.Resources, func(i, j int) bool {
		return copyOfStruct.Resources[i].HashCode < copyOfStruct.Resources[j].HashCode
	})

	gr.HashCode = hash.Get(copyOfStruct)
}

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// GobEncode implements the gob.GobEncoder interface for User entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the hash.Get function use gob to calculate the hash code.
func (m Member) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(m.IPID); err != nil {
		panic(err)
	}
	if err := enc.Encode(m.Email); err != nil {
		panic(err)
	}
	return buf.Bytes(), nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (m *Member) SetHashCode() {
	m.HashCode = hash.Get(m)
}

// MembersResult represents a member result list entity.
type MembersResult struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Resources []Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (mr *MembersResult) SetHashCode() {
	if len(mr.Resources) == 0 {
		return
	}

	copyResources := make([]Member, len(mr.Resources))
	copy(copyResources, mr.Resources)

	// only these fields are used in the hash calculation
	copyOfStruct := MembersResult{
		Items:     mr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyOfStruct.Resources, func(i, j int) bool {
		return copyOfStruct.Resources[i].IPID < copyOfStruct.Resources[j].IPID
	})

	mr.HashCode = hash.Get(copyOfStruct)
}

// GroupMembers represents a group members entity.
type GroupMembers struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode,omitempty"`
	Group     Group    `json:"group"`
	Resources []Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (gm *GroupMembers) SetHashCode() {
	if len(gm.Resources) == 0 {
		return
	}

	copyResources := make([]Member, len(gm.Resources))
	copy(copyResources, gm.Resources)

	// only these fields are used in the hash calculation
	copyOfStruct := GroupMembers{
		Items:     gm.Items,
		Group:     gm.Group,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never coulb be empty and it is unique
	sort.Slice(copyOfStruct.Resources, func(i, j int) bool {
		return copyOfStruct.Resources[i].Email < copyOfStruct.Resources[j].Email
	})

	gm.HashCode = hash.Get(copyOfStruct)
}

// GroupsMembersResult represents a group members result list entity.
type GroupsMembersResult struct {
	Items     int            `json:"items"`
	HashCode  string         `json:"hashCode"`
	Resources []GroupMembers `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsMembersResult entity.
func (gur *GroupsMembersResult) MarshalJSON() ([]byte, error) {
	if gur.Resources == nil {
		gur.Resources = make([]GroupMembers, 0)
	}
	return json.MarshalIndent(*gur, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (gmr *GroupsMembersResult) SetHashCode() {
	if len(gmr.Resources) == 0 {
		return
	}

	copyResources := make([]GroupMembers, len(gmr.Resources))
	copy(copyResources, gmr.Resources)

	// only these fields are used in the hash calculation
	copyOfStruct := GroupsMembersResult{
		Items:     gmr.Items,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never coulb be empty and it is unique
	sort.Slice(copyOfStruct.Resources, func(i, j int) bool {
		return copyOfStruct.Resources[i].HashCode < copyOfStruct.Resources[j].HashCode
	})

	gmr.HashCode = hash.Get(copyOfStruct)
}
