package google

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
)

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
