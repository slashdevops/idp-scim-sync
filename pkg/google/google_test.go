package google

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func TestNewService(t *testing.T) {
	t.Run("Should return a new Service with mocked parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"
		scopes := []string{
			"https://www.googleapis.com/auth/admin.directory.group.readonly",
			"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
			"https://www.googleapis.com/auth/admin.directory.user.readonly",
		}

		serviceAccount, err := os.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		config := DirectoryServiceConfig{
			Client:         &http.Client{},
			UserEmail:      userEmail,
			ServiceAccount: serviceAccount,
			Scopes:         scopes,
			UserAgent:      "test-agent",
		}

		svc, err := NewService(ctx, config)
		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return a new Service with empty service account parameter", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := ""
		scopes := []string{}

		config := DirectoryServiceConfig{
			Client:         &http.Client{},
			UserEmail:      userEmail,
			ServiceAccount: nil,
			Scopes:         scopes,
			UserAgent:      "test-agent",
		}

		svc, err := NewService(ctx, config)
		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("Should return an error when scope is nil", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "test@example.com"

		config := DirectoryServiceConfig{
			Client:         &http.Client{},
			UserEmail:      userEmail,
			ServiceAccount: []byte("{}"), // provide minimal valid JSON
			Scopes:         nil,
			UserAgent:      "test-agent",
		}

		svc, err := NewService(ctx, config)
		assert.Error(t, err)
		assert.Nil(t, svc)
		assert.ErrorIs(t, err, ErrGoogleClientScopeNil)
	})
}

