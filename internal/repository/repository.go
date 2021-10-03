package repository

const (
	StateIndexName       = "state.json"
	GroupsIndexName      = "groups.ndjson"
	GroupsMetaName       = "groups_meta.json"
	UsersIndexName       = "users.ndjson"
	UsersMetaName        = "users_meta.json"
	GroupsUsersIndexName = "groups_users.ndjson"
	GroupsUsersMetaName  = "groups_users_meta.json"
)

type StateMetaIndexResources struct {
	GroupsLocation      string `json:"groupsLocation"`
	UsersLocation       string `json:"usersLocation"`
	GroupsUsersLocation string `json:"groupsUsersLocation"`
}

type StateMetaIndex struct {
	LastSync  string                  `json:"lastSync"`
	HashCode  string                  `json:"hashCode"`
	Resources StateMetaIndexResources `json:"resources"`
}

type GroupsMetaIndex struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}

type UsersMetaIndex struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}

type GroupsUsersMetaIndex struct {
	Items    int    `json:"items"`
	HashCode string `json:"hashCode"`
	Location string `json:"location"`
}
