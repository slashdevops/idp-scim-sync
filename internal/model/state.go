package model

import "encoding/json"

const (
	StateSchemaVersion = "1.0.0"
)

type StateResources struct {
	Groups      GroupsResult      `json:"groups"`
	Users       UsersResult       `json:"users"`
	GroupsUsers GroupsUsersResult `json:"groupsUsers"`
}

type State struct {
	SchemaVersion string         `json:"schemaVersion"`
	CodeVersion   string         `json:"codeVersion"`
	LastSync      string         `json:"lastSync"`
	HashCode      string         `json:"hashCode"`
	Resources     StateResources `json:"resources"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(*s, "", "  ")
}
