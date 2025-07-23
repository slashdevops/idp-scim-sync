package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockErrReader int

func (e mockErrReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func ReadJSONFileAsString(t *testing.T, fileName string) string {
	bytes, err := os.ReadFile(fileName)
	assert.NoError(t, err)

	return string(bytes)
}

func TestNewSCIMService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should return error when url is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, "", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrURLEmpty)
		assert.Nil(t, got)
	})

	t.Run("should return error when token is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, "https://testing.com", "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrBearerTokenEmpty)
		assert.Nil(t, got)
	})

	t.Run("should return AWSSCIMProvider", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should return AWSSCIMProvider when httpClient is nil", func(t *testing.T) {
		got, err := NewSCIMService(nil, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should return error when url is bad formed", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, "https://%%testing.com", "MyToken")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when the url is empty ", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, "", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrURLEmpty)
		assert.Nil(t, got)
	})
}

func TestNewRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	t.Run("valid GET method should return valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mockMethod := http.MethodGet

		mockURL, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, mockURL)

		got, err := service.newRequest(context.Background(), mockMethod, mockURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, mockMethod, got.Method)
		assert.Equal(t, mockURL, got.URL)
		assert.Equal(t, "application/json", got.Header.Get("Accept"))
		assert.Nil(t, got.Body)
	})

	t.Run("valid POST method should return valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mockMethod := http.MethodPost

		mockURL, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, mockURL)

		mockBody := io.NopCloser(strings.NewReader("Hello, test world!"))

		got, err := service.newRequest(context.Background(), mockMethod, mockURL, mockBody)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, mockMethod, got.Method)
		assert.Equal(t, mockURL, got.URL)
		assert.Equal(t, "application/json", got.Header.Get("Accept"))
		assert.NotNil(t, got.Body)
	})

	t.Run("invalid method should return error", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mockMethod := "this is and invalid method"

		mockURL, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, mockURL)

		got, err := service.newRequest(context.Background(), mockMethod, mockURL, nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("valid method should return error when body is wrong", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mockMethod := http.MethodPost

		mockURL, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, mockURL)

		mockBody := map[string]any{
			"this will fail when is serialize": make(chan int),
		}

		got, err := service.newRequest(context.Background(), mockMethod, mockURL, mockBody)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("valid POST method should return valid request and valid userAgent", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mockMethod := http.MethodPost
		service.UserAgent = "MyUserAgent"

		mockURL, err := url.Parse(endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, mockURL)

		mockBody := io.NopCloser(strings.NewReader("Hello, test world!"))

		got, err := service.newRequest(context.Background(), mockMethod, mockURL, mockBody)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, mockMethod, got.Method)
		assert.Equal(t, mockURL, got.URL)
		assert.Equal(t, "application/json", got.Header.Get("Accept"))
		assert.Equal(t, "MyUserAgent", got.Header.Get("User-Agent"))
		assert.NotNil(t, got.Body)
	})
}

func TestDo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	t.Run("should return error when error come from request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("test error"))

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		got, err := service.do(context.Background(), req)
		assert.Error(t, err)

		assert.Nil(t, got)
	})

	t.Run("should return valid response", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		mockResp := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("Hello, test world!")),
			ContentLength: int64(len("Hello, test world!")),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		got, err := service.do(context.Background(), req)
		assert.NoError(t, err)

		assert.NotNil(t, got)
		assert.Equal(t, mockResp, got)
	})
}

func TestCheckHTTPResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	type httpCodes struct {
		code    int
		message string
	}

	validHTTPCodesList := []httpCodes{
		{code: http.StatusOK, message: "200 OK"},
		{code: http.StatusCreated, message: "201 Created"},
		{code: http.StatusAccepted, message: "202 Accepted"},
		{code: http.StatusNonAuthoritativeInfo, message: "203 Partial Information"},
		{code: http.StatusNoContent, message: "204 No Response"},
		{code: http.StatusResetContent, message: "205 Reset Content"},
		{code: http.StatusPartialContent, message: "206 Partial Content"},
		{code: http.StatusMultiStatus, message: "207 Multi Status"},
		{code: http.StatusAlreadyReported, message: "208 Already Reported"},
		{code: http.StatusIMUsed, message: "226 IM Used"},
		{code: http.StatusMultipleChoices, message: "300 Multiple Choices"},
		{code: http.StatusMovedPermanently, message: "301 Moved Permanently"},
		{code: http.StatusFound, message: "302 Found"},
		{code: http.StatusSeeOther, message: "303 See Other"},
		{code: http.StatusNotModified, message: "304 Not Modified"},
		{code: http.StatusUseProxy, message: "305 Use Proxy"},
		{code: http.StatusTemporaryRedirect, message: "307 Temporary Redirect"},
		{code: http.StatusPermanentRedirect, message: "308 Permanent Redirect"},
	}

	invalidHTTPCodesList := []httpCodes{
		{code: http.StatusBadRequest, message: "400 Bad Request"},
		{code: http.StatusUnauthorized, message: "401 Unauthorized"},
		{code: http.StatusPaymentRequired, message: "402 Payment Required"},
		{code: http.StatusForbidden, message: "403 Forbidden"},
		{code: http.StatusNotFound, message: "404 Not Found"},
		{code: http.StatusMethodNotAllowed, message: "405 Method Not Allowed"},
		{code: http.StatusNotAcceptable, message: "406 Not Acceptable"},
		{code: http.StatusProxyAuthRequired, message: "407 Proxy Authentication Required"},
		{code: http.StatusRequestTimeout, message: "408 Request Timeout"},
		{code: http.StatusConflict, message: "409 Conflict"},
		{code: http.StatusGone, message: "410 Gone"},
		{code: http.StatusLengthRequired, message: "411 Length Required"},
		{code: http.StatusPreconditionFailed, message: "412 Precondition Failed"},
		{code: http.StatusRequestEntityTooLarge, message: "413 Request Entity Too Large"},
		{code: http.StatusRequestURITooLong, message: "414 Request URI Too Long"},
		{code: http.StatusUnsupportedMediaType, message: "415 Unsupported Media Type"},
		{code: http.StatusRequestedRangeNotSatisfiable, message: "416 Requested Range Not Satisfiable"},
		{code: http.StatusExpectationFailed, message: "417 Expectation Failed"},
		{code: http.StatusTeapot, message: "418 I'm a teapot"},
		{code: http.StatusUnprocessableEntity, message: "422 Unprocessable Entity"},
		{code: http.StatusLocked, message: "423 Locked"},
		{code: http.StatusFailedDependency, message: "424 Failed Dependency"},
		{code: http.StatusUpgradeRequired, message: "426 Upgrade Required"},
		{code: http.StatusPreconditionRequired, message: "428 Precondition Required"},
		{code: http.StatusTooManyRequests, message: "429 Too Many Requests"},
		{code: http.StatusRequestHeaderFieldsTooLarge, message: "431 Request Header Fields Too Large"},
		{code: http.StatusUnavailableForLegalReasons, message: "451 Unavailable For Legal Reasons"},
		{code: http.StatusInternalServerError, message: "500 Internal Server Error"},
		{code: http.StatusNotImplemented, message: "501 Not Implemented"},
		{code: http.StatusBadGateway, message: "502 Bad Gateway"},
		{code: http.StatusServiceUnavailable, message: "503 Service Unavailable"},
		{code: http.StatusGatewayTimeout, message: "504 Gateway Timeout"},
		{code: http.StatusHTTPVersionNotSupported, message: "505 HTTP Version Not Supported"},
	}

	t.Run("should return nil error when respond is 200", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		mockBody := `{"Message": "Hello, test world!"}`

		mockResp := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(mockBody)),
			ContentLength: int64(len(mockBody)),
		}

		err = got.checkHTTPResponse(mockResp)
		assert.NoError(t, err)
	})

	t.Run("should return nil error when respond code >= 200 and < 400", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		for _, httpCode := range validHTTPCodesList {
			mockResp := &http.Response{
				Status:        httpCode.message,
				StatusCode:    httpCode.code,
				Proto:         "HTTP/1.1",
				Body:          io.NopCloser(strings.NewReader(httpCode.message)),
				ContentLength: int64(len(httpCode.message)),
			}

			err = got.checkHTTPResponse(mockResp)
			assert.NoError(t, err)
		}
	})

	t.Run("should return error when respond code < 200 and >= 400", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		for _, httpCode := range invalidHTTPCodesList {
			mockResp := &http.Response{
				Status:        httpCode.message,
				StatusCode:    httpCode.code,
				Proto:         "HTTP/1.1",
				Body:          io.NopCloser(strings.NewReader(httpCode.message)),
				ContentLength: int64(len(httpCode.message)),
			}

			gotErr := got.checkHTTPResponse(mockResp)
			assert.Error(t, gotErr)

			httpErr := new(HTTPResponseError)
			if errors.As(gotErr, &httpErr) {
				assert.Equal(t, httpCode.code, httpErr.StatusCode)
			}
		}
	})

	t.Run("should return error when response and body has error", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		mockResp := &http.Response{
			Status:        "400 Bad Request",
			StatusCode:    http.StatusBadRequest,
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(mockErrReader(0)),
			ContentLength: int64(0),
		}

		gotErr := got.checkHTTPResponse(mockResp)
		assert.Error(t, gotErr)
	})
}

func TestCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	CreateUserResponseFile := "testdata/CreateUserResponse_Active.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, CreateUserResponseFile)

		httpResp := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1@mail.com",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "test",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "9067729b3d-94f1e0b3-c394-48d5-8ab1-2c122a167074", got.ID)
		assert.Equal(t, "701984", got.ExternalID)
		assert.Equal(t, "bjensen", got.UserName)
		assert.Equal(t, "Barbara", got.Name.GivenName)
		assert.Equal(t, "Jensen", got.Name.FamilyName)
		assert.Equal(t, "Babs Jensen", got.DisplayName)
		assert.Equal(t, "bjensen@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})

	t.Run("should return an error when usr is nil", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.CreateUser(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrCreateUserRequestEmpty)
	})

	t.Run("should return an error when usr.UserName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "test",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrUserNameEmpty)
	})

	t.Run("should return an error when usr.DisplayName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "test",
			},
			DisplayName: "",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrDisplayNameEmpty)
	})

	t.Run("should return an error when usr.GivenName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrGivenNameEmpty)
	})

	t.Run("should return an error when usr.FamilyName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrFamilyNameEmpty)
	})

	t.Run("should return an error when usr.Emails == 0", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Active:      true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrEmailsEmpty)
	})

	t.Run("should return an error when usr.Emails has no Primary", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: false,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrPrimaryEmailEmpty)
	})

	t.Run("should return an error when usr.Emails > 1", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
				{
					Value:   "alias+user.1@mail.com",
					Type:    "work",
					Primary: false,
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrEmailsTooMany)
	})

	t.Run("should return an error when usr.Addresses > 1", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Addresses: []Address{
				{
					StreetAddress: "street 1",
					Locality:      "locality 1",
					Region:        "region 1",
					PostalCode:    "postal code 1",
					Country:       "country 1",
				},
				{
					StreetAddress: "street 2",
					Locality:      "locality 2",
					Region:        "region 2",
					PostalCode:    "postal code 2",
					Country:       "country 2",
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrAddressesTooMany)
	})

	t.Run("should return an error when usr.PhoneNumbers > 1", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			PhoneNumbers: []PhoneNumber{
				{
					Value: "123456789",
					Type:  "work",
				},
				{
					Value: "987654321",
					Type:  "home",
				},
			},
			Active: true,
		}

		got, err := service.CreateUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrPhoneNumbersTooMany)
	})
}

func TestCreateOrGetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		CreateUserResponseFile := "testdata/CreateUserResponse_Active.json"
		jsonResp := ReadJSONFileAsString(t, CreateUserResponseFile)

		httpResp := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1@mail.com",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "test",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateOrGetUser(context.Background(), usrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "9067729b3d-94f1e0b3-c394-48d5-8ab1-2c122a167074", got.ID)
		assert.Equal(t, "701984", got.ExternalID)
		assert.Equal(t, "bjensen", got.UserName)
		assert.Equal(t, "Barbara", got.Name.GivenName)
		assert.Equal(t, "Jensen", got.Name.FamilyName)
		assert.Equal(t, "Babs Jensen", got.DisplayName)
		assert.Equal(t, "bjensen@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})

	t.Run("should return a 409 response and execute the get user when not field changed", func(t *testing.T) {
		CreateUserResponseConflictFile := "testdata/CreateUserResponse_Conflict.json"
		ListUserResponseFile := "testdata/ListUserResponse.json"
		PutUserResponseFile := "testdata/PutUserResponse.json"

		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonRespConflict := ReadJSONFileAsString(t, CreateUserResponseConflictFile)
		jsonRespOK := ReadJSONFileAsString(t, ListUserResponseFile)
		jsonPutUserRespOK := ReadJSONFileAsString(t, PutUserResponseFile)

		httpRespConflict := &http.Response{
			Status:     "409 Conflict",
			StatusCode: http.StatusConflict,
			Header: http.Header{
				"Date":             []string{"Fri, 18 Mar 2022 10:57:08 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"81abca44-4ee3-47fa-b4d9-729908ef1dd9"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonRespConflict)),
			ContentLength: int64(len(jsonRespConflict)),
		}

		httpRespOK := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonRespOK)),
			ContentLength: int64(len(jsonRespOK)),
		}

		PutUserRespOK := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonPutUserRespOK)),
			ContentLength: int64(len(jsonPutUserRespOK)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpRespConflict, nil).Times(1)
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpRespOK, nil).Times(1)
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(PutUserRespOK, nil).Times(1)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447",
			ExternalID: "702135",
			UserName:   "mjack",
			Name: &Name{
				FamilyName: "Jackson",
				GivenName:  "Mark",
			},
			DisplayName: "Mark Jackson",
			Emails: []Email{
				{
					Value:   "mjack@example.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateOrGetUser(context.Background(), usrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447", got.ID)
		assert.Equal(t, "702135", got.ExternalID)
		assert.Equal(t, "mjack", got.UserName)
		assert.Equal(t, "Mark", got.Name.GivenName)
		assert.Equal(t, "Jackson", got.Name.FamilyName)
		assert.Equal(t, "Mark Jackson", got.DisplayName)
		assert.Equal(t, "mjack@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})

	t.Run("should return a 409 response and execute the get user when fields changed", func(t *testing.T) {
		CreateUserResponseConflictFile := "testdata/CreateUserResponse_Conflict.json"
		ListUserResponseFile := "testdata/ListUserResponse_fields_changes.json"
		PutUserResponseFile := "testdata/PutUserResponse.json"

		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonRespConflict := ReadJSONFileAsString(t, CreateUserResponseConflictFile)
		jsonListRespOK := ReadJSONFileAsString(t, ListUserResponseFile)
		jsonPutUserRespOK := ReadJSONFileAsString(t, PutUserResponseFile)

		httpRespConflict := &http.Response{
			Status:     "409 Conflict",
			StatusCode: http.StatusConflict,
			Header: http.Header{
				"Date":             []string{"Fri, 18 Mar 2022 10:57:08 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"81abca44-4ee3-47fa-b4d9-729908ef1dd9"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonRespConflict)),
			ContentLength: int64(len(jsonRespConflict)),
		}

		ListRespOK := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonListRespOK)),
			ContentLength: int64(len(jsonListRespOK)),
		}

		PutUserRespOK := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonPutUserRespOK)),
			ContentLength: int64(len(jsonPutUserRespOK)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpRespConflict, nil).Times(1)
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(ListRespOK, nil).Times(1)
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(PutUserRespOK, nil).Times(1)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447",
			ExternalID: "702135",
			UserName:   "mjack",
			Name: &Name{
				FamilyName: "Jackson",
				GivenName:  "Mark",
			},
			DisplayName: "Mark Jackson",
			Emails: []Email{
				{
					Value:   "mjack@example.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.CreateOrGetUser(context.Background(), usrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447", got.ID)
		assert.Equal(t, "702135", got.ExternalID)
		assert.Equal(t, "mjack", got.UserName)
		assert.Equal(t, "Mark", got.Name.GivenName)
		assert.Equal(t, "Jackson", got.Name.FamilyName)
		assert.Equal(t, "Mark Jackson", got.DisplayName)
		assert.Equal(t, "mjack@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})
}

func TestDeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		userID := "1"
		reqURL.Path = path.Join(reqURL.Path, fmt.Sprintf("/Users/%s", userID))

		httpReq, err := http.NewRequestWithContext(context.Background(), "DELETE", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "204 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("")),
			ContentLength: int64(len("")),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		err = service.DeleteUser(context.Background(), userID)
		assert.NoError(t, err)
	})
}

func TestGetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	GetUserResponseFile := "testdata/GetUserResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, GetUserResponseFile)

		userID := "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447"
		reqURL.Path = path.Join(reqURL.Path, fmt.Sprintf("/Users/%s", userID))

		httpReq, err := http.NewRequestWithContext(context.Background(), "GET", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.GetUser(context.Background(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, userID, got.ID)
		assert.Equal(t, "702135", got.ExternalID)
		assert.Equal(t, "mjack", got.UserName)
		assert.Equal(t, "Jackson", got.Name.GivenName)
		assert.Equal(t, "Mark", got.Name.FamilyName)
		assert.Equal(t, "mjack", got.DisplayName)
		assert.Equal(t, "mjack@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, false, got.Active)
	})
}

func TestGetUserByUserName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	ListUserResponseFile := "testdata/ListUserResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, ListUserResponseFile)

		userID := "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447"
		userName := "mjack"

		reqURL.Path = path.Join(reqURL.Path, "/Users")

		filter := fmt.Sprintf("userName eq %q", userName)
		q := reqURL.Query()
		q.Add("filter", filter)
		reqURL.RawQuery = q.Encode()

		httpReq, err := http.NewRequestWithContext(context.Background(), "GET", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.GetUserByUserName(context.Background(), userName)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, userID, got.ID)
		assert.Equal(t, "702135", got.ExternalID)
		assert.Equal(t, userName, got.UserName)
		assert.Equal(t, "Jackson", got.Name.GivenName)
		assert.Equal(t, "Mark", got.Name.FamilyName)
		assert.Equal(t, "mjack", got.DisplayName)
		assert.Equal(t, "mjack@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, false, got.Active)
	})
}

func TestListUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	ListUserResponseFile := "testdata/ListUserResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, ListUserResponseFile)

		userID := "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447"
		filter := "userName eq mjack"

		reqURL.Path = path.Join(reqURL.Path, "/Users")

		q := reqURL.Query()
		q.Add("filter", filter)
		reqURL.RawQuery = q.Encode()

		httpReq, err := http.NewRequestWithContext(context.Background(), "GET", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.ListUsers(context.Background(), filter)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "urn:ietf:params:scim:api:messages:2.0:ListResponse", got.Schemas[0])
		assert.Equal(t, userID, got.Resources[0].ID)
		assert.Equal(t, "702135", got.Resources[0].ExternalID)
		assert.Equal(t, "mjack", got.Resources[0].UserName)
		assert.Equal(t, "Jackson", got.Resources[0].Name.GivenName)
		assert.Equal(t, "Mark", got.Resources[0].Name.FamilyName)
		assert.Equal(t, "mjack", got.Resources[0].DisplayName)
		assert.Equal(t, "mjack@example.com", got.Resources[0].Emails[0].Value)
		assert.Equal(t, "work", got.Resources[0].Emails[0].Type)
		assert.Equal(t, true, got.Resources[0].Emails[0].Primary)
		assert.Equal(t, false, got.Resources[0].Active)
	})
}

func TestPutUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	PutUserResponseFile := "testdata/PutUserResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, PutUserResponseFile)

		httpResp := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pusrr := &PutUserRequest{
			ID:         "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447",
			ExternalID: "702135",
			UserName:   "mjack",
			Name: &Name{
				FamilyName: "Jackson",
				GivenName:  "Mark",
			},
			DisplayName: "Mark Jackson",
			Emails: []Email{
				{
					Value:   "mjack@example.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.PutUser(context.Background(), pusrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447", got.ID)
		assert.Equal(t, "702135", got.ExternalID)
		assert.Equal(t, "mjack", got.UserName)
		assert.Equal(t, "Mark Jackson", got.DisplayName)
		assert.Equal(t, true, got.Active)
	})

	t.Run("should return an error when usr is nil", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.PutUser(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrPutUserRequestEmpty)
	})

	t.Run("should return an error when usr.DisplayName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pusrr := &PutUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "test",
			},
			DisplayName: "",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.PutUser(context.Background(), pusrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrDisplayNameEmpty)
	})

	t.Run("should return an error when usr.GivenName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &PutUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.PutUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrGivenNameEmpty)
	})

	t.Run("should return an error when usr.FamilyName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &PutUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
			},
			Active: true,
		}

		got, err := service.PutUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrFamilyNameEmpty)
	})

	t.Run("should return an error when usr.Emails > 1", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &PutUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1",
			Name: &Name{
				FamilyName: "1",
				GivenName:  "user",
			},
			DisplayName: "user 1",
			Emails: []Email{
				{
					Value:   "user.1@mail.com",
					Type:    "work",
					Primary: true,
				},
				{
					Value:   "alias+user.1@mail.com",
					Type:    "work",
					Primary: false,
				},
			},
			Active: true,
		}

		got, err := service.PutUser(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrEmailsTooMany)
	})
}

func TestPatchUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	PatchUserResponseFile := "testdata/PatchUserResponse.json"

	t.Run("should return an error when usr is nil", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		err = service.PatchUser(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, PatchUserResponseFile)

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pur := &PatchUserRequest{
			User: User{
				ID: "9067729b3d-94f1e0b3-c394-48d5-8ab1-2c122a167074",
			},
			Patch: Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []*Operation{
					{
						OP:    "replace",
						Path:  "active",
						Value: true,
					},
				},
			},
		}

		err = service.PatchUser(context.Background(), pur)
		assert.NoError(t, err)
	})

	t.Run("should return an error when usr.ID is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pur := &PatchUserRequest{
			User: User{},
			Patch: Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []*Operation{
					{
						OP:    "replace",
						Path:  "active",
						Value: true,
					},
				},
			},
		}

		err = service.PatchUser(context.Background(), pur)
		assert.Error(t, err)
	})
}

func TestCreateGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	CreateGroupResponseFile := "testdata/CreateGroupResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, CreateGroupResponseFile)

		httpResp := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Mon, 06 Apr 2020 16:48:19 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateGroupRequest{
			ID:          "9067729b3d-a2cfc8a5-f4ab-4443-9d7d-b32a9013c554",
			DisplayName: "Group Bar",
		}

		got, err := service.CreateGroup(context.Background(), usrr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "9067729b3d-a2cfc8a5-f4ab-4443-9d7d-b32a9013c554", got.ID)
		assert.Equal(t, "Group Bar", got.DisplayName)
	})

	t.Run("should return an error when usr is nil", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.CreateGroup(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrCreateGroupRequestEmpty)
	})

	t.Run("should return an error when group.DisplayName is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateGroupRequest{
			ID:          "9067729b3d-a2cfc8a5-f4ab-4443-9d7d-b32a9013c554",
			DisplayName: "",
		}

		got, err := service.CreateGroup(context.Background(), usrr)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrDisplayNameEmpty)
	})
}

func TestDeleteGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		groupID := "1"
		reqURL.Path = path.Join(reqURL.Path, fmt.Sprintf("/Groups/%s", groupID))

		httpReq, err := http.NewRequestWithContext(context.Background(), "DELETE", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "204 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Mon, 06 Apr 2020 22:21:24 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("")),
			ContentLength: int64(len("")),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		err = service.DeleteGroup(context.Background(), groupID)
		assert.NoError(t, err)
	})
}

func TestGetGroupByDisplayName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	ListGroupsResponseFile := "testdata/ListGroupsResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, ListGroupsResponseFile)

		groupID := "90677c608a-ef9cb2da-d480-422b-9901-451b1bf9e607"
		displayName := "Group Foo"

		reqURL.Path = path.Join(reqURL.Path, "/Groups")

		filter := fmt.Sprintf("displayName eq %q", displayName)
		q := reqURL.Query()
		q.Add("filter", filter)
		reqURL.RawQuery = q.Encode()

		httpReq, err := http.NewRequestWithContext(context.Background(), "GET", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Thu, 23 Jul 2020 00:37:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"e01400a1-0f10-4e90-ba58-ea1766a009d7"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.GetGroupByDisplayName(context.Background(), displayName)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, groupID, got.ID)
		assert.Equal(t, displayName, got.DisplayName)
	})
}

func TestCreateOrGetGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	CreateGroupResponseFile := "testdata/CreateGroupResponse.json"
	CreateGroupResponseConflictFile := "testdata/CreateGroupResponse_Conflict.json"
	ListGroupsResponseFile := "testdata/ListGroupsResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, CreateGroupResponseFile)

		httpResp := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Mon, 06 Apr 2020 16:48:19 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		grpr := &CreateGroupRequest{
			DisplayName: "Group Foo",
			Members:     []*Member{},
		}

		got, err := service.CreateOrGetGroup(context.Background(), grpr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "9067729b3d-a2cfc8a5-f4ab-4443-9d7d-b32a9013c554", got.ID)
		assert.Equal(t, "Group Bar", got.DisplayName)
	})

	t.Run("should return a 409 response and execute the get user", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonRespConflict := ReadJSONFileAsString(t, CreateGroupResponseConflictFile)
		jsonRespOK := ReadJSONFileAsString(t, ListGroupsResponseFile)

		httpRespConflict := &http.Response{
			Status:     "409 Conflict",
			StatusCode: http.StatusConflict,
			Header: http.Header{
				"Date":             []string{"Fri, 18 Mar 2022 10:57:08 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"81abca44-4ee3-47fa-b4d9-729908ef1dd9"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonRespConflict)),
			ContentLength: int64(len(jsonRespConflict)),
		}

		httpRespOK := &http.Response{
			Status:     "201 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 31 Mar 2020 02:36:15 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"abbf9e53-9ecc-46d2-8efe-104a66ff128f"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonRespOK)),
			ContentLength: int64(len(jsonRespOK)),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpRespConflict, nil).Times(1)
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpRespOK, nil).Times(1)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		grpr := &CreateGroupRequest{
			DisplayName: "Group Foo",
			Members:     []*Member{},
		}

		got, err := service.CreateOrGetGroup(context.Background(), grpr)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "90677c608a-ef9cb2da-d480-422b-9901-451b1bf9e607", got.ID)
		assert.Equal(t, "Group Foo", got.DisplayName)
	})
}

func TestPatchGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	t.Run("should return an error when usr is nil", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		err = service.PatchGroup(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		httpResp := &http.Response{
			Status:     "204 OK",
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Date":             []string{"Tue, 07 Apr 2020 23:59:09 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"dad0c91c-1ea8-4b36-9fdb-4f099b59c1c9"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("")),
			ContentLength: int64(len("")),
		}

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pur := &PatchGroupRequest{
			Group: Group{
				ID: "9067729b3d-94f1e0b3-c394-48d5-8ab1-2c122a167074",
			},
			Patch: Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []*Operation{
					{
						OP:    "replace",
						Value: struct{ DisplayName string }{"Group Foo"},
					},
				},
			},
		}

		err = service.PatchGroup(context.Background(), pur)
		assert.NoError(t, err)
	})

	t.Run("should return an error when group.ID is empty", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		pur := &PatchGroupRequest{
			Group: Group{},
			Patch: Patch{
				Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
				Operations: []*Operation{
					{
						OP:    "replace",
						Value: struct{ DisplayName string }{""},
					},
				},
			},
		}

		err = service.PatchGroup(context.Background(), pur)
		assert.Error(t, err)
	})
}

func TestListGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	ListUserResponseFile := "testdata/ListGroupsResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFileAsString(t, ListUserResponseFile)

		groupID := "90677c608a-ef9cb2da-d480-422b-9901-451b1bf9e607"
		filter := "displayName eq \"Group Foo\""

		reqURL.Path = path.Join(reqURL.Path, "/Groups")

		q := reqURL.Query()
		q.Add("filter", filter)
		reqURL.RawQuery = q.Encode()

		httpReq, err := http.NewRequestWithContext(context.Background(), "GET", reqURL.String(), nil)
		assert.NoError(t, err)

		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Bearer MyToken")

		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusNoContent,
			Header: http.Header{
				"Date":             []string{"Wed, 22 Jul 2020 23:06:38 GMT"},
				"Content-Type":     []string{"application/json"},
				"x-amzn-RequestId": []string{"45995b44-02cd-419f-87f4-ff8fa323448d"},
			},
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader(jsonResp)),
			ContentLength: int64(len(jsonResp)),
		}

		mockHTTPClient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPClient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		got, err := service.ListGroups(context.Background(), filter)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, "urn:ietf:params:scim:api:messages:2.0:ListResponse", got.Schemas[0])
		assert.Equal(t, groupID, got.Resources[0].ID)
		assert.Equal(t, "Group Foo", got.Resources[0].DisplayName)
	})
}

// Enhanced test cases for improved coverage

