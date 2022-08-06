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
