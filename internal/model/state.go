package model

import "encoding/json"

type State struct {
	LastSync  string         `json:"lastSync,omitempty"`
	HashCode  string         `json:"hashCode"`
	Resources StateResources `json:"resources,omitempty"`
}

type StateResources struct {
	Groups      GroupsResult      `json:"groups,omitempty"`
	Users       UsersResult       `json:"users,omitempty"`
	GroupsUsers GroupsUsersResult `json:"groupsUsers,omitempty"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

type StateMetadataResources struct {
	GroupsLocation      string `json:"groupsLocation"`
	UsersLocation       string `json:"usersLocation"`
	GroupsUsersLocation string `json:"groupsUsersLocation"`
}

type StateMetadata struct {
	LastSync  string                 `json:"lastSync"`
	HashCode  string                 `json:"hashCode"`
	Resources StateMetadataResources `json:"resources"`
}
