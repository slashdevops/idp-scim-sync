package model

// MemberBuilderChoice is used to build a Member entity and ensure the calculated hash code is set.
type MemberBuilderChoice struct {
	m *Member
}

// MemberBuilder creates a new MemberBuilderChoice entity.
func MemberBuilder() *MemberBuilderChoice {
	return &MemberBuilderChoice{
		m: &Member{},
	}
}

// WithIPID sets the IPID field of the Member entity.
func (b *MemberBuilderChoice) WithIPID(ipid string) *MemberBuilderChoice {
	b.m.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Member entity.
func (b *MemberBuilderChoice) WithSCIMID(scimid string) *MemberBuilderChoice {
	b.m.SCIMID = scimid
	return b
}

// WithEmail sets the Email field of the Member entity.
func (b *MemberBuilderChoice) WithEmail(email string) *MemberBuilderChoice {
	b.m.Email = email
	return b
}

// WithStatus sets the Status field of the Member entity.
func (b *MemberBuilderChoice) WithStatus(status string) *MemberBuilderChoice {
	b.m.Status = status
	return b
}

// Build returns the Member entity.
func (b *MemberBuilderChoice) Build() *Member {
	m := b.m
	m.SetHashCode()
	return m
}

// MembersResultBuilderChoice is used to build a MembersResult entity and ensure the calculated hash code and items is set.
type MembersResultBuilderChoice struct {
	mr *MembersResult
}

// MembersResultBuilder creates a new MembersResultBuilderChoice entity.
func MembersResultBuilder() *MembersResultBuilderChoice {
	return &MembersResultBuilderChoice{
		mr: &MembersResult{
			Resources: make([]*Member, 0),
		},
	}
}

// WithResources sets the Resources field of the MembersResult entity.
func (b *MembersResultBuilderChoice) WithResources(resources []*Member) *MembersResultBuilderChoice {
	b.mr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the MembersResult entity.
func (b *MembersResultBuilderChoice) WithResource(resource *Member) *MembersResultBuilderChoice {
	b.mr.Resources = append(b.mr.Resources, resource)
	return b
}

// Build returns the MembersResult entity.
func (b *MembersResultBuilderChoice) Build() *MembersResult {
	mr := b.mr
	mr.Items = len(mr.Resources)
	mr.SetHashCode()
	return mr
}

// GroupMembersBuilderChoice is used to build a GroupMembers entity and ensure the calculated hash code is set.
type GroupMembersBuilderChoice struct {
	gm *GroupMembers
}

// GroupMembersBuilder creates a new GroupMembersBuilderChoice entity.
func GroupMembersBuilder() *GroupMembersBuilderChoice {
	return &GroupMembersBuilderChoice{
		gm: &GroupMembers{
			Resources: make([]*Member, 0),
		},
	}
}

// WithGroup sets the Group field of the GroupMembers entity.
func (b *GroupMembersBuilderChoice) WithGroup(group *Group) *GroupMembersBuilderChoice {
	b.gm.Group = group
	return b
}

// WithResources sets the Resources field of the GroupMembers entity.
func (b *GroupMembersBuilderChoice) WithResources(resources []*Member) *GroupMembersBuilderChoice {
	b.gm.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupMembers entity.
func (b *GroupMembersBuilderChoice) WithResource(resource *Member) *GroupMembersBuilderChoice {
	b.gm.Resources = append(b.gm.Resources, resource)
	return b
}

// Build returns the GroupMembers entity.
func (b *GroupMembersBuilderChoice) Build() *GroupMembers {
	gm := b.gm
	gm.Items = len(gm.Resources)
	gm.SetHashCode()
	return gm
}

// GroupsMembersResultBuilderChoice is used to build a GroupsMembersResult entity and ensure the calculated hash code and items is set.
type GroupsMembersResultBuilderChoice struct {
	gmr *GroupsMembersResult
}

// GroupsMembersResultBuilder creates a new GroupsMembersResultBuilderChoice entity.
func GroupsMembersResultBuilder() *GroupsMembersResultBuilderChoice {
	return &GroupsMembersResultBuilderChoice{
		gmr: &GroupsMembersResult{
			Resources: make([]*GroupMembers, 0),
		},
	}
}

// WithResources sets the Resources field of the GroupsMembersResult entity.
func (b *GroupsMembersResultBuilderChoice) WithResources(resources []*GroupMembers) *GroupsMembersResultBuilderChoice {
	b.gmr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupsMembersResult entity.
func (b *GroupsMembersResultBuilderChoice) WithResource(resource *GroupMembers) *GroupsMembersResultBuilderChoice {
	b.gmr.Resources = append(b.gmr.Resources, resource)
	return b
}

// Build returns the GroupsMembersResult entity.
func (b *GroupsMembersResultBuilderChoice) Build() *GroupsMembersResult {
	gmr := b.gmr
	gmr.Items = len(gmr.Resources)
	gmr.SetHashCode()
	return gmr
}
