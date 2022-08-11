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
