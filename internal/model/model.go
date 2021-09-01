package model

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
}

type UsersResult struct {
	Items     int     `json:"items"`
	Resources []*User `json:"resources"`
}

type Group struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type StoreGroupsResult struct {
	Place string
}

type StoreUsersResult struct {
	Place string
}

type StoreGroupsMembersResult struct {
	Place string
}

type StoreGroupsUsersResult struct {
	Place string
}

type StoreStateResult struct {
	Place string
}

type GroupsResult struct {
	Items     int      `json:"items"`
	Resources []*Group `json:"resources"`
}

type Member struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type MembersResult struct {
	Items     int       `json:"items"`
	Resources []*Member `json:"resources"`
}

// type GroupsMembers map[string][]*Member

// func (gms GroupsMembers) GetMembers(groupID string) []*Member {
// 	return gms[groupID]
// }

type GroupMembers struct {
	Items     int       `json:"items"`
	Group     Group     `json:"group"`
	Resources []*Member `json:"resources"`
}

type GroupsMembersResult struct {
	Items     int             `json:"items"`
	Resources []*GroupMembers `json:"resources"`
}

type GroupUsers struct {
	Items     int     `json:"items"`
	Group     Group   `json:"group"`
	Resources []*User `json:"resources"`
}

type GroupsUsersResult struct {
	Items     int           `json:"items"`
	Resources []*GroupUsers `json:"resources"`
}

type State struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	HashCode      string `json:"hashcode"`
	Groups        StoreGroupsResult
	Users         StoreUsersResult
	GroupsMembers StoreGroupsMembersResult
}
