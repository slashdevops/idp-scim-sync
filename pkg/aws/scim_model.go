package aws

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v3"
)

type Meta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
}

type Group struct {
	ID          string   `json:"id"`
	Meta        Meta     `json:"meta"`
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
	Members     []string `json:"members"`
}

type User struct {
	ID         string   `json:"id"`
	ExternalId string   `json:"externalId"`
	Meta       Meta     `json:"meta"`
	Schemas    []string `json:"schemas"`
	UserName   string   `json:"userName"`
	Name       struct {
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
	} `json:"name"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Emails      []struct {
		Value   string `json:"value"`
		Type    string `json:"type"`
		Primary bool   `json:"primary"`
	} `json:"emails"`
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
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Schemas      []string `json:"schemas"`
}

type GroupsResponse struct {
	GeneralResponse
	Resources []*Group `json:"Resources"`
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
	Resources []*User `json:"Resources"`
}

type APIErrorResponse struct {
	Status             string   `json:"status"`
	Detail             string   `json:"detail"`
	Timestamp          string   `json:"timestamp"`
	ExceptionRequestId string   `json:"exceptionRequestId"`
	Schema             []string `json:"schema"`
}
