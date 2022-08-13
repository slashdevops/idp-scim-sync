package model

// stateBuilder is used to build a State entity and ensure the calculated hash code is set.
type stateBuilder struct {
	s *State
}

// StateBuilder creates a new stateBuilder entity.
func StateBuilder() *stateBuilder {
	return &stateBuilder{
		s: &State{
			SchemaVersion: StateSchemaVersion,
			Resources: &StateResources{
				Groups: &GroupsResult{
					Resources: make([]*Group, 0),
				},
				Users: &UsersResult{
					Resources: make([]*User, 0),
				},
				GroupsMembers: &GroupsMembersResult{
					Resources: make([]*GroupMembers, 0),
				},
			},
		},
	}
}

// WithSchemaVersion sets the SchemaVersion field of the State entity.
func (b *stateBuilder) WithSchemaVersion(schemaVersion string) *stateBuilder {
	b.s.SchemaVersion = schemaVersion
	return b
}

// WithCodeVersion sets the CodeVersion field of the State entity.
func (b *stateBuilder) WithCodeVersion(codeVersion string) *stateBuilder {
	b.s.CodeVersion = codeVersion
	return b
}

// WithLastSync sets the LastSync field of the State entity.
func (b *stateBuilder) WithLastSync(lastSync string) *stateBuilder {
	b.s.LastSync = lastSync
	return b
}

// WithGroups sets the Groups field of the StateResources entity inside the State entity.
func (b *stateBuilder) WithGroups(groups *GroupsResult) *stateBuilder {
	b.s.Resources.Groups = groups
	return b
}

// WithUsers sets the Users field of the StateResources entity inside the State entity.
func (b *stateBuilder) WithUsers(users *UsersResult) *stateBuilder {
	b.s.Resources.Users = users
	return b
}

// WithGroupsMembers sets the GroupsMembers field of the StateResources entity inside the State entity.
func (b *stateBuilder) WithGroupsMembers(groupsMembers *GroupsMembersResult) *stateBuilder {
	b.s.Resources.GroupsMembers = groupsMembers
	return b
}

// Build returns the State entity.
func (b *stateBuilder) Build() *State {
	s := b.s
	s.SetHashCode()
	return s
}
