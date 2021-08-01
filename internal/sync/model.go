package sync

type Group struct {
	ID          string
	ExternalId  string
	Name        string
	Email       string
	Description string
}

type User struct {
	ID         string
	ExternalId string
	Name       struct {
		FamilyName string
		GivenName  string
	}
	DisplayName string
	Active      bool
	Email       string
}

type Member struct {
	ID         string
	ExternalId string
	Email      string
}

type GroupsResult struct {
	Items     int
	Resources []*Group
}

type UsersResult struct {
	Items     int
	Resources []*User
}

type MembersResult struct {
	Items     int
	Resources []*Member
}
