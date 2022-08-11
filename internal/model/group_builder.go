package model

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
