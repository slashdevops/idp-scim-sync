package sync

type Name struct {
	FamilyName string
	GivenName  string
}

type User struct {
	Id          string
	Name        Name
	DisplayName string
	Active      bool
	Email       string
}

type UsersResult struct {
	Items     int
	Resources []*User
}

type Group struct {
	Id    string
	Name  string
	Email string
}

type GroupsResult struct {
	Items     int
	Resources []*Group
}

type Member struct {
	Id    string
	Email string
}
type MembersResult struct {
	Items     int
	Resources []*Member
}

type GroupsMembers map[string][]*Member

func (gms GroupsMembers) GetMembers(groupId string) []*Member {
	return gms[groupId]
}

type GroupsMembersResult struct {
	Items     int
	Resources *GroupsMembers
}
