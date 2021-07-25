package aws

type Meta struct {
	ResourceType string `json:"resourceType,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

type Group struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta,omitempty"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
	Members     []string `json:"members,omitempty"`
}

type User struct {
	ID         string   `json:"id"`
	ExternalId string   `json:"externalId,omitempty"`
	Meta       Meta     `json:"meta,omitempty"`
	Schemas    []string `json:"schemas"`
	UserName   string   `json:"userName"`
	Name       struct {
		FamilyName string `json:"familyName,omitempty"`
		GivenName  string `json:"givenName,omitempty"`
	} `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Active      bool   `json:"active"`
	Emails      []struct {
		Value   string `json:"value,omitempty"`
		Type    string `json:"type,omitempty"`
		Primary bool   `json:"primary,omitempty"`
	} `json:"emails,omitempty"`
}

type ServiceProviderConfig struct {
	Schemas               []string `json:"schemas,omitempty"`
	DocumentationURI      string   `json:"documentationUri,omitempty"`
	AuthenticationSchemes []struct {
		Type             string `json:"type,omitempty"`
		Name             string `json:"name,omitempty"`
		Description      string `json:"description,omitempty"`
		SpecURI          string `json:"specURI,omitempty"`
		DocumentationURI string `json:"documentationUri,omitempty"`
		Primary          bool   `json:"primary,omitempty"`
	} `json:"authenticationSchemes,omitempty"`
	Patch struct {
		Supported bool `json:"supported,omitempty"`
	} `json:"patch,omitempty"`
	Bulk struct {
		Supported      bool `json:"supported,omitempty"`
		MaxOperations  int  `json:"maxOperations,omitempty"`
		MaxPayloadSize int  `json:"maxPayloadSize,omitempty"`
	} `json:"bulk,omitempty"`
	Filter struct {
		Supported  bool `json:"supported,omitempty"`
		MaxResults int  `json:"maxResults,omitempty"`
	} `json:"filter,omitempty"`
	ChangePassword struct {
		Supported bool `json:"supported,omitempty"`
	} `json:"changePassword,omitempty"`
	Sort struct {
		Supported bool `json:"supported,omitempty"`
	} `json:"sort,omitempty"`
	Etag struct {
		Supported bool `json:"supported,omitempty"`
	} `json:"etag,omitempty"`
}

type GeneralResponse struct {
	TotalResults int      `json:"totalResults,omitempty"`
	ItemsPerPage int      `json:"itemsPerPage,omitempty"`
	StartIndex   int      `json:"startIndex,omitempty"`
	Schemas      []string `json:"schemas,omitempty"`
}

type GroupsResponse struct {
	GeneralResponse
	Resources []*Group `json:"Resources,omitempty"`
}

type UsersResponse struct {
	GeneralResponse
	Resources []*User `json:"Resources,omitempty"`
}
