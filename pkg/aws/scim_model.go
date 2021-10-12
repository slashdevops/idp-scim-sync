package aws

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName"`
}

type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary"`
}

type Meta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
}

type Memeber struct {
	Value string `json:"value"`
	Ref   string `json:"$ref"`
	Type  string `json:"type"`
}

type Group struct {
	ID          string    `json:"id"`
	Meta        Meta      `json:"meta"`
	Schemas     []string  `json:"schemas"`
	DisplayName string    `json:"displayName"`
	Members     []Memeber `json:"members"`
}

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

type GeneralResponse struct {
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Schemas      []string `json:"schemas"`
}

type ListsGroupsResponse struct {
	GeneralResponse
	Resources []*Group `json:"Resources"`
}

type UsersResponse struct {
	GeneralResponse
	Resources []*User `json:"Resources"`
}

type APIErrorResponse struct {
	Status             string   `json:"status"`
	Detail             string   `json:"detail"`
	Timestamp          string   `json:"timestamp"`
	ExceptionRequestId string   `json:"exceptionRequestId"`
	Schema             []string `json:"schema"`
}

type CreateGroupResponse struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

type CreateGroupRequest struct {
	DisplayName string    `json:"displayName"`
	Memebers    []Memeber `json:"members"`
}

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
