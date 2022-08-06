package model

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
