package model

import "encoding/json"

type State struct {
	LastSync  string         `json:"lastSync"`
	Resources StateResources `json:"resources"`
}

type StateResources struct {
	Groups      *GroupsResult      `json:"groups"`
	Users       *UsersResult       `json:"users"`
	GroupsUsers *GroupsUsersResult `json:"groupsUsers"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	if s.Resources.Groups == nil {
		s.Resources.Groups = &GroupsResult{}
	}
	if s.Resources.Users == nil {
		s.Resources.Users = &UsersResult{}
	}
	if s.Resources.GroupsUsers == nil {
		s.Resources.GroupsUsers = &GroupsUsersResult{}
	}
	return json.Marshal(*s)
}
