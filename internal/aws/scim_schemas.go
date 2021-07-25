package aws

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v3"
)

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

func (c *ServiceProviderConfig) ToJSON() string {
	out, err := json.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	return string(out)
}

func (c *ServiceProviderConfig) ToYAML() string {
	out, err := yaml.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	return string(out)
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

func (c *GroupsResponse) ToJSON() string {
	out, err := json.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	return string(out)
}

func (c *GroupsResponse) ToYAML() string {
	out, err := yaml.Marshal(c)
	if err != nil {
		log.Panic(err)
	}
	return string(out)
}

type UsersResponse struct {
	GeneralResponse
	Resources []*User `json:"Resources,omitempty"`
}
