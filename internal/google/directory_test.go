package google

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDirectoryService(t *testing.T) {

	t.Run("New Client with parameters", func(t *testing.T) {

		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "testdata/service_account.json"

		serviceAccount, err := ioutil.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		client, err := NewDirectoryService(ctx, userEmail, serviceAccount)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("New Client without right parameters", func(t *testing.T) {

		ctx := context.TODO()
		userEmail := ""
		serviceAccount := []byte("")

		client, err := NewDirectoryService(ctx, userEmail, serviceAccount)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
