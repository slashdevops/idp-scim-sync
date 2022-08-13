package core

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
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

		stateFile, err := ioutil.TempFile(tmpDir, "state.json")
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

		svc := prepareService(t, ctx, svrIDP, svrSCIM, stateFile)

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

		jsonStateBytes, err := ioutil.ReadAll(jsonSate)
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
}

func prepareService(
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
