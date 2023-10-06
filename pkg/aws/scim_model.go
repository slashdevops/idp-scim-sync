package aws

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

var (
	// ErrUserIDEmpty is returned when the user id is empty.
	ErrUserIDEmpty = errors.Errorf("aws: user id may not be empty")

	// ErrEmailsTooMany is returned when the emails has more than one entity.
	ErrEmailsTooMany = errors.Errorf("aws: emails may not be more than 1")

	// ErrEmailsEmpty
	ErrEmailsEmpty = errors.Errorf("aws: emails may not be empty")

	// ErrFamilyNameEmpty is returned when the family name is empty.
	ErrFamilyNameEmpty = errors.Errorf("aws: family name may not be empty")

	// ErrDisplayNameEmpty is returned when the display name is empty.
	ErrDisplayNameEmpty = errors.Errorf("aws: display name may not be empty")

	// ErrGivenNameEmpty is returned when the given name is empty.
	ErrGivenNameEmpty = errors.Errorf("aws: given name may not be empty")

	// ErrUserNameEmpty is returned when the user name is empty.
	ErrUserNameEmpty = errors.Errorf("aws: user name may not be empty")

	// ErrUserUserNameEmpty is returned when the userName is empty.
	ErrUserUserNameEmpty = errors.Errorf("aws: userName may not be empty")

	// ErrPrimaryEmailEmpty is returned when the primary email is empty.
	ErrPrimaryEmailEmpty = errors.Errorf("aws: primary email may not be empty")

	// ErrAddressesTooMany is returned when the addresses has more than one entity.
	ErrAddressesTooMany = errors.Errorf("aws: addresses may not be more than 1")

	// ErrPhoneNumbersTooMany is returned when the phone numbers has more than one entity.
	ErrPhoneNumbersTooMany = errors.Errorf("aws: phone numbers may not be more than 1")

	// ErrTooManyPrimaryEmails when there are more than one primary email
	ErrTooManyPrimaryEmails = errors.Errorf("aws: there can only be one primary email")
)

