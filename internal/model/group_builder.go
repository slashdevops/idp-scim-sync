package model

// GroupBuilderChoice is used to build a Group entity and ensure the calculated hash code is set.
type GroupBuilderChoice struct {
	g *Group
}

// GroupBuilder creates a new GroupBuilderChoice entity.
func GroupBuilder() *GroupBuilderChoice {
	return &GroupBuilderChoice{
		g: &Group{},
	}
}

// WithIPID sets the IPID field of the Group entity.
func (b *GroupBuilderChoice) WithIPID(ipid string) *GroupBuilderChoice {
	b.g.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the Group entity.
func (b *GroupBuilderChoice) WithSCIMID(scimid string) *GroupBuilderChoice {
	b.g.SCIMID = scimid
	return b
}

// WithName sets the Name field of the Group entity.
func (b *GroupBuilderChoice) WithName(name string) *GroupBuilderChoice {
	b.g.Name = name
	return b
}

// WithEmail sets the Email field of the Group entity.
func (b *GroupBuilderChoice) WithEmail(email string) *GroupBuilderChoice {
	b.g.Email = email
	return b
}

// Build returns the Group entity.
func (b *GroupBuilderChoice) Build() *Group {
	g := b.g
	g.SetHashCode()
	return g
}

// GroupsResultBuilderChoice is used to build a GroupsResult entity and ensure the calculated hash code and items is set.
type GroupsResultBuilderChoice struct {
	gr *GroupsResult
}

// GroupsResultBuilder creates a new GroupsResultBuilderChoice entity.
func GroupsResultBuilder() *GroupsResultBuilderChoice {
	return &GroupsResultBuilderChoice{
		gr: &GroupsResult{
			Resources: make([]*Group, 0),
		},
	}
}

// WithResources sets the Resources field of the GroupsResult entity.
func (b *GroupsResultBuilderChoice) WithResources(resources []*Group) *GroupsResultBuilderChoice {
	b.gr.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the GroupsResult entity.
func (b *GroupsResultBuilderChoice) WithResource(resource *Group) *GroupsResultBuilderChoice {
	b.gr.Resources = append(b.gr.Resources, resource)
	return b
}

// Build returns the GroupsResult entity.
func (b *GroupsResultBuilderChoice) Build() *GroupsResult {
	gr := b.gr
	gr.Items = len(gr.Resources)
	gr.SetHashCode()
	return gr
}
