package model

import "encoding/json"

const (
	// StateSchemaVersion is the current schema version for the state file.
	StateSchemaVersion = "1.0.0"
)

// StateResources is a list of resources in the state, groups, users and groups and their users.
type StateResources struct {
	Groups        GroupsResult        `json:"groups"`
	Users         UsersResult         `json:"users"`
	GroupsMembers GroupsMembersResult `json:"groupsMembers"`
}

// State is the state of the system.
type State struct {
	SchemaVersion string         `json:"schemaVersion"`
	CodeVersion   string         `json:"codeVersion"`
	LastSync      string         `json:"lastSync"`
	HashCode      string         `json:"hashCode"`
	Resources     StateResources `json:"resources"`
}

// MarshalJSON marshals the State to JSON.
func (s *State) MarshalJSON() ([]byte, error) {
	if s.Resources.Groups.Resources == nil {
		s.Resources.Groups.Resources = make([]*Group, 0)
	}
	if s.Resources.Users.Resources == nil {
		s.Resources.Users.Resources = make([]*User, 0)
	}
	if s.Resources.GroupsMembers.Resources == nil {
		s.Resources.GroupsMembers.Resources = make([]*GroupMembers, 0)
	}

	return json.MarshalIndent(*s, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (s *State) SetHashCode() {
	// we need to do a deep copy of the state struct to avoid SCIMID in the hash calculation
	// because every time the idp data is compared with the state data, the SCIMID doesn't compute in the hash

	groups := make([]*Group, 0)
	for _, group := range s.Resources.Groups.Resources {
		e := &Group{
			IPID:  group.IPID,
			Name:  group.Name,
			Email: group.Email,
		}
		e.SetHashCode()
		groups = append(groups, e)
	}

	groupsResult := GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	groupsResult.SetHashCode()

	users := make([]*User, 0)
	for _, user := range s.Resources.Users.Resources {
		e := &User{
			IPID:        user.IPID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Active:      user.Active,
			Email:       user.Email,
		}
		e.SetHashCode()
		users = append(users, e)
	}
	usersResult := UsersResult{
		Items:     len(users),
		Resources: users,
	}
	usersResult.SetHashCode()

	groupsMembers := make([]*GroupMembers, 0)
	for _, groupMembers := range s.Resources.GroupsMembers.Resources {
		group := Group{
			IPID:  groupMembers.Group.IPID,
			Name:  groupMembers.Group.Name,
			Email: groupMembers.Group.Email,
		}
		group.SetHashCode()

		members := make([]*Member, 0)
		for _, member := range groupMembers.Resources {
			m := &Member{
				IPID:  member.IPID,
				Email: member.Email,
			}
			m.SetHashCode()
			members = append(members, m)
		}

		e := &GroupMembers{
			Items:     len(groupMembers.Resources),
			Group:     group,
			Resources: members,
		}
		e.SetHashCode()
		groupsMembers = append(groupsMembers, e)
	}

	groupsMembersResult := GroupsMembersResult{
		Items:     len(groupsMembers),
		Resources: groupsMembers,
	}
	groupsMembersResult.SetHashCode()

	copyState := State{
		Resources: StateResources{
			Groups:        groupsResult,
			Users:         usersResult,
			GroupsMembers: groupsMembersResult,
		},
	}

	s.HashCode = Hash(copyState)
}
