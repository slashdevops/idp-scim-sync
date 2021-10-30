package aws

// Name represent a name entity
type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName"`
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

// Patch represent a patch entity
type Patch struct {
	Schemas    []string    `json:"schemas"`
	Operations []Operation `json:"Operations"`
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
	Members     []Member `json:"members"`
}

// User represent a user entity
type User struct {
	ID          string   `json:"id"`
	ExternalId  string   `json:"externalId"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Active      bool     `json:"active"`
	Emails      []Email  `json:"emails"`
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

// APIErrorResponse represent an api error response entity
type APIErrorResponse struct {
	Status             string   `json:"status"`
	Detail             string   `json:"detail"`
	Timestamp          string   `json:"timestamp"`
	ExceptionRequestId string   `json:"exceptionRequestId"`
	Schema             []string `json:"schema"`
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
	DisplayName string   `json:"displayName"`
	Members     []Member `json:"members"`
}

// CreateUserRequest represent a create user request entity
type CreateUserRequest struct {
	ID          string  `json:"id"`
	ExternalId  string  `json:"externalId"`
	UserName    string  `json:"userName"`
	Name        Name    `json:"name"`
	DisplayName string  `json:"displayName"`
	NickName    string  `json:"nickName"`
	ProfileURL  string  `json:"profileURL"`
	Active      bool    `json:"active"`
	Emails      []Email `json:"emails"`
}

// PutUserRequest represent a put user request entity
type PutUserRequest struct {
	ID          string  `json:"id"`
	ExternalId  string  `json:"externalId"`
	UserName    string  `json:"userName"`
	Name        Name    `json:"name"`
	DisplayName string  `json:"displayName"`
	NickName    string  `json:"nickName"`
	ProfileURL  string  `json:"profileURL"`
	Active      bool    `json:"active"`
	Emails      []Email `json:"emails"`
}

// CreateUserResponse represent a create user response entity
type CreateUserResponse struct {
	ID          string   `json:"id"`
	ExternalId  string   `json:"externalId"`
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
	ExternalId  string   `json:"externalId"`
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
	Group Group `json:"group"`
	Patch Patch `json:"patch"`
}

// PatchUserRequest represent a patch user request entity
type PatchUserRequest struct {
	User  User  `json:"user"`
	Patch Patch `json:"patch"`
}