func TestNewDirectoryService(t *testing.T) {
	t.Run("Should return a new Directory Service Client with mocked parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"
		scopes := []string{
			"https://www.googleapis.com/auth/admin.directory.group.readonly",
			"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
			"https://www.googleapis.com/auth/admin.directory.user.readonly",
		}

		serviceAccount, err := os.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		config := DirectoryServiceConfig{
			Client:         &http.Client{},
			UserEmail:      userEmail,
			ServiceAccount: serviceAccount,
			Scopes:         scopes,
			UserAgent:      "test-agent",
		}

		svc, err := NewService(ctx, config)
		if err != nil {
			t.Fatalf("Error creating a service: %s", err)
		}

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Should return a new Directory Service Client", func(t *testing.T) {
		svc := &admin.Service{}

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestNewDirectoryService_ListUsers(t *testing.T) {
	t.Run("should return a valid list of two users with nil argument", func(t *testing.T) {
		ctx := context.TODO()
		urlPath := "/admin/directory/v1/users"

		userList := &admin.Users{
			Etag: "etag-users",
			Kind: "directory#users",
			Users: []*admin.User{
				{
					Id:           "123456789",
					Etag:         "etag-user-123456789",
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
				},
				{
					Id:           "987654321",
					Etag:         "etag-user-987654321",
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
				},
			},
		}
		jsonBytes, err := userList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, nil)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-user-123456789", got[0].Etag)
		assert.Equal(t, "etag-user-987654321", got[1].Etag)
		assert.Equal(t, "user.1@mail.com", got[0].PrimaryEmail)
		assert.Equal(t, "user.2@mail.com", got[1].PrimaryEmail)
		assert.Equal(t, "1", got[0].Name.FamilyName)
		assert.Equal(t, "user", got[0].Name.GivenName)
		assert.Equal(t, "2", got[1].Name.FamilyName)
		assert.Equal(t, "user", got[1].Name.GivenName)
		assert.False(t, got[0].Suspended)
		assert.False(t, got[1].Suspended)
	})

	t.Run("should return a valid list of two users with empty argument", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{""}
		urlPath := "/admin/directory/v1/users"

		userList := &admin.Users{
			Etag: "etag-users",
			Kind: "directory#users",
			Users: []*admin.User{
				{
					Id:           "123456789",
					Etag:         "etag-user-123456789",
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
				},
				{
					Id:           "987654321",
					Etag:         "etag-user-987654321",
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
				},
			},
		}
		jsonBytes, err := userList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, filter)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-user-123456789", got[0].Etag)
		assert.Equal(t, "etag-user-987654321", got[1].Etag)
		assert.Equal(t, "user.1@mail.com", got[0].PrimaryEmail)
		assert.Equal(t, "user.2@mail.com", got[1].PrimaryEmail)
		assert.Equal(t, "1", got[0].Name.FamilyName)
		assert.Equal(t, "user", got[0].Name.GivenName)
		assert.Equal(t, "2", got[1].Name.FamilyName)
		assert.Equal(t, "user", got[1].Name.GivenName)
		assert.False(t, got[0].Suspended)
		assert.False(t, got[1].Suspended)
	})

	t.Run("should return a valid list of one users with filter argument", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{"name:user* email:user*"}
		urlPath := "/admin/directory/v1/users"

		userList := &admin.Users{
			Etag: "etag-users",
			Kind: "directory#users",
			Users: []*admin.User{
				{
					Id:           "123456789",
					Etag:         "etag-user-123456789",
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
				},
			},
		}
		jsonBytes, err := userList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			assert.Equal(t, filter[0], r.URL.Query().Get("query"))
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, filter)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "etag-user-123456789", got[0].Etag)
		assert.Equal(t, "user.1@mail.com", got[0].PrimaryEmail)
		assert.Equal(t, "1", got[0].Name.FamilyName)
		assert.Equal(t, "user", got[0].Name.GivenName)
		assert.False(t, got[0].Suspended)
	})

	t.Run("show return an error when pages() return error because the json is bad", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{"name:user* email:user*"}
		urlPath := "/admin/directory/v1/users"

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			assert.Equal(t, filter[0], r.URL.Query().Get("query"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"etag":"etag-users","kind":"directory#users","users":[{"id":"123456789","etag":"etag-user-123456789","primaryEmail":"""}]}`))
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, filter)
		assert.Error(t, err)

		t.Logf("got: %+v", got)
		assert.Equal(t, 0, len(got))
	})

	t.Run("show return an error when pages() return error", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{"name:user* email:user*"}
		urlPath := "/admin/directory/v1/users"

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			assert.Equal(t, filter[0], r.URL.Query().Get("query"))
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, filter)
		assert.Error(t, err)

		t.Logf("got: %+v", got)
		assert.Equal(t, 0, len(got))
	})

	t.Run("show return an error when pages() return error", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{"name:user* email:user*"}
		urlPath := "/admin/directory/v1/users"

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			assert.Equal(t, filter[0], r.URL.Query().Get("query"))
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListUsers(ctx, filter)
		assert.Error(t, err)

		t.Logf("got: %+v", got)
		assert.Equal(t, 0, len(got))
	})
}

func TestNewDirectoryService_ListGroups(t *testing.T) {
	t.Run("should return a valid list of two groups with nil argument", func(t *testing.T) {
		ctx := context.TODO()

		groupsList := &admin.Groups{
			Etag: "etag-groups",
			Kind: "directory#groups",
			Groups: []*admin.Group{
				{
					Id:    "123456789",
					Etag:  "etag-group-123456789",
					Email: "group.1@mail.com",
					Name:  "group 1",
				},
				{
					Id:    "987654321",
					Etag:  "etag-group-987654321",
					Email: "group.2@mail.com",
					Name:  "group 2",
				},
			},
		}
		jsonBytes, err := groupsList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroups(ctx, nil)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-group-123456789", got[0].Etag)
		assert.Equal(t, "etag-group-987654321", got[1].Etag)
		assert.Equal(t, "group.1@mail.com", got[0].Email)
		assert.Equal(t, "group.2@mail.com", got[1].Email)
		assert.Equal(t, "group 1", got[0].Name)
		assert.Equal(t, "group 2", got[1].Name)
	})

	t.Run("should return a valid list of two users with empty argument", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{""}

		groupsList := &admin.Groups{
			Etag: "etag-groups",
			Kind: "directory#groups",
			Groups: []*admin.Group{
				{
					Id:    "123456789",
					Etag:  "etag-group-123456789",
					Email: "group.1@mail.com",
					Name:  "group 1",
				},
				{
					Id:    "987654321",
					Etag:  "etag-group-987654321",
					Email: "group.2@mail.com",
					Name:  "group 2",
				},
			},
		}
		jsonBytes, err := groupsList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroups(ctx, filter)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-group-123456789", got[0].Etag)
		assert.Equal(t, "etag-group-987654321", got[1].Etag)
		assert.Equal(t, "group.1@mail.com", got[0].Email)
		assert.Equal(t, "group.2@mail.com", got[1].Email)
		assert.Equal(t, "group 1", got[0].Name)
		assert.Equal(t, "group 2", got[1].Name)
	})

	t.Run("should return a valid list of one users with filter argument", func(t *testing.T) {
		ctx := context.TODO()

		filter := []string{"name:user* email:user*"}

		groupsList := &admin.Groups{
			Etag: "etag-groups",
			Kind: "directory#groups",
			Groups: []*admin.Group{
				{
					Id:    "123456789",
					Etag:  "etag-group-123456789",
					Email: "group.1@mail.com",
					Name:  "group 1",
				},
				{
					Id:    "987654321",
					Etag:  "etag-group-987654321",
					Email: "group.2@mail.com",
					Name:  "group 2",
				},
			},
		}
		jsonBytes, err := groupsList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, filter[0], r.URL.Query().Get("query"))
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroups(ctx, filter)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-group-123456789", got[0].Etag)
		assert.Equal(t, "etag-group-987654321", got[1].Etag)
		assert.Equal(t, "group.1@mail.com", got[0].Email)
		assert.Equal(t, "group.2@mail.com", got[1].Email)
		assert.Equal(t, "group 1", got[0].Name)
		assert.Equal(t, "group 2", got[1].Name)
	})
}

func TestNewDirectoryService_ListGroupMembers(t *testing.T) {
	t.Run("should return a valid list of two members with groupId empty", func(t *testing.T) {
		ctx := context.TODO()

		groupID := ""

		membersList := &admin.Members{
			Etag: "etag-members",
			Kind: "directory#members",
			Members: []*admin.Member{
				{
					Id:     "123456789",
					Etag:   "etag-member-123456789",
					Email:  "member.1@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
				{
					Id:     "987654321",
					Etag:   "etag-member-987654321",
					Email:  "member.2@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
			},
		}
		jsonBytes, err := membersList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroupMembers(ctx, groupID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return a valid list of two members with groupId", func(t *testing.T) {
		ctx := context.TODO()

		groupID := "123456789"
		urlPath := fmt.Sprintf("/admin/directory/v1/groups/%s/members", groupID)

		membersList := &admin.Members{
			Etag: "etag-members",
			Kind: "directory#members",
			Members: []*admin.Member{
				{
					Id:     "123456789",
					Etag:   "etag-member-123456789",
					Email:  "member.1@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
				{
					Id:     "987654321",
					Etag:   "etag-member-987654321",
					Email:  "member.2@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
			},
		}
		jsonBytes, err := membersList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroupMembers(ctx, groupID)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(got))
		assert.Equal(t, "123456789", got[0].Id)
		assert.Equal(t, "987654321", got[1].Id)
		assert.Equal(t, "etag-member-123456789", got[0].Etag)
		assert.Equal(t, "etag-member-987654321", got[1].Etag)
		assert.Equal(t, "member.1@mail.com", got[0].Email)
		assert.Equal(t, "member.2@mail.com", got[1].Email)
	})

	t.Run("should return a valid list of one members when a member has Status DEACTIVATE", func(t *testing.T) {
		ctx := context.TODO()

		groupID := "123456789"
		urlPath := fmt.Sprintf("/admin/directory/v1/groups/%s/members", groupID)

		membersList := &admin.Members{
			Etag: "etag-members",
			Kind: "directory#members",
			Members: []*admin.Member{
				{
					Id:     "123456789",
					Etag:   "etag-member-123456789",
					Email:  "member.1@mail.com",
					Status: "DEACTIVATE",
					Type:   "USER",
				},
				{
					Id:     "987654321",
					Etag:   "etag-member-987654321",
					Email:  "member.2@mail.com",
					Status: "ACTIVE",
					Type:   "USER",
				},
			},
		}
		jsonBytes, err := membersList.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.ListGroupMembers(ctx, groupID)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(got))
		assert.Equal(t, "987654321", got[0].Id)
		assert.Equal(t, "etag-member-987654321", got[0].Etag)
		assert.Equal(t, "member.2@mail.com", got[0].Email)
	})
}

func TestNewDirectoryService_GetUser(t *testing.T) {
	t.Run("should return error with userId empty", func(t *testing.T) {
		ctx := context.TODO()

		userID := ""

		user := &admin.User{
			Id:           "123456789",
			Etag:         "etag-user-123456789",
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

		jsonBytes, err := user.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.GetUser(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return a valid user with userId not empty", func(t *testing.T) {
		ctx := context.TODO()

		userID := "123456789"
		urlPath := fmt.Sprintf("/admin/directory/v1/users/%s", userID)

		user := &admin.User{
			Id:           "123456789",
			Etag:         "etag-user-123456789",
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

		jsonBytes, err := user.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.GetUser(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "123456789", got.Id)
		assert.Equal(t, "etag-user-123456789", got.Etag)
		assert.Equal(t, "user.1@mail.com", got.PrimaryEmail)
		assert.Equal(t, "1", got.Name.FamilyName)
		assert.Equal(t, "user", got.Name.GivenName)
		assert.False(t, got.Suspended)
	})
}

func TestNewDirectoryService_GetGroup(t *testing.T) {
	t.Run("should return a error groupId empty", func(t *testing.T) {
		ctx := context.TODO()

		groupID := ""

		group := &admin.Group{
			Id:    "123456789",
			Etag:  "etag-group-123456789",
			Email: "group.1@mail.com",
			Name:  "group 1",
		}

		jsonBytes, err := group.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.GetGroup(ctx, groupID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return a valid  group with groupId not empty", func(t *testing.T) {
		ctx := context.TODO()

		groupID := "123456789"
		urlPath := fmt.Sprintf("/admin/directory/v1/groups/%s", groupID)

		group := &admin.Group{
			Id:    "123456789",
			Etag:  "etag-group-123456789",
			Email: "group.1@mail.com",
			Name:  "group 1",
		}

		jsonBytes, err := group.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, urlPath, r.URL.Path)
			_, _ = w.Write(jsonBytes)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		got, err := client.GetGroup(ctx, groupID)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "123456789", got.Id)
		assert.Equal(t, "etag-group-123456789", got.Etag)
		assert.Equal(t, "group.1@mail.com", got.Email)
		assert.Equal(t, "group 1", got.Name)
	})
}

func TestNewDirectoryService_ListGroupMembersBatch(t *testing.T) {
	t.Run("should return empty map for empty groupIDs", func(t *testing.T) {
		ctx := context.TODO()

		svc := &admin.Service{}
		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)

		got, err := client.ListGroupMembersBatch(ctx, []string{})
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})

	t.Run("should return members for multiple groups", func(t *testing.T) {
		ctx := context.TODO()

		group1ID := "group-1"
		group2ID := "group-2"

		members1 := &admin.Members{
			Etag: "etag-1",
			Kind: "directory#members",
			Members: []*admin.Member{
				{Id: "m1", Email: "user1@mail.com", Status: "ACTIVE", Type: "USER"},
			},
		}
		members2 := &admin.Members{
			Etag: "etag-2",
			Kind: "directory#members",
			Members: []*admin.Member{
				{Id: "m2", Email: "user2@mail.com", Status: "ACTIVE", Type: "USER"},
				{Id: "m3", Email: "user3@mail.com", Status: "ACTIVE", Type: "USER"},
			},
		}

		json1, err := members1.MarshalJSON()
		assert.NoError(t, err)
		json2, err := members2.MarshalJSON()
		assert.NoError(t, err)

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			if r.URL.Path == fmt.Sprintf("/admin/directory/v1/groups/%s/members", group1ID) {
				_, _ = w.Write(json1)
			} else if r.URL.Path == fmt.Sprintf("/admin/directory/v1/groups/%s/members", group2ID) {
				_, _ = w.Write(json2)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)

		got, err := client.ListGroupMembersBatch(ctx, []string{group1ID, group2ID})
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Len(t, got, 2)
		assert.Len(t, got[group1ID], 1)
		assert.Len(t, got[group2ID], 2)
		assert.Equal(t, "user1@mail.com", got[group1ID][0].Email)
		assert.Equal(t, "user2@mail.com", got[group2ID][0].Email)
		assert.Equal(t, "user3@mail.com", got[group2ID][1].Email)
	})

	t.Run("should return error when a group fetch fails", func(t *testing.T) {
		ctx := context.TODO()

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer svr.Close()

		svc, err := admin.NewService(ctx, option.WithHTTPClient(svr.Client()), option.WithEndpoint(svr.URL), option.WithUserAgent("test"))
		assert.NoError(t, err)

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)

		got, err := client.ListGroupMembersBatch(ctx, []string{"bad-group"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestWithSyncFieldSet(t *testing.T) {
	t.Run("should apply sync field set to directory service", func(t *testing.T) {
		svc := &admin.Service{}

		fields := model.NewSyncFieldSet([]string{"phoneNumbers"})
		client, err := NewDirectoryService(svc, WithSyncFieldSet(fields))
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// The fields should now be applied — verify internal state
		// by checking that the field strings contain "phones" but not "addresses"
		listFields := string(client.listUsersRequiredFields)
		assert.Contains(t, listFields, "phones")
		assert.NotContains(t, listFields, "addresses")

		getFields := string(client.getUsersRequiredFields)
		assert.Contains(t, getFields, "phones")
		assert.NotContains(t, getFields, "addresses")
	})

	t.Run("should use default fields when no option provided", func(t *testing.T) {
		svc := &admin.Service{}

		client, err := NewDirectoryService(svc)
		assert.NoError(t, err)

		listFields := string(client.listUsersRequiredFields)
		assert.Contains(t, listFields, "phones")
		assert.Contains(t, listFields, "addresses")
		assert.Contains(t, listFields, "organizations")
	})
}
