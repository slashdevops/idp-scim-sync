package sync

type Group struct {
	Name  string
	Email string
}

type User struct {
	Name  string
	Email string
}

type Member struct {
	Name  string
	Email string
}

type GroupResult struct {
	Total     int
	Items     int
	NextItem  int
	Resources []*Group
}

type UserResult struct {
	Total     int
	Items     int
	NextItem  int
	Resources []*User
}

type MemberResult struct {
	Total     int
	Items     int
	NextItem  int
	Resources []*Member
}
