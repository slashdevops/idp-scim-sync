package aws

import (
	"encoding/json"
	"log"
)

// Name represent a name entity
type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName,omitempty"`
}

// Email represent an email entity
type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary"`
}

// Addresses represent an address entity
type Addresses struct {
	Type          string `json:"type"`
	StreetAddress string `json:"streetAddress"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
	Formatted     string `json:"formatted"`
	Primary       bool   `json:"primary"`
}

// Meta represent a meta entity
type Meta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
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
	ID          string       `json:"id"`
	ExternalID  string       `json:"externalId,omitempty"`
	Meta        Meta         `json:"meta,omitempty"`
	Schemas     []string     `json:"schemas,omitempty"`
	UserName    string       `json:"userName"`
	Name        Name         `json:"name,omitempty"`
	DisplayName string       `json:"displayName,omitempty"`
	NickName    string       `json:"nickName,omitempty"`
	ProfileURL  string       `json:"profileURL,omitempty"`
	Active      bool         `json:"active,omitempty"`
	Emails      []*Email     `json:"emails,omitempty"`
	Addresses   []*Addresses `json:"addresses,omitempty"`
}

// String is the implementation of Stringer interface
func (u *User) String() string {
	JSON, err := json.Marshal(u)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(JSON)
}

// GetUserResponse represent a get user response entity
type GetUserResponse User

// CreateUserRequest represent a create user request entity
type CreateUserRequest User

// PutUserRequest represent a put user request entity
type PutUserRequest User

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
