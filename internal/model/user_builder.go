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

// WithUserName sets the UserName field of the User entity.
func (b *UserBuilderChoice) WithUserName(userName string) *UserBuilderChoice {
	b.u.UserName = userName
	return b
}

// WithFormattedName sets the Formatted field of the User entity.
func (b *UserBuilderChoice) WithFormattedName(formatted string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.Formatted = formatted
	return b
}

// WithFamilyName sets the FamilyName field of the User entity.
func (b *UserBuilderChoice) WithFamilyName(familyName string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.FamilyName = familyName
	return b
}

// WithGivenName sets the GivenName field of the User entity.
func (b *UserBuilderChoice) WithGivenName(givenName string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.GivenName = givenName
	return b
}

// WithMiddleName sets the MiddleName field of the User entity.
func (b *UserBuilderChoice) WithMiddleName(middleName string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.MiddleName = middleName
	return b
}

// WithHonorificPrefixName sets the HonorificPrefix field of the User entity.
func (b *UserBuilderChoice) WithHonorificPrefixName(honorificPrefix string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.HonorificPrefix = honorificPrefix
	return b
}

// WithHonorificSuffixName sets the HonorificSuffix field of the User entity.
func (b *UserBuilderChoice) WithHonorificSuffixName(honorificSuffix string) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name.HonorificSuffix = honorificSuffix
	return b
}

// WithDisplayName sets the DisplayName field of the User entity.
func (b *UserBuilderChoice) WithDisplayName(displayName string) *UserBuilderChoice {
	b.u.DisplayName = displayName
	return b
}

// WithNickName sets the NickName field of the User entity.
func (b *UserBuilderChoice) WithNickName(nickName string) *UserBuilderChoice {
	b.u.NickName = nickName
	return b
}

// WithProfileURL sets the ProfileURL field of the User entity.
func (b *UserBuilderChoice) WithProfileURL(profileURL string) *UserBuilderChoice {
	b.u.ProfileURL = profileURL
	return b
}

// WithTitle sets the Title field of the User entity.
func (b *UserBuilderChoice) WithTitle(title string) *UserBuilderChoice {
	b.u.Title = title
	return b
}

// WithUserType sets the UserType field of the User entity.
func (b *UserBuilderChoice) WithUserType(userType string) *UserBuilderChoice {
	b.u.UserType = userType
	return b
}

// WithPreferredLanguage sets the PreferredLanguage field of the User entity.
func (b *UserBuilderChoice) WithPreferredLanguage(preferredLanguage string) *UserBuilderChoice {
	b.u.PreferredLanguage = preferredLanguage
	return b
}

// WithLocale sets the Locale field of the User entity.
func (b *UserBuilderChoice) WithLocale(locale string) *UserBuilderChoice {
	b.u.Locale = locale
	return b
}

// WithTimezone sets the Timezone field of the User entity.
func (b *UserBuilderChoice) WithTimezone(timezone string) *UserBuilderChoice {
	b.u.Timezone = timezone
	return b
}

// WithEmail sets the Email field of the User entity.
// if the Emails contains an email, it will be replaced by the new one.
func (b *UserBuilderChoice) WithEmail(email Email) *UserBuilderChoice {
	if len(b.u.Emails) == 0 {
		b.u.Emails = append(b.u.Emails, email)
	} else {
		b.u.Emails[0] = email
	}
	return b
}

// WithEmails sets the Emails field of the User entity.
func (b *UserBuilderChoice) WithEmails(emails []Email) *UserBuilderChoice {
	b.u.Emails = emails
	return b
}

// WithAddress sets the Address field of the User entity.
// if the Addresses contains an address, it will be replaced by the new one.
func (b *UserBuilderChoice) WithAddress(address Address) *UserBuilderChoice {
	if len(b.u.Addresses) == 0 {
		b.u.Addresses = append(b.u.Addresses, address)
	} else {
		b.u.Addresses[0] = address
	}
	return b
}

// WithPhoneNumber sets the PhoneNumber field of the User entity.
// if the PhoneNumbers contains a phone number, it will be replaced by the new one.
func (b *UserBuilderChoice) WithPhoneNumber(phoneNumber PhoneNumber) *UserBuilderChoice {
	if len(b.u.PhoneNumbers) == 0 {
		b.u.PhoneNumbers = append(b.u.PhoneNumbers, phoneNumber)
	} else {
		b.u.PhoneNumbers[0] = phoneNumber
	}
	return b
}

