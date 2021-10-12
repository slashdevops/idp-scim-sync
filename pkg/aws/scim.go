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
	ErrURLEmpty         = errors.Errorf("aws: url may not be empty")
	ErrDisplayNameEmpty = errors.Errorf("aws: display name may not be empty")
	ErrGivenNameEmpty   = errors.Errorf("aws: given name may not be empty")
	ErrFamilyNameEmpty  = errors.Errorf("aws: family name may not be empty")
	ErrEmailsTooMany    = errors.Errorf("aws: emails may not be more than 1")
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

func NewSCIMService(httpClient HTTPClient, urlStr string, token string) (*AWSSCIMProvider, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if urlStr == "" {
		return nil, ErrURLEmpty
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %v", err)
	}

	return &AWSSCIMProvider{
		httpClient:  httpClient,
		url:         url,
		bearerToken: token,
	}, nil
}

func (s *AWSSCIMProvider) newRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	u, err := s.url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %v", err)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		err := enc.Encode(body)
		if err != nil {
			return nil, fmt.Errorf("aws: error encoding request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request: %v", err)
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
			return fmt.Errorf("aws: error decoding response body: %v", err)
		}
		return fmt.Errorf("aws: error, code %s, response %s", r.Status, utils.ToJSON(errResp))
	}

	return nil
}

func (s *AWSSCIMProvider) do(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request: %v", err)
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
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}
	defer resp.Body.Close()

	var response UsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) ListGroups(ctx context.Context, filter string) (*ListsGroupsResponse, error) {
	s.url.Path = path.Join(s.url.Path, "/Groups")

	if filter != "" {
		q := s.url.Query()
		q.Add("filter", filter)
		s.url.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, s.url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}
	defer resp.Body.Close()

	var response ListsGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) ServiceProviderConfig(ctx context.Context) (*ServiceProviderConfig, error) {
	s.url.Path = path.Join(s.url.Path, "/ServiceProviderConfig")

	req, err := s.newRequest(http.MethodGet, s.url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %v", http.MethodGet, s.url.String(), err)
	}
	defer resp.Body.Close()

	var response ServiceProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) CreateGroup(ctx context.Context, g *CreateGroupRequest) (*CreateGroupResponse, error) {
	if g == nil {
		return nil, fmt.Errorf("aws: error creating group, group is nil")
	}
	if g.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}

	s.url.Path = path.Join(s.url.Path, "/Groups")

	req, err := s.newRequest(http.MethodPost, s.url.String(), *g)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %v", http.MethodPost, s.url.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %v", http.MethodPost, s.url.String(), err)
	}
	defer resp.Body.Close()

	var response CreateGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}

// CreateUser creates a new user in the SCIM provider.
// references:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/createuser.html
func (s *AWSSCIMProvider) CreateUser(ctx context.Context, u *CreateUserRequest) (*CreateUserResponse, error) {
	if u == nil {
		return nil, fmt.Errorf("aws: error creating user, user is nil")
	}

	// Check constraints in the reference document
	if u.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}
	if u.Name.GivenName == "" {
		return nil, ErrGivenNameEmpty
	}
	if u.Name.FamilyName == "" {
		return nil, ErrFamilyNameEmpty
	}

	if len(u.Emails) > 1 {
		return nil, ErrEmailsTooMany
	}
	u.Emails[0].Primary = true

	s.url.Path = path.Join(s.url.Path, "/Users")

	req, err := s.newRequest(http.MethodPost, s.url.String(), *u)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %v", http.MethodPost, s.url.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %v", http.MethodPost, s.url.String(), err)
	}
	defer resp.Body.Close()

	var response CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}
