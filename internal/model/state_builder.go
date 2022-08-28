package model

// StateBuilderChoice is used to build a State entity and ensure the calculated hash code is set.
type StateBuilderChoice struct {
	s *State
}

// StateBuilder creates a new StateBuilderChoice entity.
func StateBuilder() *StateBuilderChoice {
	return &StateBuilderChoice{
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
func (b *StateBuilderChoice) WithSchemaVersion(schemaVersion string) *StateBuilderChoice {
	b.s.SchemaVersion = schemaVersion
	return b
}

// WithCodeVersion sets the CodeVersion field of the State entity.
func (b *StateBuilderChoice) WithCodeVersion(codeVersion string) *StateBuilderChoice {
	b.s.CodeVersion = codeVersion
	return b
}

// WithLastSync sets the LastSync field of the State entity.
func (b *StateBuilderChoice) WithLastSync(lastSync string) *StateBuilderChoice {
	b.s.LastSync = lastSync
	return b
}

// WithGroups sets the Groups field of the StateResources entity inside the State entity.
func (b *StateBuilderChoice) WithGroups(groups *GroupsResult) *StateBuilderChoice {
	b.s.Resources.Groups = groups
	return b
}

// WithUsers sets the Users field of the StateResources entity inside the State entity.
func (b *StateBuilderChoice) WithUsers(users *UsersResult) *StateBuilderChoice {
	b.s.Resources.Users = users
	return b
}

// WithGroupsMembers sets the GroupsMembers field of the StateResources entity inside the State entity.
func (b *StateBuilderChoice) WithGroupsMembers(groupsMembers *GroupsMembersResult) *StateBuilderChoice {
	b.s.Resources.GroupsMembers = groupsMembers
	return b
}

// Build returns the State entity.
func (b *StateBuilderChoice) Build() *State {
	s := b.s
	s.SetHashCode()
	return s
}
