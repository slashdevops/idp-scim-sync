package model

// UserBuilder is the builder of User entity.
type UserBuilder struct {
	u *User
}

// NewUserBuilder creates a new UserBuilder entity.
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		u: &User{},
	}
}

// WithIPID sets the IPID field of the User entity.
func (b *UserBuilder) WithIPID(ipid string) *UserBuilder {
	b.u.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the User entity.
func (b *UserBuilder) WithSCIMID(scimid string) *UserBuilder {
	b.u.SCIMID = scimid
	return b
}

// WithFamilyName sets the FamilyName field of the User entity.
func (b *UserBuilder) WithFamilyName(familyName string) *UserBuilder {
	b.u.Name.FamilyName = familyName
	return b
}

// WithGivenName sets the GivenName field of the User entity.
func (b *UserBuilder) WithGivenName(givenName string) *UserBuilder {
	b.u.Name.GivenName = givenName
	return b
}

// WithDisplayName sets the DisplayName field of the User entity.
func (b *UserBuilder) WithDisplayName(displayName string) *UserBuilder {
	b.u.DisplayName = displayName
	return b
}

// WithActive sets the Active field of the User entity.
func (b *UserBuilder) WithActive(active bool) *UserBuilder {
	b.u.Active = active
	return b
}

// WithEmail sets the Email field of the User entity.
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.u.Email = email
	return b
}

// Build returns the User entity.
func (b *UserBuilder) Build() *User {
	u := b.u
	u.SetHashCode()
	return u
}

// UsersResultBuilder is used to build a UsersResult entity and ensure the calculated hash code and items.
type UsersResultBuilder struct {
	ur *UsersResult
}

// NewUsersResultBuilder creates a new UsersResultBuilder entity.
func NewUsersResultBuilder() *UsersResultBuilder {
	return &UsersResultBuilder{
		ur: &UsersResult{
			Resources: make([]*User, 0),
		},
	}
}

// WithResources sets the Resources field of the UsersResult entity.
func (b *UsersResultBuilder) WithResources(resources []*User) *UsersResultBuilder {
	b.ur.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the UsersResult entity.
func (b *UsersResultBuilder) WithResource(resource *User) *UsersResultBuilder {
	b.ur.Resources = append(b.ur.Resources, resource)
	return b
}

// Build returns the UserResult entity.
func (b *UsersResultBuilder) Build() *UsersResult {
	ur := b.ur
	ur.Items = len(ur.Resources)
	ur.SetHashCode()
	return ur
}

// GroupBuilder is used to build a Group entity and ensure the calculated hash code is set.
type GroupBuilder struct {
	g *Group
}

// NewGroupBuilder creates a new GroupBuilder entity.
func NewGroupBuilder() *GroupBuilder {
	return &GroupBuilder{
		g: &Group{},
	}
}

// WithIPID sets the IPID field of the Group entity.
func (b *GroupBuilder) WithIPID(ipid string) *GroupBuilder {
	b.g.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Group entity.
func (b *GroupBuilder) WithSCIMID(scimid string) *GroupBuilder {
	b.g.SCIMID = scimid
	return b
}

// WithName sets the Name field of the Group entity.
func (b *GroupBuilder) WithName(name string) *GroupBuilder {
	b.g.Name = name
	return b
}

// WithEmail sets the Email field of the Group entity.
func (b *GroupBuilder) WithEmail(email string) *GroupBuilder {
	b.g.Email = email
	return b
}

// Build returns the Group entity.
func (b *GroupBuilder) Build() *Group {
	g := b.g
	g.SetHashCode()
	return g
}

// GroupsResultBuilder is used to build a GroupsResult entity and ensure the calculated hash code and items is set.
type GroupsResultBuilder struct {
	gr *GroupsResult
}

// NewGroupsResultBuilder creates a new GroupsResultBuilder entity.
func NewGroupsResultBuilder() *GroupsResultBuilder {
	return &GroupsResultBuilder{
		gr: &GroupsResult{
			Resources: make([]*Group, 0),
		},
	}
}

// WithResources sets the Resources field of the GroupsResult entity.
func (b *GroupsResultBuilder) WithResources(resources []*Group) *GroupsResultBuilder {
	b.gr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupsResult entity.
func (b *GroupsResultBuilder) WithResource(resource *Group) *GroupsResultBuilder {
	b.gr.Resources = append(b.gr.Resources, resource)
	return b
}

// Build returns the GroupsResult entity.
func (b *GroupsResultBuilder) Build() *GroupsResult {
	gr := b.gr
	gr.Items = len(gr.Resources)
	gr.SetHashCode()
	return gr
}

// MemberBuilder is used to build a Member entity and ensure the calculated hash code is set.
type MemberBuilder struct {
	m *Member
}

// NewMemberBuilder creates a new MemberBuilder entity.
func NewMemberBuilder() *MemberBuilder {
	return &MemberBuilder{
		m: &Member{},
	}
}

// WithIPID sets the IPID field of the Member entity.
func (b *MemberBuilder) WithIPID(ipid string) *MemberBuilder {
	b.m.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Member entity.
func (b *MemberBuilder) WithSCIMID(scimid string) *MemberBuilder {
	b.m.SCIMID = scimid
	return b
}

// WithEmail sets the Email field of the Member entity.
func (b *MemberBuilder) WithEmail(email string) *MemberBuilder {
	b.m.Email = email
	return b
}

// WithStatus sets the Status field of the Member entity.
func (b *MemberBuilder) WithStatus(status string) *MemberBuilder {
	b.m.Status = status
	return b
}

// Build returns the Member entity.
func (b *MemberBuilder) Build() *Member {
	m := b.m
	m.SetHashCode()
	return m
}

// MembersResultBuilder is used to build a MembersResult entity and ensure the calculated hash code and items is set.
type MembersResultBuilder struct {
	mr *MembersResult
}

// NewMembersResultBuilder creates a new MembersResultBuilder entity.
func NewMembersResultBuilder() *MembersResultBuilder {
	return &MembersResultBuilder{
		mr: &MembersResult{
			Resources: make([]*Member, 0),
		},
	}
}

// WithResources sets the Resources field of the MembersResult entity.
func (b *MembersResultBuilder) WithResources(resources []*Member) *MembersResultBuilder {
	b.mr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the MembersResult entity.
func (b *MembersResultBuilder) WithResource(resource *Member) *MembersResultBuilder {
	b.mr.Resources = append(b.mr.Resources, resource)
	return b
}

// Build returns the MembersResult entity.
func (b *MembersResultBuilder) Build() *MembersResult {
	mr := b.mr
	mr.Items = len(mr.Resources)
	mr.SetHashCode()
	return mr
}
