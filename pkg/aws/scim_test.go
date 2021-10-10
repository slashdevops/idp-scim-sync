package aws

import (
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
)

func Test_NewSCIMService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return AWSSCIMProvider", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("Should return ErrParsingURL error", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "https://%%testing.com", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrParsingURL)
		assert.Nil(t, got)
	})

	t.Run("Should return ErrEndpointEmpty error ", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrEndpointEmpty)
		assert.Nil(t, got)
	})
}

func Test_AWSSCIMProvider_EndpointURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return valid url", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		endpoint := "https://testing.com"

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		url, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, url)

		assert.Equal(t, url, got.EndpointURL())
	})
}
