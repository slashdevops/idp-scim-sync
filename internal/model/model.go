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

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// GobEncode implements the gob.GobEncoder interface for Group entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (g *Group) GobEncode() ([]byte, error) {
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
// only fields coming from the Identity Provider are used.
func (g *Group) SetHashCode() {
	g.HashCode = Hash(g)
}

// GroupsResult represents a group result list entity.
type GroupsResult struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Resources []*Group `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsResult entity.
func (gr *GroupsResult) MarshalJSON() ([]byte, error) {
	if gr.Resources == nil {
		gr.Resources = make([]*Group, 0)
	}
	return json.MarshalIndent(*gr, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gr *GroupsResult) SetHashCode() {
	copyResources := make([]*Group, len(gr.Resources))
	copy(copyResources, gr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &GroupsResult{
		Items:     gr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	gr.HashCode = Hash(copyStruct)
}

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	HashCode string `json:"hashCode"`
}

// GobEncode implements the gob.GobEncoder interface for Member entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (m *Member) GobEncode() ([]byte, error) {
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
// only fields coming from the Identity Provider are used.
func (m *Member) SetHashCode() {
	m.HashCode = Hash(m)
}

// MembersResult represents a member result list entity.
type MembersResult struct {
	Items     int       `json:"items"`
	HashCode  string    `json:"hashCode"`
	Resources []*Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (mr *MembersResult) SetHashCode() {
	copyResources := make([]*Member, len(mr.Resources))
	copy(copyResources, mr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &MembersResult{
		Items:     mr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].IPID < copyStruct.Resources[j].IPID
	})

	mr.HashCode = Hash(copyStruct)
}

// GroupMembers represents a group members entity.
type GroupMembers struct {
	Items     int       `json:"items"`
	HashCode  string    `json:"hashCode,omitempty"`
	Group     *Group    `json:"group"`
	Resources []*Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gm *GroupMembers) SetHashCode() {
	copyResources := make([]*Member, len(gm.Resources))
	copy(copyResources, gm.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &GroupMembers{
		Items:     gm.Items,
		Group:     gm.Group,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].Email < copyStruct.Resources[j].Email
	})

	gm.HashCode = Hash(copyStruct)
}

// GroupsMembersResult represents a group members result list entity.
type GroupsMembersResult struct {
	Items     int             `json:"items"`
	HashCode  string          `json:"hashCode"`
	Resources []*GroupMembers `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsMembersResult entity.
func (gmr *GroupsMembersResult) MarshalJSON() ([]byte, error) {
	if gmr.Resources == nil {
		gmr.Resources = make([]*GroupMembers, 0)
	}
	return json.MarshalIndent(*gmr, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gmr *GroupsMembersResult) SetHashCode() {
	copyResources := make([]*GroupMembers, len(gmr.Resources))
	copy(copyResources, gmr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := GroupsMembersResult{
		Items:     gmr.Items,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	gmr.HashCode = Hash(copyStruct)
}
