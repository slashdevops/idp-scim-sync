package google

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
)

func TestNewService(t *testing.T) {
	t.Run("New Service with parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"
		scope := "admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope, admin.AdminDirectoryUserReadonlyScope"

		serviceAccount, err := ioutil.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		svc, err := NewService(ctx, userEmail, serviceAccount, scope)
		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("New Service without parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := ""
		scope := ""

		svc, err := NewService(ctx, userEmail, nil, scope)
		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

func TestNewDirectoryService(t *testing.T) {
	t.Run("New Directory Service Client with parameters", func(t *testing.T) {
		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"
		scope := "admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope, admin.AdminDirectoryUserReadonlyScope"

		serviceAccount, err := ioutil.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		svc, err := NewService(ctx, userEmail, serviceAccount, scope)
		if err != nil {
			t.Fatalf("Error creating a service: %s", err)
		}

		client, err := NewDirectoryService(ctx, svc)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("New Directory Service Client without right parameters", func(t *testing.T) {
		ctx := context.TODO()
		svc := &admin.Service{}

		client, err := NewDirectoryService(ctx, svc)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
