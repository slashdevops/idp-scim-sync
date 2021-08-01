package sync

type Id struct {
	IdentityProvider string
	SCIM             string
}

type Name struct {
	FamilyName string
	GivenName  string
}

type Group struct {
	Id    Id
	Name  string
	Email string
}

type User struct {
	Id          Id
	Name        Name
	DisplayName string
	Active      bool
	Email       string
}

type Member struct {
	Id    Id
	Email string
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
