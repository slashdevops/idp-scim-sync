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

type SyncState struct {
	Version       string `json:"version"`
	Checksum      string `json:"checksum"`
	Groups        StoreGroupsResult
	Users         StoreUsersResult
	GroupsMembers StoreGroupsMembersResult
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

type GroupsMembers map[string][]*Member

func (gms GroupsMembers) GetMembers(groupID string) []*Member {
	return gms[groupID]
}

type GroupsMembersResult struct {
	Items     int            `json:"items"`
	Resources *GroupsMembers `json:"resources"`
}