package model

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
