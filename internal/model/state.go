package model

import "encoding/json"

const (
	// StateSchemaVersion is the current schema version for the state file.
	StateSchemaVersion = "1.0.0"
)

// StateResources is a list of resources in the state, groups, users and groups and their users.
type StateResources struct {
	Groups        *GroupsResult        `json:"groups"`
	Users         *UsersResult         `json:"users"`
	GroupsMembers *GroupsMembersResult `json:"groupsMembers"`
}

// State is the state of the system.
type State struct {
	SchemaVersion string          `json:"schemaVersion"`
	CodeVersion   string          `json:"codeVersion"`
	LastSync      string          `json:"lastSync"`
	HashCode      string          `json:"hashCode"`
	Resources     *StateResources `json:"resources"`
}

// MarshalJSON marshals the State to JSON.
func (s *State) MarshalJSON() ([]byte, error) {
	if s.Resources == nil {
		s.Resources = &StateResources{}
	}
	if s.Resources.Groups == nil {
		s.Resources.Groups = &GroupsResult{}
	}
	if s.Resources.Groups.Resources == nil {
		s.Resources.Groups.Resources = make([]*Group, 0)
	}

	if s.Resources.Users == nil {
		s.Resources.Users = &UsersResult{}
	}
	if s.Resources.Users.Resources == nil {
		s.Resources.Users.Resources = make([]*User, 0)
	}

	if s.Resources.GroupsMembers == nil {
		s.Resources.GroupsMembers = &GroupsMembersResult{}
	}
	if s.Resources.GroupsMembers.Resources == nil {
		s.Resources.GroupsMembers.Resources = make([]*GroupMembers, 0)
	}

	for mbrs := range s.Resources.GroupsMembers.Resources {
		if s.Resources.GroupsMembers.Resources[mbrs] == nil {
			s.Resources.GroupsMembers.Resources[mbrs] = &GroupMembers{}
		}

		if s.Resources.GroupsMembers.Resources[mbrs].Resources == nil {
			s.Resources.GroupsMembers.Resources[mbrs].Resources = make([]*Member, 0)
		}

		for member := range s.Resources.GroupsMembers.Resources[mbrs].Resources {
			if s.Resources.GroupsMembers.Resources[mbrs].Resources[member] == nil {
				s.Resources.GroupsMembers.Resources[mbrs].Resources[member] = &Member{}
			}
		}
	}

	return json.MarshalIndent(*s, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (s *State) SetHashCode() {
	// we need to do a deep copy of the state struct to avoid SCIMID in the hash calculation
	// because every time the idp data is compared with the state data, the SCIMID doesn't compute in the hash

	if s.Resources == nil {
		s.Resources = &StateResources{
			Groups: &GroupsResult{
				Resources: make([]*Group, 0),
			},
			Users: &UsersResult{
				Resources: make([]*User, 0),
			},
			GroupsMembers: &GroupsMembersResult{
				Resources: make([]*GroupMembers, 0),
			},
		}
	}

	if s.Resources.Groups == nil {
		s.Resources.Groups = &GroupsResult{}
	}
	groups := make([]*Group, 0)
	for _, group := range s.Resources.Groups.Resources {
		e := GroupBuilder().
			WithIPID(group.IPID).
			WithName(group.Name).
			WithEmail(group.Email).
			Build()

		groups = append(groups, e)
	}

	groupsResult := GroupsResultBuilder().WithResources(groups).Build()

	if s.Resources.Users == nil {
		s.Resources.Users = &UsersResult{}
	}
	users := make([]*User, 0)
	for _, user := range s.Resources.Users.Resources {
		e := UserBuilder().
			WithIPID(user.IPID).
			WithName(user.Name).
			WithDisplayName(user.DisplayName).
			WithEmail(
				EmailBuilder().
					WithValue(user.GetPrimaryEmailAddress()).
					WithType("work").
					WithPrimary(true).
					Build(),
			).
			WithActive(user.Active).
			Build()

		users = append(users, e)
	}
	usersResult := UsersResultBuilder().WithResources(users).Build()

	if s.Resources.GroupsMembers == nil {
		s.Resources.GroupsMembers = &GroupsMembersResult{}
	}
	groupsMembers := make([]*GroupMembers, 0)
	for _, groupMembers := range s.Resources.GroupsMembers.Resources {
		group := GroupBuilder().
			WithIPID(groupMembers.Group.IPID).
			WithName(groupMembers.Group.Name).
			WithEmail(groupMembers.Group.Email).
			Build()

		members := make([]*Member, 0)
		for _, member := range groupMembers.Resources {
			m := MemberBuilder().
				WithIPID(member.IPID).
				WithEmail(member.Email).
				Build()

			members = append(members, m)
		}

		e := GroupMembersBuilder().WithGroup(group).WithResources(members).Build()
		groupsMembers = append(groupsMembers, e)
	}

	groupsMembersResult := GroupsMembersResultBuilder().WithResources(groupsMembers).Build()

	// The hash code of the state only depends on Resources changes not in metadata changes.
	copyState := State{
		Resources: &StateResources{
			Groups:        groupsResult,
			Users:         usersResult,
			GroupsMembers: groupsMembersResult,
		},
	}

	s.HashCode = Hash(copyState)
}
