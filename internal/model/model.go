package model

import (
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

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (u *User) SetHashCode() {
	u.HashCode = hash.Get(User{
		IPID:        u.IPID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Active:      u.Active,
		Email:       u.Email,
	})
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
	sort.Slice(ur.Resources, func(i, j int) bool {
		return ur.Resources[i].HashCode < ur.Resources[j].HashCode
	})

	ur.HashCode = hash.Get(UsersResult{
		Items:     ur.Items,
		Resources: ur.Resources,
	})
}

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (g *Group) SetHashCode() {
	g.HashCode = hash.Get(Group{
		IPID:  g.IPID,
		Name:  g.Name,
		Email: g.Email,
	})
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
	sort.Slice(gr.Resources, func(i, j int) bool {
		return gr.Resources[i].HashCode < gr.Resources[j].HashCode
	})

	gr.HashCode = hash.Get(GroupsResult{
		Items:     gr.Items,
		Resources: gr.Resources,
	})
}

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (m *Member) SetHashCode() {
	m.HashCode = hash.Get(Member{
		IPID:  m.IPID,
		Email: m.Email,
	})
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
	sort.Slice(mr.Resources, func(i, j int) bool {
		return mr.Resources[i].HashCode < mr.Resources[j].HashCode
	})

	mr.HashCode = hash.Get(MembersResult{
		Items:     mr.Items,
		Resources: mr.Resources,
	})
}

// GroupMembers represents a group members entity.
type GroupMembers struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Group     Group    `json:"group"`
	Resources []Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields comming from the Identity Provider are used.
func (gm *GroupMembers) SetHashCode() {
	sort.Slice(gm.Resources, func(i, j int) bool {
		return gm.Resources[i].HashCode < gm.Resources[j].HashCode
	})

	gm.HashCode = hash.Get(GroupMembers{
		Items:     gm.Items,
		Group:     gm.Group,
		Resources: gm.Resources,
	})
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
	sort.Slice(gmr.Resources, func(i, j int) bool {
		return gmr.Resources[i].HashCode < gmr.Resources[j].HashCode
	})

	gmr.HashCode = hash.Get(GroupsMembersResult{
		Items:     gmr.Items,
		Resources: gmr.Resources,
	})
}
