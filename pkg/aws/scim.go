package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

// Consume http methods
// implement scim.AWSSCIMProvider interface

// AWS SSO SCIM API
// reference: https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html

var (
	ErrParsingURL          = errors.Errorf("error parsing url")
	ErrEndpointEmpty       = errors.Errorf("endpoint is empty and it is required")
	ErrCallingDo           = errors.Errorf("error calling http client Do method")
	ErrReadingResponseBody = errors.Errorf("error reading response body")
	ErrSendingRequest      = errors.Errorf("error sending request")
	ErrDecodingResponse    = errors.Errorf("error decoding response")
	ErrEncodingRequest     = errors.Errorf("error encoding request")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/aws/scim_mocks.go -source=scim.go HTTPClient

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AWSSCIMProvider struct {
	httpClient  HTTPClient
	url         *url.URL
	UserAgent   string
	bearerToken string
}

func NewSCIMService(httpClient HTTPClient, endpoint string, token string) (*AWSSCIMProvider, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if endpoint == "" {
		return nil, errors.Wrapf(ErrEndpointEmpty, "NewSCIMService")
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(ErrParsingURL, "NewSCIMService -> "+err.Error())
	}

	return &AWSSCIMProvider{
		httpClient:  httpClient,
		url:         url,
		bearerToken: token,
	}, nil
}

func (s *AWSSCIMProvider) newRequest(method string, url string, body interface{}) (*http.Request, error) {
	u, err := s.url.Parse(url)
	if err != nil {
		return nil, errors.Wrapf(ErrParsingURL, "NewRequest -> "+err.Error())
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		err := enc.Encode(body)
		if err != nil {
			return nil, errors.Wrapf(ErrEncodingRequest, "NewRequest -> "+err.Error())
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	if s.UserAgent != "" {
		req.Header.Set("User-Agent", s.UserAgent)
	}
	return req, nil
}

func (s *AWSSCIMProvider) GetURL() *url.URL {
	return s.url
}

func (s *AWSSCIMProvider) checkHTTPResponse(r *http.Response) error {
	if r.StatusCode < http.StatusOK || r.StatusCode >= http.StatusBadRequest {
		var errResp APIErrorResponse
		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			defer r.Body.Close()
			return errors.Wrapf(ErrDecodingResponse, "checkHTTPResponse -> "+err.Error())
		}
		return errors.Errorf("%s: %s", r.Status, utils.ToJSON(errResp))
	}

	return nil
}

func (s *AWSSCIMProvider) request(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(ErrCallingDo, "request -> "+err.Error())
	}

	if err := s.checkHTTPResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AWSSCIMProvider) ListUsers(ctx context.Context, filter string) (*UsersResponse, error) {
	s.url.Path = path.Join(s.url.Path, "/Users")

	if filter != "" {
		q := s.url.Query()
		q.Add("filter", filter)
		s.url.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, s.url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.url.String())
	}

	resp, err := s.request(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ListUsers -> "+err.Error())
	}
	defer resp.Body.Close()

	var response UsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ListUsers -> "+err.Error())
	}

	return &response, nil
}

func (s *AWSSCIMProvider) ListGroups(ctx context.Context, filter string) (*GroupsResponse, error) {
	s.url.Path = path.Join(s.url.Path, "/Groups")

	if filter != "" {
		q := s.url.Query()
		q.Add("filter", filter)
		s.url.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, s.url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.url.String())
	}

	resp, err := s.request(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ListGroups -> "+err.Error())
	}
	defer resp.Body.Close()

	var response GroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ListGroups -> "+err.Error())
	}

	return &response, nil
}

func (s *AWSSCIMProvider) ServiceProviderConfig(ctx context.Context) (*ServiceProviderConfig, error) {
	s.url.Path = path.Join(s.url.Path, "/ServiceProviderConfig")

	req, err := s.newRequest(http.MethodGet, s.url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.url.String())
	}

	resp, err := s.request(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ServiceProviderConfig -> "+err.Error())
	}
	defer resp.Body.Close()

	var response ServiceProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ServiceProviderConfig -> "+err.Error())
	}

	return &response, nil
}