// WithPhoneNumbers sets the PhoneNumbers field of the User entity.
func (b *UserBuilderChoice) WithPhoneNumbers(phoneNumbers []PhoneNumber) *UserBuilderChoice {
	b.u.PhoneNumbers = phoneNumbers
	return b
}

// WithName sets the Name field of the User entity.
func (b *UserBuilderChoice) WithName(name Name) *UserBuilderChoice {
	if b.u.Name == nil {
		b.u.Name = &Name{}
	}
	b.u.Name = &name
	return b
}

// WithEnterpriseData sets the EnterpriseData field of the User entity.
func (b *UserBuilderChoice) WithEnterpriseData(enterpriseData EnterpriseData) *UserBuilderChoice {
	b.u.EnterpriseData = &enterpriseData
	return b
}

// WithActive sets the Active field of the User entity.
func (b *UserBuilderChoice) WithActive(active bool) *UserBuilderChoice {
	b.u.Active = active
	return b
}

// Build returns the User entity.
func (b *UserBuilderChoice) Build() *User {
	b.u.SetHashCode()
	return b.u
}

// NameBuilderChoice is the builder of Name entity.
type NameBuilderChoice struct {
	n *Name
}

// NameBuilder creates a new NameBuilderChoice entity.
func NameBuilder() *NameBuilderChoice {
	return &NameBuilderChoice{
		n: &Name{},
	}
}

// WithFormatted sets the Formatted field of the Name entity.
func (b *NameBuilderChoice) WithFormatted(formatted string) *NameBuilderChoice {
	b.n.Formatted = formatted
	return b
}

// WithFamilyName sets the FamilyName field of the Name entity.
func (b *NameBuilderChoice) WithFamilyName(familyName string) *NameBuilderChoice {
	b.n.FamilyName = familyName
	return b
}

// WithGivenName sets the GivenName field of the Name entity.
func (b *NameBuilderChoice) WithGivenName(givenName string) *NameBuilderChoice {
	b.n.GivenName = givenName
	return b
}

// WithMiddleName sets the MiddleName field of the Name entity.
func (b *NameBuilderChoice) WithMiddleName(middleName string) *NameBuilderChoice {
	b.n.MiddleName = middleName
	return b
}

// WithHonorificPrefix sets the HonorificPrefix field of the Name entity.
func (b *NameBuilderChoice) WithHonorificPrefix(honorificPrefix string) *NameBuilderChoice {
	b.n.HonorificPrefix = honorificPrefix
	return b
}

// WithHonorificSuffix sets the HonorificSuffix field of the Name entity.
func (b *NameBuilderChoice) WithHonorificSuffix(honorificSuffix string) *NameBuilderChoice {
	b.n.HonorificSuffix = honorificSuffix
	return b
}

// Build returns the Name entity.
func (b *NameBuilderChoice) Build() *Name {
	return b.n
}

// EnterpriseDataBuilderChoice is the builder of EnterpriseData entity.
type EnterpriseDataBuilderChoice struct {
	ed *EnterpriseData
}

// EnterpriseDataBuilder creates a new EnterpriseDataBuilderChoice entity.
func EnterpriseDataBuilder() *EnterpriseDataBuilderChoice {
	return &EnterpriseDataBuilderChoice{
		ed: &EnterpriseData{},
	}
}

// WithEmployeeNumber sets the EmployeeNumber field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithEmployeeNumber(employeeNumber string) *EnterpriseDataBuilderChoice {
	b.ed.EmployeeNumber = employeeNumber
	return b
}

// WithCostCenter sets the CostCenter field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithCostCenter(costCenter string) *EnterpriseDataBuilderChoice {
	b.ed.CostCenter = costCenter
	return b
}

// WithOrganization sets the Organization field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithOrganization(organization string) *EnterpriseDataBuilderChoice {
	b.ed.Organization = organization
	return b
}

// WithDivision sets the Division field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithDivision(division string) *EnterpriseDataBuilderChoice {
	b.ed.Division = division
	return b
}

// WithDepartment sets the Department field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithDepartment(department string) *EnterpriseDataBuilderChoice {
	b.ed.Department = department
	return b
}

// WithManager sets the Manager field of the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) WithManager(manager Manager) *EnterpriseDataBuilderChoice {
	b.ed.Manager = manager
	return b
}

// Build returns the EnterpriseData entity.
func (b *EnterpriseDataBuilderChoice) Build() *EnterpriseData {
	return b.ed
}

// ManagerBuilderChoice is the builder of Manager entity.
type ManagerBuilderChoice struct {
	m *Manager
}

// ManagerBuilder creates a new ManagerBuilderChoice entity.
func ManagerBuilder() *ManagerBuilderChoice {
	return &ManagerBuilderChoice{
		m: &Manager{},
	}
}

