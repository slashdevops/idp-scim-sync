package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

// Consume HTTPClient interface

var (
	ErrParsingURL          = errors.Errorf("error parsing url")
	ErrEndpointEmpty       = errors.Errorf("endpoint is empty and it is required")
	ErrCallingDo           = errors.Errorf("error calling http client Do method")
	ErrReadingResponseBody = errors.Errorf("error reading response body")
	ErrSendingRequest      = errors.Errorf("error sending request")
	ErrDecodingResponse    = errors.Errorf("error decoding response")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/aws/scim_mocks.go -source=scim.go HTTPClient

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AWSSCIMProvider struct {
	httpClient  HTTPClient
	endpointURL *url.URL
	bearerToken string
}

func NewSCIMService(http HTTPClient, endpoint string, token string) (*AWSSCIMProvider, error) {
	if endpoint == "" {
		return nil, errors.Wrapf(ErrEndpointEmpty, "NewSCIMService")
	}

	scimURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(ErrParsingURL, "NewSCIMService -> "+err.Error())
	}

	return &AWSSCIMProvider{
		httpClient:  http,
		endpointURL: scimURL,
		bearerToken: token,
	}, nil
}

func (s *AWSSCIMProvider) EndpointURL() *url.URL {
	return s.endpointURL
}

func (s *AWSSCIMProvider) sendRequest(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "idp-scim-sync/1.0") // TODO: add right user agent
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(ErrCallingDo, "sendRequest -> "+err.Error())
	}

	return resp, nil
}

func (s *AWSSCIMProvider) ListUsers(ctx context.Context, filter string) (*UsersResponse, error) {
	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/Users")
	var uResponse UsersResponse

	if filter != "" {
		q := s.endpointURL.Query()
		q.Add("filter", filter)
		s.endpointURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.endpointURL.String())
	}

	resp, err := s.sendRequest(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ListUsers -> "+err.Error())
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&uResponse); err == nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ListUsers -> "+err.Error())
	}

	return &uResponse, nil
}

func (s *AWSSCIMProvider) ListGroups(ctx context.Context, filter string) (*GroupsResponse, error) {
	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/Groups")
	var gResponse GroupsResponse

	if filter != "" {
		q := s.endpointURL.Query()
		q.Add("filter", filter)
		s.endpointURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.endpointURL.String())
	}

	resp, err := s.sendRequest(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ListGroups -> "+err.Error())
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&gResponse); err != nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ListGroups -> "+err.Error())
	}

	return &gResponse, nil
}

func (s *AWSSCIMProvider) ServiceProviderConfig(ctx context.Context) (*ServiceProviderConfig, error) {
	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/ServiceProviderConfig")

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: http method -> %s, endpoint -> %v", err, http.MethodGet, s.endpointURL.String())
	}

	resp, err := s.sendRequest(ctx, req, nil)
	if err != nil {
		return nil, errors.Wrapf(ErrSendingRequest, "ServiceProviderConfig -> "+err.Error())
	}
	defer resp.Body.Close()

	var config ServiceProviderConfig
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&config); err != nil {
		return nil, errors.Wrapf(ErrDecodingResponse, "ServiceProviderConfig -> "+err.Error())
	}

	return &config, nil
}
