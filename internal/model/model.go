package model

import (
	"encoding/json"
)

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type User struct {
	ID          string `json:"id"`
	Name        Name   `json:"name"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Email       string `json:"email"`
	HashCode    string `json:"hashCode"`
}

type UsersResult struct {
	Items     int    `json:"items"`
	HashCode  string `json:"hashCode"`
	Resources []User `json:"resources"`
}

func (ur *UsersResult) MarshalJSON() ([]byte, error) {
	if ur.Resources == nil {
		ur.Resources = make([]User, 0)
	}
	return json.Marshal(*ur)
}

type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

type GroupsResult struct {
	Items     int     `json:"items"`
	HashCode  string  `json:"hashCode"`
	Resources []Group `json:"resources"`
}

func (gr *GroupsResult) MarshalJSON() ([]byte, error) {
	if gr.Resources == nil {
		gr.Resources = make([]Group, 0)
	}
	return json.Marshal(*gr)
}

type Member struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	HashCode string `json:"hashCode"`
}

type MembersResult struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Resources []Member `json:"resources"`
}

type GroupMembers struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Group     Group    `json:"group"`
	Resources []Member `json:"resources"`
}

type GroupsMembersResult struct {
	Items     int            `json:"items"`
	HashCode  string         `json:"hashCode"`
	Resources []GroupMembers `json:"resources"`
}

type GroupUsers struct {
	Items     int    `json:"items"`
	HashCode  string `json:"hashCode"`
	Group     Group  `json:"group"`
	Resources []User `json:"resources"`
}

type GroupsUsersResult struct {
	Items     int          `json:"items"`
	HashCode  string       `json:"hashCode"`
	Resources []GroupUsers `json:"resources"`
}

func (gur *GroupsUsersResult) MarshalJSON() ([]byte, error) {
	if gur.Resources == nil {
		gur.Resources = make([]GroupUsers, 0)
	}
	return json.Marshal(*gur)
}

type GroupsMetadata struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}

type UsersMetadata struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}

type GroupsUsersMetadata struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}