// Name represent a name entity
type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// Email represent an email entity
type Email struct {
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// Addresses represent an address entity
type Address struct {
	Type          string `json:"type,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Country       string `json:"country,omitempty"`
}

type PhoneNumber struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Manager struct {
	Value string `json:"value,omitempty"`
	Ref   string `json:"$ref,omitempty"`
}

type SchemaEnterpriseUser struct {
	EmployeeNumber string   `json:"employeeNumber,omitempty"`
	CostCenter     string   `json:"costCenter,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Division       string   `json:"division,omitempty"`
	Department     string   `json:"department,omitempty"`
	Manager        *Manager `json:"manager,omitempty"`
}

// Meta represent a meta entity
type Meta struct {
	ResourceType string `json:"resourceType,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

// Operation represent an operation entity
type Operation struct {
	OP    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// Patch represent a patch entity and its operations
type Patch struct {
	Schemas    []string     `json:"schemas"`
	Operations []*Operation `json:"Operations"`
}

// ListResponse represent a general response entity
type ListResponse struct {
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Schemas      []string `json:"schemas"`
}

// User represent a user entity
type User struct {
	ID                   string                `json:"id,omitempty"`
	ExternalID           string                `json:"externalId,omitempty"`
	UserName             string                `json:"userName,omitempty"`
	DisplayName          string                `json:"displayName,omitempty"`
	NickName             string                `json:"nickName,omitempty"`
	ProfileURL           string                `json:"profileURL,omitempty"`
	UserType             string                `json:"userType,omitempty"`
	Title                string                `json:"title,omitempty"`
	PreferredLanguage    string                `json:"preferredLanguage,omitempty"`
	Locale               string                `json:"locale,omitempty"`
	Timezone             string                `json:"timezone,omitempty"`
	Name                 *Name                 `json:"name,omitempty"`
	Meta                 *Meta                 `json:"meta,omitempty"`
	SchemaEnterpriseUser *SchemaEnterpriseUser `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	Schemas              []string              `json:"schemas,omitempty"`
	Addresses            []Address             `json:"addresses,omitempty"`
	Emails               []Email               `json:"emails,omitempty"`
	PhoneNumbers         []PhoneNumber         `json:"phoneNumbers,omitempty"`
	Active               bool                  `json:"active,omitempty"`
}

// Validate check if the user entity is valid according to the SCIM spec constraints
// Reference: https://docs.aws.amazon.com/singlesignon/latest/developerguide/createuser.html
func (u *User) Validate() error {
	if u.UserName == "" {
		return ErrUserNameEmpty
	}
	if u.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	if u.Name.GivenName == "" {
		return ErrGivenNameEmpty
	}
	if u.Name.FamilyName == "" {
		return ErrFamilyNameEmpty
	}
	if len(u.Emails) == 0 {
		return ErrEmailsEmpty
	}
	if len(u.Emails) > 1 {
		return ErrEmailsTooMany
	}

	primaryCount := 0
	for _, email := range u.Emails {
		if !email.Primary {
			return ErrPrimaryEmailEmpty
		}
		primaryCount++
	}
	if primaryCount != 1 {
		return ErrTooManyPrimaryEmails
	}

	if len(u.Addresses) > 1 {
		return ErrAddressesTooMany
	}
	if len(u.PhoneNumbers) > 1 {
		return ErrPhoneNumbersTooMany
	}

	return nil
}

// String is the implementation of Stringer interface
func (u *User) String() string {
	JSON, err := json.Marshal(u)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(JSON)
}

// GetPrimaryEmail returns the primary email of the user
func (u *User) GetPrimaryEmail() *Email {
	for _, email := range u.Emails {
		if email.Primary {
			return &email
		}
	}
	return nil
}

// GetPrimaryEmailAddress returns the primary email address of the user
func (u *User) GetPrimaryEmailAddress() string {
	for _, email := range u.Emails {
		if email.Primary {
			return email.Value
		}
	}
	return ""
}

// GetPrimaryAddress returns the primary address of the user
func (u *User) GetPrimaryAddress() *Address {
	if len(u.Addresses) > 0 {
		return &u.Addresses[0]
	}

	return nil
}

// GetUserResponse represent a get user response entity
type GetUserResponse User

// CreateUserRequest represent a create user request entity
type CreateUserRequest User

func (u *CreateUserRequest) Validate() error {
	return (*User)(u).Validate()
}

// PutUserRequest represent a put user request entity
type PutUserRequest User

func (u *PutUserRequest) Validate() error {
	return (*User)(u).Validate()
}

// CreateUserResponse represent a create user response entity
type CreateUserResponse User

// PutUserResponse represent a put user response entity
type PutUserResponse User

// PatchUserResponse represent a put user response entity
type PatchUserResponse User

// PatchUserRequest represent a patch user request entity
type PatchUserRequest struct {
	User  User  `json:"user"`
	Patch Patch `json:"patch"`
}

func (u *PatchUserRequest) Validate() error {
	if u.User.ID == "" {
		return ErrUserIDEmpty
	}
	return nil
}

// ListUsersResponse represent a list users response entity
type ListUsersResponse struct {
	ListResponse
	Resources []*User `json:"Resources"`
}

// Member represent a member group entity
type Member struct {
	Value string `json:"value"`
	Ref   string `json:"$ref"`
	Type  string `json:"type"`
}

// Group represent a group entity
type Group struct {
	ID          string    `json:"id"`
	Meta        Meta      `json:"meta,omitempty"`
	Schemas     []string  `json:"schemas,omitempty"`
	DisplayName string    `json:"displayName"`
	ExternalID  string    `json:"externalId,omitempty"`
	Members     []*Member `json:"members,omitempty"`
}

// Validate check if the group entity is valid according to the SCIM spec constraints
// Reference: https://docs.aws.amazon.com/singlesignon/latest/developerguide/creategroup.html
func (g *Group) Validate() error {
	if g.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}

// String is the implementation of Stringer interface
func (g *Group) String() string {
	JSON, err := json.Marshal(g)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(JSON)
}

// GetGroupResponse represent a group user response entity
type GetGroupResponse Group

// CreateGroupRequest represent a create group request entity
type CreateGroupRequest Group

func (g *CreateGroupRequest) Validate() error {
	return (*Group)(g).Validate()
}

// CreateGroupResponse represent a create group response entity
type CreateGroupResponse Group

// ListGroupsResponse represent a list groups response entity
type ListGroupsResponse struct {
	ListResponse
	Resources []*Group `json:"Resources"`
}

// PatchGroupRequest represent a patch group request entity
type PatchGroupRequest struct {
	Group Group `json:"group"`
	Patch Patch `json:"patch"`
}

// ServiceProviderConfig represent a service provider config entity
type ServiceProviderConfig struct {
	Schemas               []string `json:"schemas"`
	DocumentationURI      string   `json:"documentationUri"`
	AuthenticationSchemes []struct {
		Type             string `json:"type"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		SpecURI          string `json:"specURI"`
		DocumentationURI string `json:"documentationUri"`
		Primary          bool   `json:"primary"`
	} `json:"authenticationSchemes"`
	Patch struct {
		Supported bool `json:"supported"`
	} `json:"patch"`
	Bulk struct {
		Supported      bool `json:"supported"`
		MaxOperations  int  `json:"maxOperations"`
		MaxPayloadSize int  `json:"maxPayloadSize"`
	} `json:"bulk"`
	Filter struct {
		Supported  bool `json:"supported"`
		MaxResults int  `json:"maxResults"`
	} `json:"filter"`
	ChangePassword struct {
		Supported bool `json:"supported"`
	} `json:"changePassword"`
	Sort struct {
		Supported bool `json:"supported"`
	} `json:"sort"`
	Etag struct {
		Supported bool `json:"supported"`
	} `json:"etag"`
}
