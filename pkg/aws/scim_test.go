package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
	"github.com/stretchr/testify/assert"
)

func ReadJSONFIleAsString(t *testing.T, fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)

	return string(bytes)
}

func TestNewSCIMService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should return AWSSCIMProvider", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should return AWSSCIMProvider when httpClient is nil", func(t *testing.T) {
		got, err := NewSCIMService(nil, "https://testing.com", "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should return error when url is bad formed", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "https://%%testing.com", "MyToken")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when the url is empty ", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, "", "MyToken")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrURLEmpty)
		assert.Nil(t, got)
	})
}

func TestDo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"

	t.Run("should return error when error come from request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("test error"))

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, got)

		req := httptest.NewRequest(http.MethodGet, endpoint, nil)

		resp, err := got.do(context.Background(), req)
		assert.Error(t, err)

		assert.Nil(t, resp)
	})

	t.Run("should return valid response", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		mockResp := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "HTTP/1.1",
			Body:          io.NopCloser(strings.NewReader("Hello, test world!")),
			ContentLength: int64(len("Hello, test world!")),
		}

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

		got, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
}

func TestCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	CreateUserResponseFile := "testdata/CreateUserResponse_Active.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, CreateUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1@mail.com",
			Name: Name{
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

		assert.Equal(t, "1", got.ID)
		assert.Equal(t, "1", got.ExternalID)
		assert.Equal(t, "user.1@mail.com", got.UserName)
		assert.Equal(t, "user", got.Name.GivenName)
		assert.Equal(t, "1", got.Name.FamilyName)
		assert.Equal(t, "user 1", got.DisplayName)
		assert.Equal(t, "user.1@mail.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})
}

func TestCreateOrGetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	CreateUserResponseFile := "testdata/CreateUserResponse_Active.json"
	CreateUserResponseConflictFile := "testdata/CreateUserResponse_Conflict.json"
	ListUserResponseFile := "testdata/ListUserResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, CreateUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "1",
			ExternalID: "1",
			UserName:   "user.1@mail.com",
			Name: Name{
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

		assert.Equal(t, "1", got.ID)
		assert.Equal(t, "1", got.ExternalID)
		assert.Equal(t, "user.1@mail.com", got.UserName)
		assert.Equal(t, "user", got.Name.GivenName)
		assert.Equal(t, "1", got.Name.FamilyName)
		assert.Equal(t, "user 1", got.DisplayName)
		assert.Equal(t, "user.1@mail.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, true, got.Active)
	})

	t.Run("should return a 409 response and execute the get user", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonRespConflict := ReadJSONFIleAsString(t, CreateUserResponseConflictFile)
		jsonRespOK := ReadJSONFIleAsString(t, ListUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpRespConflict, nil).Times(1)
		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpRespOK, nil).Times(1)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
		assert.NoError(t, err)
		assert.NotNil(t, service)

		usrr := &CreateUserRequest{
			ID:         "90677c608a-7afcdc23-0bd4-4fb7-b2ff-10ccffdff447",
			ExternalID: "702135",
			UserName:   "mjack",
			Name: Name{
				FamilyName: "Mark",
				GivenName:  "Jackson",
			},
			DisplayName: "mjack",
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
		assert.Equal(t, "Jackson", got.Name.GivenName)
		assert.Equal(t, "Mark", got.Name.FamilyName)
		assert.Equal(t, "mjack", got.DisplayName)
		assert.Equal(t, "mjack@example.com", got.Emails[0].Value)
		assert.Equal(t, "work", got.Emails[0].Type)
		assert.Equal(t, true, got.Emails[0].Primary)
		assert.Equal(t, false, got.Active)
	})
}

func TestDeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)

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

		mockHTTPCLient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, GetUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, ListUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, ListUserResponseFile)

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

		mockHTTPCLient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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

func TestGetGroupByDisplayName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	endpoint := "https://testing.com"
	reqURL, err := url.Parse(endpoint)
	assert.NoError(t, err)

	ListGroupsResponseFile := "testdata/ListGroupsResponse.json"

	t.Run("should return a valid response with a valid request", func(t *testing.T) {
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, ListGroupsResponseFile)

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

		mockHTTPCLient.EXPECT().Do(httpReq).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonResp := ReadJSONFIleAsString(t, CreateGroupResponseFile)

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

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpResp, nil)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
		mockHTTPCLient := mocks.NewMockHTTPClient(mockCtrl)
		jsonRespConflict := ReadJSONFIleAsString(t, CreateGroupResponseConflictFile)
		jsonRespOK := ReadJSONFIleAsString(t, ListGroupsResponseFile)

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

		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpRespConflict, nil).Times(1)
		mockHTTPCLient.EXPECT().Do(gomock.Any()).Return(httpRespOK, nil).Times(1)

		service, err := NewSCIMService(mockHTTPCLient, endpoint, "MyToken")
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
