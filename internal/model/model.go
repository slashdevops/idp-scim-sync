package model

import (
	"encoding/json"
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

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
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

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

// MembersResult represents a member result list entity.
type MembersResult struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Resources []Member `json:"resources"`
}

// GroupMembers represents a group members entity.
type GroupMembers struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Group     Group    `json:"group"`
	Resources []Member `json:"resources"`
}

// GroupsMembersResult represents a group members result list entity.
type GroupsMembersResult struct {
	Items     int            `json:"items"`
	HashCode  string         `json:"hashCode"`
	Resources []GroupMembers `json:"resources"`
}

// GroupUsers represents a group users entity.
type GroupUsers struct {
	Items     int    `json:"items"`
	HashCode  string `json:"hashCode"`
	Group     Group  `json:"group"`
	Resources []User `json:"resources"`
}

// GroupsUsersResult represents a group users result list entity.
type GroupsUsersResult struct {
	Items     int          `json:"items"`
	HashCode  string       `json:"hashCode"`
	Resources []GroupUsers `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsUsersResult entity.
func (gur *GroupsUsersResult) MarshalJSON() ([]byte, error) {
	if gur.Resources == nil {
		gur.Resources = make([]GroupUsers, 0)
	}
	return json.MarshalIndent(*gur, "", "  ")
}
