package aws

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("Should return AWSSCIMProvider when httpClient is nil", func(t *testing.T) {
		got, err := NewSCIMService(nil, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("Should return error when url is bad formed", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "https://%%testing.com", "MyToken")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Should return error when the url is empty ", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrURLEmpty)
		assert.Nil(t, got)
	})
}

func Test_AWSSCIMProvider_request(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Should return error when error come from request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		endpoint := "https://testing.com"

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("test error"))

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		resp, err := got.do(context.TODO(), req, nil)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Should return valid response", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		endpoint := "https://testing.com"

		mockResp := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("Hello, test world!")),
			ContentLength: int64(len("Hello, test world!")),
		}

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		resp, err := got.do(context.TODO(), req, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockResp, resp)
	})
}