// WithValue sets the Value field of the Manager entity.
func (b *ManagerBuilderChoice) WithValue(value string) *ManagerBuilderChoice {
	b.m.Value = value
	return b
}

// WithRef sets the Ref field of the Manager entity.
func (b *ManagerBuilderChoice) WithRef(ref string) *ManagerBuilderChoice {
	b.m.Ref = ref
	return b
}

// Build returns the Manager entity.
func (b *ManagerBuilderChoice) Build() *Manager {
	return b.m
}

// EmailBuilderChoice is the builder of Email entity.
type EmailBuilderChoice struct {
	e Email
}

// EmailBuilder creates a new EmailBuilderChoice entity.
func EmailBuilder() *EmailBuilderChoice {
	return &EmailBuilderChoice{
		e: Email{},
	}
}

// WithValue sets the Value field of the Email entity.
func (b *EmailBuilderChoice) WithValue(value string) *EmailBuilderChoice {
	b.e.Value = value
	return b
}

// WithType sets the Type field of the Email entity.
func (b *EmailBuilderChoice) WithType(emailType string) *EmailBuilderChoice {
	b.e.Type = emailType
	return b
}

// WithPrimary sets the Primary field of the Email entity.
func (b *EmailBuilderChoice) WithPrimary(primary bool) *EmailBuilderChoice {
	b.e.Primary = primary
	return b
}

// Build returns the Email entity.
func (b *EmailBuilderChoice) Build() Email {
	return b.e
}

// AddressBuilderChoice is the builder of Address entity.
type AddressBuilderChoice struct {
	a Address
}

// AddressBuilder creates a new AddressBuilderChoice entity.
func AddressBuilder() *AddressBuilderChoice {
	return &AddressBuilderChoice{
		a: Address{},
	}
}

// WithType sets the Type field of the Address entity.
func (b *AddressBuilderChoice) WithType(addressType string) *AddressBuilderChoice {
	b.a.Type = addressType
	return b
}

// WithFormatted sets the Formatted field of the Address entity.
func (b *AddressBuilderChoice) WithFormatted(formatted string) *AddressBuilderChoice {
	b.a.Formatted = formatted
	return b
}

// WithStreetAddress sets the StreetAddress field of the Address entity.
func (b *AddressBuilderChoice) WithStreetAddress(streetAddress string) *AddressBuilderChoice {
	b.a.StreetAddress = streetAddress
	return b
}

// WithLocality sets the Locality field of the Address entity.
func (b *AddressBuilderChoice) WithLocality(locality string) *AddressBuilderChoice {
	b.a.Locality = locality
	return b
}

// WithRegion sets the Region field of the Address entity.
func (b *AddressBuilderChoice) WithRegion(region string) *AddressBuilderChoice {
	b.a.Region = region
	return b
}

// WithPostalCode sets the PostalCode field of the Address entity.
func (b *AddressBuilderChoice) WithPostalCode(postalCode string) *AddressBuilderChoice {
	b.a.PostalCode = postalCode
	return b
}

// WithCountry sets the Country field of the Address entity.
func (b *AddressBuilderChoice) WithCountry(country string) *AddressBuilderChoice {
	b.a.Country = country
	return b
}

// WithPrimary sets the Primary field of the Address entity.
func (b *AddressBuilderChoice) WithPrimary(primary bool) *AddressBuilderChoice {
	b.a.Primary = primary
	return b
}

// Build returns the Address entity.
func (b *AddressBuilderChoice) Build() Address {
	return b.a
}

// PhoneNumberBuilderChoice is the builder of PhoneNumber entity.
type PhoneNumberBuilderChoice struct {
	pn PhoneNumber
}

// PhoneNumberBuilder creates a new PhoneNumberBuilderChoice entity.
func PhoneNumberBuilder() *PhoneNumberBuilderChoice {
	return &PhoneNumberBuilderChoice{
		pn: PhoneNumber{},
	}
}

// WithValue sets the Value field of the PhoneNumber entity.
func (b *PhoneNumberBuilderChoice) WithValue(value string) *PhoneNumberBuilderChoice {
	b.pn.Value = value
	return b
}

// WithType sets the Type field of the PhoneNumber entity.
func (b *PhoneNumberBuilderChoice) WithType(phoneNumberType string) *PhoneNumberBuilderChoice {
	b.pn.Type = phoneNumberType
	return b
}

// Build returns the PhoneNumber entity.
func (b *PhoneNumberBuilderChoice) Build() PhoneNumber {
	return b.pn
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
