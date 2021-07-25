package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type SCIMService interface {
	ListUsers(filter string) (*UsersResponse, error)
	ListGroups(filter string) (*GroupsResponse, error)
	ServiceProviderConfig() (*ServiceProviderConfig, error)
}

type scim struct {
	ctx         *context.Context
	httpClient  *http.Client
	endpointURL *url.URL
	bearerToken string
}

func NewSCIMService(ctx *context.Context, http *http.Client, endpoint string, token string) (SCIMService, error) {

	scimURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &scim{
		ctx:         ctx,
		httpClient:  http,
		endpointURL: scimURL,
		bearerToken: token,
	}, nil
}

func (s *scim) EndpointURL() *url.URL {
	return s.endpointURL
}

func (s *scim) sendRequest(req *http.Request, body interface{}) (resp *http.Response, err error) {

	req = req.WithContext(*s.ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "aws-sso-gws-sync/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err = s.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("unknown error, response body: %s status code: %d", string(bodyBytes), resp.StatusCode)
	}

	return
}

func (s *scim) ListUsers(filter string) (*UsersResponse, error) {

	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/Users")
	var uResponse UsersResponse

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.sendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&uResponse); err == nil {
		return nil, fmt.Errorf("error decoding response, error: %s", err)
	}

	return &uResponse, nil
}

func (s *scim) ListGroups(filter string) (*GroupsResponse, error) {

	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/Groups")
	var gResponse GroupsResponse

	if filter != "" {
		q := s.endpointURL.Query()
		q.Add("filter", filter)
		s.endpointURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.sendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&gResponse); err != nil {
		return nil, fmt.Errorf("error decoding response, error: %s", err)
	}

	return &gResponse, nil
}

func (s *scim) ServiceProviderConfig() (*ServiceProviderConfig, error) {

	s.endpointURL.Path = path.Join(s.endpointURL.Path, "/ServiceProviderConfig")

	req, err := http.NewRequest(http.MethodGet, s.endpointURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.sendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("unknown error, response body: %s status code: %d", string(bodyBytes), resp.StatusCode)
	}

	var config ServiceProviderConfig
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding response, error: %s", err)
	}

	return &config, nil
}
