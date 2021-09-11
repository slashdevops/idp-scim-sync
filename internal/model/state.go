package model

import "encoding/json"

type State struct {
	Name          string         `json:"name"`
	SchemaVersion string         `json:"version"`
	LastSync      string         `json:"lastSync"`
	HashCode      string         `json:"hashCode"`
	Resources     StateResources `json:"resources"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

type StateResources struct {
	StoreGroupsResult      *StoreGroupsResult      `json:"storeGroupsResult"`
	StoreUsersResult       *StoreUsersResult       `json:"storeUsersResult"`
	StoreGroupsUsersResult *StoreGroupsUsersResult `json:"storeGroupsUsersResult"`
}

type StoreGroupsResult struct {
	Location string
}

type StoreUsersResult struct {
	Location string
}

type StoreGroupsUsersResult struct {
	Location string
}

type StoreStateResult struct {
	Location string
}
