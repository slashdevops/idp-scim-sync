package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/slashdevops/idp-scim-sync/internal/idp"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/scim"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/core"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func TestSyncService_NewSyncService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("New Service with parameters", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)
		mockStateRepository := mocks.NewMockStateRepository(mockCtrl)

		svc, err := NewSyncService(mockProviderService, mockSCIMService, mockStateRepository)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("New Service without parameters", func(t *testing.T) {
		svc, err := NewSyncService(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("New Service without IdentityProviderService, return specific error", func(t *testing.T) {
		svc, err := NewSyncService(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrIdentityProviderServiceNil)
	})

	t.Run("New Service without SCIMService, return context specific error", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)

		svc, err := NewSyncService(mockProviderService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrSCIMServiceNil)
	})

	t.Run("New Service without Repository, return context specific error", func(t *testing.T) {
		mockProviderService := mocks.NewMockIdentityProviderService(mockCtrl)
		mockSCIMService := mocks.NewMockSCIMService(mockCtrl)

		svc, err := NewSyncService(mockProviderService, mockSCIMService, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.Equal(t, err, ErrStateRepositoryNil)
	})
}

func TestSyncService_SyncGroupsAndTheirMembers(t *testing.T) {
	ctx := context.TODO()

	t.Run("create empty state file when no date came from idp and scim", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		stateFile, err := os.CreateTemp(tmpDir, "state.json")
		assert.NoError(t, err)
		assert.NotNil(t, stateFile)
		defer stateFile.Close()
		defer os.Remove(stateFile.Name())

		// mock Google Workspace API calls
		svrIDP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Calling IDP API with method: %s, path: %s, query: %s", r.Method, r.URL.Path, r.URL.RawQuery)
			w.Write([]byte(`{}`))
		}))
		defer svrIDP.Close()

		// mock Google Workspace API calls
		svrSCIM := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Calling SCIM API with method: %s, path: %s, query: %s", r.Method, r.URL.Path, r.URL.RawQuery)
			w.Write([]byte(`{}`))
		}))
		defer svrSCIM.Close()

		svc := createService(t, ctx, svrIDP, svrSCIM, stateFile)

		err = svc.SyncGroupsAndTheirMembers(ctx)
		assert.NoError(t, err)

		// check if state file is created
		stateFileCreated, err := os.Stat(stateFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, stateFileCreated)
		assert.Greater(t, stateFileCreated.Size(), int64(0))

		// check if state contains expected data
		jsonSate, err := os.Open(stateFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, jsonSate)
		defer jsonSate.Close()

		jsonStateBytes, err := io.ReadAll(jsonSate)
		assert.NoError(t, err)
		assert.NotNil(t, jsonStateBytes)
		assert.Greater(t, len(jsonStateBytes), 0)

		var state model.State
		err = json.Unmarshal(jsonStateBytes, &state)
		assert.NoError(t, err)
		assert.NotNil(t, state)

		assert.Equal(t, 0, len(state.Resources.Groups.Resources))
		assert.Equal(t, 0, len(state.Resources.Users.Resources))
		assert.Equal(t, 0, len(state.Resources.GroupsMembers.Resources))
		assert.NotEqual(t, "", state.LastSync)
		assert.NotEqual(t, "", state.HashCode)
		assert.Equal(t, "", state.CodeVersion)
		assert.Equal(t, model.StateSchemaVersion, state.SchemaVersion)
	})

	t.Run("create state file when date came from idp and no data from scim", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		stateFile, err := os.CreateTemp(tmpDir, "state.json")
		assert.NoError(t, err)
		assert.NotNil(t, stateFile)
		defer stateFile.Close()
		defer os.Remove(stateFile.Name())

		groupsList := &admin.Groups{
			Etag: "etag-groups",
			Kind: "directory#groups",
			Groups: []*admin.Group{
				{
					Id:    "group-1",
					Etag:  "etag-group-1",
					Email: "group.1@mail.com",
					Name:  "group 1",
				},
				{
					Id:    "group-2",
					Etag:  "etag-group-2",
					Email: "group.2@mail.com",
					Name:  "group 2",
				},
			},
		}
		groupsListJSONBytes, err := groupsList.MarshalJSON()
		assert.NoError(t, err)

		membersList := &admin.Members{
			Etag: "etag-members",
			Kind: "directory#members",
			Members: []*admin.Member{
				{
					Id:     "user-1",
					Etag:   "etag-member-1",
					Email:  "user.1@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
				{
					Id:     "user-2",
					Etag:   "etag-member-2",
					Email:  "user.2@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
			},
		}
		membersListJSONBytes, err := membersList.MarshalJSON()
		assert.NoError(t, err)

		user1 := &admin.User{
			Id:           "user-1",
			Etag:         "etag-user-1",
			PrimaryEmail: "user.1@mail.com",
			Name: &admin.UserName{
				FamilyName: "1",
				GivenName:  "user",
			},
			Suspended: false,
			Emails: []*admin.UserEmail{
				{
					Address: "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
		}

		user1JSONBytes, err := user1.MarshalJSON()
		assert.NoError(t, err)

		user2 := &admin.User{
			Id:           "user-2",
			Etag:         "etag-user-2",
			PrimaryEmail: "user.2@mail.com",
			Name: &admin.UserName{
				FamilyName: "2",
				GivenName:  "user",
			},
			Suspended: false,
			Emails: []*admin.UserEmail{
				{
					Address: "user.2@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
		}

		user2JSONBytes, err := user2.MarshalJSON()
		assert.NoError(t, err)

		createGroup1Response := &aws.CreateGroupResponse{
			ID: "group-1",
			Meta: aws.Meta{
				ResourceType: "Group",
				Created:      "2020-01-01T00:00:00Z",
				LastModified: "2020-01-01T00:00:00Z",
			},
			Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
			DisplayName: "group 1",
		}
		createGroup1ResponseJSONBytes, err := json.Marshal(createGroup1Response)
		assert.NoError(t, err)

		createGroup2Response := &aws.CreateGroupResponse{
			ID: "group-2",
			Meta: aws.Meta{
				ResourceType: "Group",
				Created:      "2020-01-01T00:00:00Z",
				LastModified: "2020-01-01T00:00:00Z",
			},
			Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
			DisplayName: "group 2",
		}
		createGroup2ResponseJSONBytes, err := json.Marshal(createGroup2Response)
		assert.NoError(t, err)

		createUser1Response := &aws.CreateUserResponse{
			ID:         "user-1",
			ExternalID: "user-1",
			Meta: aws.Meta{
				ResourceType: "User",
				Created:      "2020-01-01T00:00:00Z",
				LastModified: "2020-01-01T00:00:00Z",
			},
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			UserName: "user.1@mail.com",
			Name: aws.Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Active:      true,
			Emails: []aws.Email{{
				Value:   "user.1@mail.com",
				Type:    "work",
				Primary: true,
			}},
		}
		createUser1ResponseJSONBytes, err := json.Marshal(createUser1Response)
		assert.NoError(t, err)

		createUser2Response := &aws.CreateUserResponse{
			ID:         "user-2",
			ExternalID: "user-2",
			Meta: aws.Meta{
				ResourceType: "User",
				Created:      "2020-01-01T00:00:00Z",
				LastModified: "2020-01-01T00:00:00Z",
			},
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			UserName: "user.2@mail.com",
			Name: aws.Name{
				FamilyName: "2",
				GivenName:  "user",
			},
			DisplayName: "user 2",
			Active:      true,
			Emails: []aws.Email{{
				Value:   "user.2@mail.com",
				Type:    "work",
				Primary: true,
			}},
		}
		createUser2ResponseJSONBytes, err := json.Marshal(createUser2Response)
		assert.NoError(t, err)

		listGroupsResponseGroup1User1 := &aws.ListGroupsResponse{
			ListResponse: aws.ListResponse{
				StartIndex:   1,
				ItemsPerPage: 1,
				TotalResults: 1,
				Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:ListResponse"},
			},
			Resources: []*aws.Group{
				{
					ID: "group-1",
					Meta: aws.Meta{
						ResourceType: "Group",
						Created:      "2020-01-01T00:00:00Z",
						LastModified: "2020-01-01T00:00:00Z",
					},
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					DisplayName: "group 1",
					Members:     []*aws.Member{}, // AWS SSO SCIM API don't return members for list groups only TotalResults is returned
				},
			},
		}
		listGroupsResponseGroup1User1JSONBytes, err := json.Marshal(listGroupsResponseGroup1User1)
		assert.NoError(t, err)

		listGroupsResponseGroup1User2 := &aws.ListGroupsResponse{
			ListResponse: aws.ListResponse{
				StartIndex:   1,
				ItemsPerPage: 1,
				TotalResults: 0,
				Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:ListResponse"},
			},
			Resources: []*aws.Group{
				{
					ID: "group-1",
					Meta: aws.Meta{
						ResourceType: "Group",
						Created:      "2020-01-01T00:00:00Z",
						LastModified: "2020-01-01T00:00:00Z",
					},
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					DisplayName: "group 1",
					Members:     []*aws.Member{}, // AWS SSO SCIM API don't return members for list groups only TotalResults is returned
				},
			},
		}
		listGroupsResponseGroup1User2JSONBytes, err := json.Marshal(listGroupsResponseGroup1User2)
		assert.NoError(t, err)

		listGroupsResponseGroup2User1 := &aws.ListGroupsResponse{
			ListResponse: aws.ListResponse{
				StartIndex:   1,
				ItemsPerPage: 1,
				TotalResults: 0,
				Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:ListResponse"},
			},
			Resources: []*aws.Group{
				{
					ID: "group-2",
					Meta: aws.Meta{
						ResourceType: "Group",
						Created:      "2020-01-01T00:00:00Z",
						LastModified: "2020-01-01T00:00:00Z",
					},
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					DisplayName: "group 2",
					Members:     []*aws.Member{}, // AWS SSO SCIM API don't return members for list groups only TotalResults is returned
				},
			},
		}
		listGroupsResponseGroup2User1JSONBytes, err := json.Marshal(listGroupsResponseGroup2User1)
		assert.NoError(t, err)

		listGroupsResponseGroup2User2 := &aws.ListGroupsResponse{
			ListResponse: aws.ListResponse{
				StartIndex:   1,
				ItemsPerPage: 1,
				TotalResults: 1,
				Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:ListResponse"},
			},
			Resources: []*aws.Group{
				{
					ID: "group-2",
					Meta: aws.Meta{
						ResourceType: "Group",
						Created:      "2020-01-01T00:00:00Z",
						LastModified: "2020-01-01T00:00:00Z",
					},
					Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
					DisplayName: "group 2",
					Members:     []*aws.Member{}, // AWS SSO SCIM API don't return members for list groups only TotalResults is returned
				},
			},
		}
		listGroupsResponseGroup2User2JSONBytes, err := json.Marshal(listGroupsResponseGroup2User2)
		assert.NoError(t, err)

		// mock Google Workspace API calls
		svrIDP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Calling IDP API with method: %s, path: %s, query: %s", r.Method, r.URL.Path, r.URL.RawQuery)
			switch r.URL.Path {
			case "/admin/directory/v1/groups":
				w.Write(groupsListJSONBytes)
			case "/admin/directory/v1/groups/group-1/members":
				w.Write(membersListJSONBytes)
			case "/admin/directory/v1/groups/group-2/members":
				w.Write(membersListJSONBytes)
			case "/admin/directory/v1/users/user.1@mail.com":
				w.Write(user1JSONBytes)
			case "/admin/directory/v1/users/user.2@mail.com":
				w.Write(user2JSONBytes)
			}
		}))
		defer svrIDP.Close()

		// mock Google Workspace API calls
		svrSCIM := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Calling SCIM API with method: %s, path: %s, query: %s", r.Method, r.URL.Path, r.URL.RawQuery)
			switch r.Method {
			case "GET":
				switch r.URL.Path {
				case "/Groups":
					filter := r.URL.Query().Get("filter")

					switch filter {
					case "": // first time getting groups
						w.Write([]byte(`{}`))
					case "id eq \"group-1\" and members eq \"user-1\"":
						w.Write(listGroupsResponseGroup1User1JSONBytes)
					case "id eq \"group-1\" and members eq \"user-2\"":
						w.Write(listGroupsResponseGroup1User2JSONBytes) // user 2 is not in group 1
					case "id eq \"group-2\" and members eq \"user-1\"":
						w.Write(listGroupsResponseGroup2User1JSONBytes) // user 1 is not in group 2
					case "id eq \"group-2\" and members eq \"user-2\"":
						w.Write(listGroupsResponseGroup2User2JSONBytes)
					default:
						w.WriteHeader(http.StatusBadRequest)
					}
				case "/Users":
					w.Write([]byte(`{}`))
				}
			case "POST":
				var bodyData map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
					t.Errorf("Error decoding body: %s", err)
				}
				// t.Logf("Body: %s", utils.ToJSON(bodyData))

				switch r.URL.Path {
				case "/Groups":
					switch bodyData["displayName"] {
					case "group 1":
						w.Write(createGroup1ResponseJSONBytes)
					case "group 2":
						w.Write(createGroup2ResponseJSONBytes)
					default:
						w.WriteHeader(http.StatusBadRequest)
					}
				case "/Users":
					switch bodyData["userName"] {
					case "user.1@mail.com":
						w.Write(createUser1ResponseJSONBytes)
					case "user.2@mail.com":
						w.Write(createUser2ResponseJSONBytes)
					default:
						w.WriteHeader(http.StatusBadRequest)
					}
				}
			}
		}))
		defer svrSCIM.Close()

		svc := createService(t, ctx, svrIDP, svrSCIM, stateFile)

		err = svc.SyncGroupsAndTheirMembers(ctx)
		assert.NoError(t, err)

		// check if state file is created
		stateFileCreated, err := os.Stat(stateFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, stateFileCreated)
		assert.Greater(t, stateFileCreated.Size(), int64(0))

		// check if state contains expected data
		jsonSate, err := os.Open(stateFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, jsonSate)
		defer jsonSate.Close()

		jsonStateBytes, err := io.ReadAll(jsonSate)
		assert.NoError(t, err)
		assert.NotNil(t, jsonStateBytes)
		assert.Greater(t, len(jsonStateBytes), 0)

		var state model.State
		err = json.Unmarshal(jsonStateBytes, &state)
		assert.NoError(t, err)
		assert.NotNil(t, state)

		// t.Logf("State: %s", utils.ToJSON(state))
		assert.Equal(t, 2, len(state.Resources.Groups.Resources))
		assert.Equal(t, 2, len(state.Resources.Users.Resources))
		assert.Equal(t, 4, len(state.Resources.GroupsMembers.Resources))
		assert.NotEqual(t, "", state.LastSync)
		assert.NotEqual(t, "", state.HashCode)
		assert.Equal(t, "", state.CodeVersion)
		assert.Equal(t, model.StateSchemaVersion, state.SchemaVersion)
		assert.Equal(t, 2, state.Resources.Groups.Items)
		assert.Equal(t, 2, state.Resources.Users.Items)
	})
}

// createService helper function to create a new SyncService instance
func createService(
	t *testing.T,
	ctx context.Context,
	idpSRV *httptest.Server,
	scimSRV *httptest.Server,
	stateFile io.ReadWriter,
) *SyncService {
	googleSvc, err := admin.NewService(ctx, option.WithHTTPClient(idpSRV.Client()), option.WithEndpoint(idpSRV.URL), option.WithUserAgent("test"))
	assert.NoError(t, err)

	gwsDS, err := google.NewDirectoryService(googleSvc)
	assert.NoError(t, err)
	assert.NotNil(t, gwsDS)

	awsSCIM, err := aws.NewSCIMService(scimSRV.Client(), scimSRV.URL, "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, awsSCIM)

	// Identity Provider Service
	idpService, err := idp.NewIdentityProvider(gwsDS)
	assert.NoError(t, err)
	assert.NotNil(t, idpService)

	// AWS SCIM Service
	scimService, err := scim.NewProvider(awsSCIM)
	assert.NoError(t, err)
	assert.NotNil(t, scimService)

	// Disk State Repository
	repo, err := repository.NewDiskRepository(stateFile)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	svc, err := NewSyncService(idpService, scimService, repo)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	return svc
}
