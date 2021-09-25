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
