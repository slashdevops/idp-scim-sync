package google

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func TestListUsers(t *testing.T) {
}

func TestNewService(t *testing.T) {
	t.Run("Should return a new Service with mocked parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"
		scope := "admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope, admin.AdminDirectoryUserReadonlyScope"

		serviceAccount, err := os.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		svc, err := NewService(ctx, userEmail, serviceAccount, scope)
		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("Should return a new Service with empty service account parameter", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := ""
		scope := ""

		svc, err := NewService(ctx, userEmail, nil, scope)
		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("Should return an error when scope is nil", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := ""

		svc, err := NewService(ctx, userEmail, nil)
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
		scope := "admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope, admin.AdminDirectoryUserReadonlyScope"

		serviceAccount, err := os.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		svc, err := NewService(ctx, userEmail, serviceAccount, scope)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write([]byte(`{"etag":"etag-users","kind":"directory#users","users":[{"id":"123456789","etag":"etag-user-123456789","primaryEmail":"""}]}`))
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
			w.Write(jsonBytes)
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
