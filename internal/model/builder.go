package model

// userBuilder is the builder of User entity.
type userBuilder struct {
	u *User
}

// UserBuilder creates a new userBuilder entity.
func UserBuilder() *userBuilder {
	return &userBuilder{
		u: &User{},
	}
}

// WithIPID sets the IPID field of the User entity.
func (b *userBuilder) WithIPID(ipid string) *userBuilder {
	b.u.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the User entity.
func (b *userBuilder) WithSCIMID(scimid string) *userBuilder {
	b.u.SCIMID = scimid
	return b
}

// WithFamilyName sets the FamilyName field of the User entity.
func (b *userBuilder) WithFamilyName(familyName string) *userBuilder {
	b.u.Name.FamilyName = familyName
	return b
}

// WithGivenName sets the GivenName field of the User entity.
func (b *userBuilder) WithGivenName(givenName string) *userBuilder {
	b.u.Name.GivenName = givenName
	return b
}

// WithDisplayName sets the DisplayName field of the User entity.
func (b *userBuilder) WithDisplayName(displayName string) *userBuilder {
	b.u.DisplayName = displayName
	return b
}

// WithActive sets the Active field of the User entity.
func (b *userBuilder) WithActive(active bool) *userBuilder {
	b.u.Active = active
	return b
}

// WithEmail sets the Email field of the User entity.
func (b *userBuilder) WithEmail(email string) *userBuilder {
	b.u.Email = email
	return b
}

// Build returns the User entity.
func (b *userBuilder) Build() *User {
	u := b.u
	u.SetHashCode()
	return u
}

// usersResultBuilder is used to build a UsersResult entity and ensure the calculated hash code and items.
type usersResultBuilder struct {
	ur *UsersResult
}

// UsersResultBuilder creates a new usersResultBuilder entity.
func UsersResultBuilder() *usersResultBuilder {
	return &usersResultBuilder{
		ur: &UsersResult{
			Resources: make([]*User, 0),
		},
	}
}

// WithResources sets the Resources field of the UsersResult entity.
func (b *usersResultBuilder) WithResources(resources []*User) *usersResultBuilder {
	b.ur.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the UsersResult entity.
func (b *usersResultBuilder) WithResource(resource *User) *usersResultBuilder {
	b.ur.Resources = append(b.ur.Resources, resource)
	return b
}

// Build returns the UserResult entity.
func (b *usersResultBuilder) Build() *UsersResult {
	ur := b.ur
	ur.Items = len(ur.Resources)
	ur.SetHashCode()
	return ur
}

// groupBuilder is used to build a Group entity and ensure the calculated hash code is set.
type groupBuilder struct {
	g *Group
}

// GroupBuilder creates a new groupBuilder entity.
func GroupBuilder() *groupBuilder {
	return &groupBuilder{
		g: &Group{},
	}
}

// WithIPID sets the IPID field of the Group entity.
func (b *groupBuilder) WithIPID(ipid string) *groupBuilder {
	b.g.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Group entity.
func (b *groupBuilder) WithSCIMID(scimid string) *groupBuilder {
	b.g.SCIMID = scimid
	return b
}

// WithName sets the Name field of the Group entity.
func (b *groupBuilder) WithName(name string) *groupBuilder {
	b.g.Name = name
	return b
}

// WithEmail sets the Email field of the Group entity.
func (b *groupBuilder) WithEmail(email string) *groupBuilder {
	b.g.Email = email
	return b
}

// Build returns the Group entity.
func (b *groupBuilder) Build() *Group {
	g := b.g
	g.SetHashCode()
	return g
}

// groupsResultBuilder is used to build a GroupsResult entity and ensure the calculated hash code and items is set.
type groupsResultBuilder struct {
	gr *GroupsResult
}

// GroupsResultBuilder creates a new groupsResultBuilder entity.
func GroupsResultBuilder() *groupsResultBuilder {
	return &groupsResultBuilder{
		gr: &GroupsResult{
			Resources: make([]*Group, 0),
		},
	}
}

// WithResources sets the Resources field of the GroupsResult entity.
func (b *groupsResultBuilder) WithResources(resources []*Group) *groupsResultBuilder {
	b.gr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupsResult entity.
func (b *groupsResultBuilder) WithResource(resource *Group) *groupsResultBuilder {
	b.gr.Resources = append(b.gr.Resources, resource)
	return b
}

// Build returns the GroupsResult entity.
func (b *groupsResultBuilder) Build() *GroupsResult {
	gr := b.gr
	gr.Items = len(gr.Resources)
	gr.SetHashCode()
	return gr
}

// memberBuilder is used to build a Member entity and ensure the calculated hash code is set.
type memberBuilder struct {
	m *Member
}

// MemberBuilder creates a new memberBuilder entity.
func MemberBuilder() *memberBuilder {
	return &memberBuilder{
		m: &Member{},
	}
}

// WithIPID sets the IPID field of the Member entity.
func (b *memberBuilder) WithIPID(ipid string) *memberBuilder {
	b.m.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Member entity.
func (b *memberBuilder) WithSCIMID(scimid string) *memberBuilder {
	b.m.SCIMID = scimid
	return b
}

// WithEmail sets the Email field of the Member entity.
func (b *memberBuilder) WithEmail(email string) *memberBuilder {
	b.m.Email = email
	return b
}

// WithStatus sets the Status field of the Member entity.
func (b *memberBuilder) WithStatus(status string) *memberBuilder {
	b.m.Status = status
	return b
}

// Build returns the Member entity.
func (b *memberBuilder) Build() *Member {
	m := b.m
	m.SetHashCode()
	return m
}

// membersResultBuilder is used to build a MembersResult entity and ensure the calculated hash code and items is set.
type membersResultBuilder struct {
	mr *MembersResult
}

// MembersResultBuilder creates a new membersResultBuilder entity.
func MembersResultBuilder() *membersResultBuilder {
	return &membersResultBuilder{
		mr: &MembersResult{
			Resources: make([]*Member, 0),
		},
	}
}

// WithResources sets the Resources field of the MembersResult entity.
func (b *membersResultBuilder) WithResources(resources []*Member) *membersResultBuilder {
	b.mr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the MembersResult entity.
func (b *membersResultBuilder) WithResource(resource *Member) *membersResultBuilder {
	b.mr.Resources = append(b.mr.Resources, resource)
	return b
}

// Build returns the MembersResult entity.
func (b *membersResultBuilder) Build() *MembersResult {
	mr := b.mr
	mr.Items = len(mr.Resources)
	mr.SetHashCode()
	return mr
}

// groupMembersBuilder is used to build a GroupMembers entity and ensure the calculated hash code is set.
type groupMembersBuilder struct {
	gm *GroupMembers
}

// GroupMembersBuilder creates a new groupMembersBuilder entity.
func GroupMembersBuilder() *groupMembersBuilder {
	return &groupMembersBuilder{
		gm: &GroupMembers{
			Resources: make([]*Member, 0),
		},
	}
}

// WithGroup sets the Group field of the GroupMembers entity.
func (b *groupMembersBuilder) WithGroup(group *Group) *groupMembersBuilder {
	b.gm.Group = group
	return b
}

// WithResources sets the Resources field of the GroupMembers entity.
func (b *groupMembersBuilder) WithResources(resources []*Member) *groupMembersBuilder {
	b.gm.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupMembers entity.
func (b *groupMembersBuilder) WithResource(resource *Member) *groupMembersBuilder {
	b.gm.Resources = append(b.gm.Resources, resource)
	return b
}

// Build returns the GroupMembers entity.
func (b *groupMembersBuilder) Build() *GroupMembers {
	gm := b.gm
	gm.Items = len(gm.Resources)
	gm.SetHashCode()
	return gm
}

// groupsMembersResultBuilder is used to build a GroupsMembersResult entity and ensure the calculated hash code and items is set.
type groupsMembersResultBuilder struct {
	gmr *GroupsMembersResult
}

// GroupsMembersResultBuilder creates a new groupsMembersResultBuilder entity.
func GroupsMembersResultBuilder() *groupsMembersResultBuilder {
	return &groupsMembersResultBuilder{
		gmr: &GroupsMembersResult{
			Resources: make([]*GroupMembers, 0),
		},
	}
}

// WithResources sets the Resources field of the GroupsMembersResult entity.
func (b *groupsMembersResultBuilder) WithResources(resources []*GroupMembers) *groupsMembersResultBuilder {
	b.gmr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupsMembersResult entity.
func (b *groupsMembersResultBuilder) WithResource(resource *GroupMembers) *groupsMembersResultBuilder {
	b.gmr.Resources = append(b.gmr.Resources, resource)
	return b
}

// Build returns the GroupsMembersResult entity.
func (b *groupsMembersResultBuilder) Build() *GroupsMembersResult {
	gmr := b.gmr
	gmr.Items = len(gmr.Resources)
	gmr.SetHashCode()
	return gmr
}

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
