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

// Meta represent a meta entity
type Meta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
}

// Operation represent an operation entity
type Operation struct {
	OP    string   `json:"op"`
	Path  string   `json:"path"`
	Value []string `json:"value"`
}

// OperationGroup represent an operation group entity
type OperationGroup struct {
	OP    string      `json:"op"`
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value"`
}

// Patch represent a patch entity
type Patch struct {
	Schemas    []string    `json:"schemas"`
	Operations []Operation `json:"Operations"`
}

// PatchGroup represent a path group entity
type PatchGroup struct {
	Schemas    []string         `json:"schemas"`
	Operations []OperationGroup `json:"Operations"`
}

// Member represent a member entity
type Member struct {
	Value string `json:"value"`
	Ref   string `json:"$ref"`
	Type  string `json:"type"`
}

// Group represent a group entity
type Group struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
	ExternalID  string   `json:"externalId"`
	Members     []Member `json:"members"`
}

func (g *Group) String() string {
	JSON, err := json.Marshal(g)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(JSON)
}

// User represent a user entity
type User struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Active      bool     `json:"active"`
	Emails      []Email  `json:"emails"`
}

func (u *User) String() string {
	JSON, err := json.Marshal(u)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(JSON)
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

// GeneralResponse represent a general response entity
type GeneralResponse struct {
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Schemas      []string `json:"schemas"`
}

// GetGroupResponse represent a group user response entity
type GetGroupResponse struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

// ListGroupsResponse represent a list groups response entity
type ListGroupsResponse struct {
	GeneralResponse
	Resources []*Group `json:"Resources"`
}

// ListUsersResponse represent a list users response entity
type ListUsersResponse struct {
	GeneralResponse
	Resources []*User `json:"Resources"`
}

// CreateGroupResponse represent a create group response entity
type CreateGroupResponse struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

// CreateGroupRequest represent a create group request entity
type CreateGroupRequest struct {
	DisplayName string    `json:"displayName"`
	ExternalID  string    `json:"externalId,omitempty"`
	Members     []*Member `json:"members,omitempty"`
}

// CreateUserRequest represent a create user request entity
type CreateUserRequest struct {
	ID          string  `json:"id,omitempty"`
	ExternalID  string  `json:"externalId"`
	UserName    string  `json:"userName"`
	Name        Name    `json:"name"`
	DisplayName string  `json:"displayName"`
	NickName    string  `json:"nickName,omitempty"`
	ProfileURL  string  `json:"profileURL,omitempty"`
	Active      bool    `json:"active"`
	Emails      []Email `json:"emails"`
}

// PutUserRequest represent a put user request entity
type PutUserRequest struct {
	ID          string  `json:"id"`
	ExternalID  string  `json:"externalId"`
	UserName    string  `json:"userName"`
	Name        Name    `json:"name"`
	DisplayName string  `json:"displayName"`
	NickName    string  `json:"nickName,omitempty"`
	ProfileURL  string  `json:"profileURL,omitempty"`
	Active      bool    `json:"active"`
	Emails      []Email `json:"emails"`
}

// CreateUserResponse represent a create user response entity
type CreateUserResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Active      bool     `json:"active"`
	Emails      []Email  `json:"emails"`
}

// PutUserResponse represent a put user response entity
type PutUserResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Active      bool     `json:"active"`
	Emails      []Email  `json:"emails"`
}

// PatchGroupRequest represent a patch group request entity
type PatchGroupRequest struct {
	Group Group      `json:"group"`
	Patch PatchGroup `json:"patch"`
}

// PatchUserRequest represent a patch user request entity
type PatchUserRequest struct {
	User  User  `json:"user"`
	Patch Patch `json:"patch"`
}

// GetUserResponse represent a get user response entity
type GetUserResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Active      bool     `json:"active"`
	Emails      []Email  `json:"emails"`
}
