package model

// UserBuilderChoice is the builder of User entity.
type UserBuilderChoice struct {
	u *User
}

// UserBuilder creates a new UserBuilderChoice entity.
func UserBuilder() *UserBuilderChoice {
	return &UserBuilderChoice{
		u: &User{},
	}
}

// WithIPID sets the IPID field of the User entity.
func (b *UserBuilderChoice) WithIPID(ipid string) *UserBuilderChoice {
	b.u.IPID = ipid
	return b
}

// WithSCIMID sets the SCIMID field of the User entity.
func (b *UserBuilderChoice) WithSCIMID(scimid string) *UserBuilderChoice {
	b.u.SCIMID = scimid
	return b
}

// WithFamilyName sets the FamilyName field of the User entity.
func (b *UserBuilderChoice) WithFamilyName(familyName string) *UserBuilderChoice {
	b.u.Name.FamilyName = familyName
	return b
}

// WithGivenName sets the GivenName field of the User entity.
func (b *UserBuilderChoice) WithGivenName(givenName string) *UserBuilderChoice {
	b.u.Name.GivenName = givenName
	return b
}

// WithDisplayName sets the DisplayName field of the User entity.
func (b *UserBuilderChoice) WithDisplayName(displayName string) *UserBuilderChoice {
	b.u.DisplayName = displayName
	return b
}

// WithActive sets the Active field of the User entity.
func (b *UserBuilderChoice) WithActive(active bool) *UserBuilderChoice {
	b.u.Active = active
	return b
}

// WithEmail sets the Email field of the User entity.
func (b *UserBuilderChoice) WithEmail(email string) *UserBuilderChoice {
	b.u.Email = email
	return b
}

// Build returns the User entity.
func (b *UserBuilderChoice) Build() *User {
	b.u.SetHashCode()
	return b.u
}

// UsersResultBuilderChoice is used to build a UsersResult entity and ensure the calculated hash code and items.
type UsersResultBuilderChoice struct {
	ur *UsersResult
}

// UsersResultBuilder creates a new UsersResultBuilderChoice entity.
func UsersResultBuilder() *UsersResultBuilderChoice {
	return &UsersResultBuilderChoice{
		ur: &UsersResult{
			Resources: make([]*User, 0),
		},
	}
}

// WithResources sets the Resources field of the UsersResult entity.
func (b *UsersResultBuilderChoice) WithResources(resources []*User) *UsersResultBuilderChoice {
	b.ur.Resources = resources
	return b
}

// WithResource add the resource to a Resources field of the UsersResult entity.
func (b *UsersResultBuilderChoice) WithResource(resource *User) *UsersResultBuilderChoice {
	b.ur.Resources = append(b.ur.Resources, resource)
	return b
}

// Build returns the UserResult entity.
func (b *UsersResultBuilderChoice) Build() *UsersResult {
	b.ur.Items = len(b.ur.Resources)
	b.ur.SetHashCode()
	return b.ur
}
