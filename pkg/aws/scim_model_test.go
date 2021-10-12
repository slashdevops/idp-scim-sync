package aws

import (
	"testing"
)

func Test_ServiceProviderConfig_ToJSON(t *testing.T) {
	type obj struct {
		spc ServiceProviderConfig
	}
	tests := []struct {
		name string
		obj  obj
		want string
	}{
		{
			name: "empty",
			obj: obj{
				spc: ServiceProviderConfig{},
			},
			want: `{"schemas":null,"documentationUri":"","authenticationSchemes":null,"patch":{"supported":false},"bulk":{"supported":false,"maxOperations":0,"maxPayloadSize":0},"filter":{"supported":false,"maxResults":0},"changePassword":{"supported":false},"sort":{"supported":false},"etag":{"supported":false}}`,
		},
		{
			name: "full",
			obj: obj{
				spc: ServiceProviderConfig{
					Schemas:          []string{"urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"},
					DocumentationURI: "http://uri.com/docs",
					AuthenticationSchemes: []struct {
						Type             string `json:"type"`
						Name             string `json:"name"`
						Description      string `json:"description"`
						SpecURI          string `json:"specURI"`
						DocumentationURI string `json:"documentationUri"`
						Primary          bool   `json:"primary"`
					}{
						{
							Type:             "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
							Name:             "email",
							Description:      "Email address",
							SpecURI:          "http://openid.net/specs/openid-connect-core-1_0.html#Claims",
							DocumentationURI: "http://uri.com/docs",
							Primary:          true,
						},
					},
					Patch: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Bulk: struct {
						Supported      bool `json:"supported"`
						MaxOperations  int  `json:"maxOperations"`
						MaxPayloadSize int  `json:"maxPayloadSize"`
					}{
						Supported:      false,
						MaxOperations:  0,
						MaxPayloadSize: 0,
					},
					Filter: struct {
						Supported  bool `json:"supported"`
						MaxResults int  `json:"maxResults"`
					}{
						Supported:  false,
						MaxResults: 0,
					},
					ChangePassword: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Sort: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Etag: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
				},
			},
			want: `{"schemas":["urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"],"documentationUri":"http://uri.com/docs","authenticationSchemes":[{"type":"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress","name":"email","description":"Email address","specURI":"http://openid.net/specs/openid-connect-core-1_0.html#Claims","documentationUri":"http://uri.com/docs","primary":true}],"patch":{"supported":false},"bulk":{"supported":false,"maxOperations":0,"maxPayloadSize":0},"filter":{"supported":false,"maxResults":0},"changePassword":{"supported":false},"sort":{"supported":false},"etag":{"supported":false}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tt.obj.spc
			if got := c.ToJSON(); got != tt.want {
				t.Errorf("ServiceProviderConfig.ToJSON() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_ServiceProviderConfig_ToYAML(t *testing.T) {
	type obj struct {
		spc ServiceProviderConfig
	}
	tests := []struct {
		name string
		obj  obj
		want string
	}{
		{
			name: "empty",
			obj: obj{
				spc: ServiceProviderConfig{},
			},
			want: `schemas: []
documentationuri: ""
authenticationschemes: []
patch:
    supported: false
bulk:
    supported: false
    maxoperations: 0
    maxpayloadsize: 0
filter:
    supported: false
    maxresults: 0
changepassword:
    supported: false
sort:
    supported: false
etag:
    supported: false
`,
		},
		{
			name: "full",
			obj: obj{
				spc: ServiceProviderConfig{
					Schemas:          []string{"urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"},
					DocumentationURI: "http://uri.com/docs",
					AuthenticationSchemes: []struct {
						Type             string `json:"type"`
						Name             string `json:"name"`
						Description      string `json:"description"`
						SpecURI          string `json:"specURI"`
						DocumentationURI string `json:"documentationUri"`
						Primary          bool   `json:"primary"`
					}{
						{
							Type:             "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
							Name:             "email",
							Description:      "Email address",
							SpecURI:          "http://openid.net/specs/openid-connect-core-1_0.html#Claims",
							DocumentationURI: "http://uri.com/docs",
							Primary:          true,
						},
					},
					Patch: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Bulk: struct {
						Supported      bool `json:"supported"`
						MaxOperations  int  `json:"maxOperations"`
						MaxPayloadSize int  `json:"maxPayloadSize"`
					}{
						Supported:      false,
						MaxOperations:  0,
						MaxPayloadSize: 0,
					},
					Filter: struct {
						Supported  bool `json:"supported"`
						MaxResults int  `json:"maxResults"`
					}{
						Supported:  false,
						MaxResults: 0,
					},
					ChangePassword: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Sort: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
					Etag: struct {
						Supported bool `json:"supported"`
					}{
						Supported: false,
					},
				},
			},
			want: `schemas:
    - urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig
documentationuri: http://uri.com/docs
authenticationschemes:
    - type: http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress
      name: email
      description: Email address
      specuri: http://openid.net/specs/openid-connect-core-1_0.html#Claims
      documentationuri: http://uri.com/docs
      primary: true
patch:
    supported: false
bulk:
    supported: false
    maxoperations: 0
    maxpayloadsize: 0
filter:
    supported: false
    maxresults: 0
changepassword:
    supported: false
sort:
    supported: false
etag:
    supported: false
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tt.obj.spc
			if got := c.ToYAML(); got != tt.want {
				t.Errorf("ServiceProviderConfig.ToYAML() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_ListsGroupsResponse_ToJSON(t *testing.T) {
	type obj struct {
		gr ListsGroupsResponse
	}
	tests := []struct {
		name string
		obj  obj
		want string
	}{
		{
			name: "empty",
			obj: obj{
				gr: ListsGroupsResponse{},
			},
			want: `{"totalResults":0,"itemsPerPage":0,"startIndex":0,"schemas":null,"Resources":null}`,
		},
		{
			name: "full",
			obj: obj{
				gr: ListsGroupsResponse{
					GeneralResponse: GeneralResponse{
						TotalResults: 1,
						ItemsPerPage: 1,
						StartIndex:   1,
					},
					Resources: []*Group{
						{
							ID: "id",
							Meta: Meta{
								ResourceType: "ResourceType",
								Created:      "Created",
								LastModified: "LastModified",
							},
							Schemas:     []string{"Schemas"},
							DisplayName: "displayName",
							Members:     []Memeber{{Value: "Value", Ref: "Ref", Type: "Type"}},
						},
					},
				},
			},
			want: `{"totalResults":1,"itemsPerPage":1,"startIndex":1,"schemas":null,"Resources":[{"id":"id","meta":{"resourceType":"ResourceType","created":"Created","lastModified":"LastModified"},"schemas":["Schemas"],"displayName":"displayName","members":[{"value":"Value","$ref":"Ref","type":"Type"}]}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tt.obj.gr
			if got := c.ToJSON(); got != tt.want {
				t.Errorf("ListsGroupsResponse.ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ListsGroupsResponse_ToYAML(t *testing.T) {
	type obj struct {
		gr ListsGroupsResponse
	}
	tests := []struct {
		name string
		obj  obj
		want string
	}{
		{
			name: "empty",
			obj: obj{
				gr: ListsGroupsResponse{},
			},
			want: `generalresponse:
    totalresults: 0
    itemsperpage: 0
    startindex: 0
    schemas: []
resources: []
`,
		},
		{
			name: "full",
			obj: obj{
				gr: ListsGroupsResponse{
					GeneralResponse: GeneralResponse{
						TotalResults: 1,
						ItemsPerPage: 1,
						StartIndex:   1,
					},
					Resources: []*Group{
						{
							ID: "id",
							Meta: Meta{
								ResourceType: "ResourceType",
								Created:      "Created",
								LastModified: "LastModified",
							},
							Schemas:     []string{"Schemas"},
							DisplayName: "displayName",
							Members:     []Memeber{{Value: "Value", Ref: "Ref", Type: "Type"}},
						},
					},
				},
			},
			want: `generalresponse:
    totalresults: 1
    itemsperpage: 1
    startindex: 1
    schemas: []
resources:
    - id: id
      meta:
        resourcetype: ResourceType
        created: Created
        lastmodified: LastModified
      schemas:
        - Schemas
      displayname: displayName
      members:
        - value: Value
          ref: Ref
          type: Type
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tt.obj.gr
			if got := c.ToYAML(); got != tt.want {
				t.Errorf("ListsGroupsResponse.ToYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}