func TestNewSCIMServiceWithHTTPConfig(t *testing.T) {
	t.Run("should create service with custom HTTP config", func(t *testing.T) {
		got, err := NewSCIMServiceWithHTTPConfig("https://testing.com", "MyToken", 10, 5)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should return error when url is empty", func(t *testing.T) {
		got, err := NewSCIMServiceWithHTTPConfig("", "MyToken", 10, 5)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrURLEmpty)
		assert.Nil(t, got)
	})

	t.Run("should return error when token is empty", func(t *testing.T) {
		got, err := NewSCIMServiceWithHTTPConfig("https://testing.com", "", 10, 5)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrBearerTokenEmpty)
		assert.Nil(t, got)
	})
}

func TestSCIMService_ContextCancellation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle context cancellation", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		reqURL, _ := url.Parse("https://testing.com/Users")
		req, _ := service.newRequest(ctx, http.MethodGet, reqURL, nil)

		_, err = service.do(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestSCIMService_RateLimiting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle 429 Too Many Requests", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		// Create a mock response with 429 status
		resp := &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Status:     "429 Too Many Requests",
			Body:       io.NopCloser(strings.NewReader(`{"detail":"Rate limit exceeded"}`)),
		}

		err = service.checkHTTPResponse(resp)
		assert.Error(t, err)

		var httpErr *HTTPResponseError
		assert.True(t, errors.As(err, &httpErr))
		assert.Equal(t, http.StatusTooManyRequests, httpErr.StatusCode)
	})
}

func TestSCIMService_LargeResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle large response bodies", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		// Create a large response body (10KB)
		largeBody := strings.Repeat("x", 10*1024)
		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader(largeBody)),
		}

		err = service.checkHTTPResponse(resp)
		assert.Error(t, err)

		var httpErr *HTTPResponseError
		assert.True(t, errors.As(err, &httpErr))
		assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
		assert.Equal(t, largeBody, httpErr.Message)
	})
}

func TestSCIMService_MalformedJSON(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle malformed JSON in error response", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		malformedJSON := `{"detail":"Error","invalid"}`
		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader(malformedJSON)),
		}

		err = service.checkHTTPResponse(resp)
		assert.Error(t, err)

		var httpErr *HTTPResponseError
		assert.True(t, errors.As(err, &httpErr))
		assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
		assert.Equal(t, malformedJSON, httpErr.Message)
	})
}

func TestSCIMService_StructuredErrorResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should parse structured SCIM error response", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		scimError := `{
			"schemas": ["urn:ietf:params:scim:api:messages:2.0:Error"],
			"scimType": "invalidValue",
			"detail": "Invalid filter syntax",
			"status": "400"
		}`
		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader(scimError)),
		}

		err = service.checkHTTPResponse(resp)
		assert.Error(t, err)

		var httpErr *HTTPResponseError
		assert.True(t, errors.As(err, &httpErr))
		assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
		assert.Equal(t, "invalidValue", httpErr.Code)
		assert.Equal(t, "Invalid filter syntax", httpErr.Message)
	})
}

func TestSCIMService_HTTPClientError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle HTTP client errors", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		reqURL, _ := url.Parse("https://testing.com/Users")
		req, _ := service.newRequest(context.Background(), http.MethodGet, reqURL, nil)

		// Mock HTTP client to return error
		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network error"))

		_, err = service.do(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network error")
	})
}

func TestSCIMService_EmptyResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should handle empty response body", func(t *testing.T) {
		mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
		service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
		assert.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader("")),
		}

		err = service.checkHTTPResponse(resp)
		assert.Error(t, err)

		var httpErr *HTTPResponseError
		assert.True(t, errors.As(err, &httpErr))
		assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
		assert.Equal(t, "", httpErr.Message)
	})
}

func TestConstants(t *testing.T) {
	t.Run("should have correct constant values", func(t *testing.T) {
		assert.Equal(t, "/Users", UsersPath)
		assert.Equal(t, "/Groups", GroupsPath)
		assert.Equal(t, "/ServiceProviderConfig", ServiceProviderConfigPath)
		assert.Equal(t, "application/scim+json", ContentTypeSCIMJSON)
		assert.Equal(t, "application/json", ContentTypeJSON)
		assert.Equal(t, int64(30*1000000000), int64(DefaultTimeout))
		assert.Equal(t, 3, MaxRetries)
	})
}
