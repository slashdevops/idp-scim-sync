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

	ErrCreateGroupRequestEmpty = errors.Errorf("aws: create group request may not be empty")
	ErrCreateUserRequestEmpty  = errors.Errorf("aws: create user request may not be empty")

	ErrPatchGroupRequestEmpty = errors.Errorf("aws: patch group request may not be empty")
	ErrGroupIDEmpty           = errors.Errorf("aws: group id may not be empty")
	ErrUserIDEmpty            = errors.Errorf("aws: user id may not be empty")

	ErrPatchUserRequestEmpty = errors.Errorf("aws: patch user request may not be empty")
	ErrPutUserRequestEmpty   = errors.Errorf("aws: put user request may not be empty")
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
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	return &AWSSCIMProvider{
		httpClient:  httpClient,
		url:         url,
		bearerToken: token,
	}, nil
}

func (s *AWSSCIMProvider) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		err := enc.Encode(body)
		if err != nil {
			return nil, fmt.Errorf("aws: error encoding request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request: %w", err)
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

func (s *AWSSCIMProvider) checkHTTPResponse(r *http.Response) error {
	if r.StatusCode < http.StatusOK || r.StatusCode >= http.StatusBadRequest {
		var errResp APIErrorResponse
		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			defer r.Body.Close()
			return fmt.Errorf("aws: error decoding response body: %w", err)
		}
		return fmt.Errorf("aws: error, code %s, response %s", r.Status, utils.ToJSON(errResp))
	}

	return nil
}

// do sends an HTTP request and returns an HTTP response, following policy (e.g. redirects, cookies, auth) as configured on the client.
func (s *AWSSCIMProvider) do(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request: %w", err)
	}

	if err := s.checkHTTPResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateUser creates a new user in the AWS SSO Using the API.
// references:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/createuser.html
func (s *AWSSCIMProvider) CreateUser(ctx context.Context, usr *CreateUserRequest) (*CreateUserResponse, error) {
	if usr == nil {
		return nil, ErrCreateUserRequestEmpty
	}
	// Check constraints in the reference document
	if usr.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}
	if usr.Name.GivenName == "" {
		return nil, ErrGivenNameEmpty
	}
	if usr.Name.FamilyName == "" {
		return nil, ErrFamilyNameEmpty
	}

	if len(usr.Emails) > 1 {
		return nil, ErrEmailsTooMany
	}
	usr.Emails[0].Primary = true

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Users")

	req, err := s.newRequest(http.MethodPost, reqUrl, *usr)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %w", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrUserIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Users/%s", id))

	req, err := s.newRequest(http.MethodPatch, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// ListUsers returns a list of users from the AWS SSO Using the API
func (s *AWSSCIMProvider) ListUsers(ctx context.Context, filter string) (*ListUsersResponse, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Users")

	if filter != "" {
		q := reqUrl.Query()
		q.Add("filter", filter)
		reqUrl.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ListUsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %w", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) PatchUser(ctx context.Context, pur *PatchUserRequest) error {
	if pur == nil {
		return ErrPatchUserRequestEmpty
	}
	if pur.User.ID == "" {
		return ErrUserIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Groups/%s", pur.User.ID))

	req, err := s.newRequest(http.MethodPatch, reqUrl, pur.Patch)
	if err != nil {
		return fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

func (s *AWSSCIMProvider) PutUser(ctx context.Context, usr *PutUserRequest) (*PutUserResponse, error) {
	if usr == nil {
		return nil, ErrPutUserRequestEmpty
	}
	// Check constraints in the reference document
	if usr.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}
	if usr.Name.GivenName == "" {
		return nil, ErrGivenNameEmpty
	}
	if usr.Name.FamilyName == "" {
		return nil, ErrFamilyNameEmpty
	}

	if len(usr.Emails) > 1 {
		return nil, ErrEmailsTooMany
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Users/%s", usr.ID))

	req, err := s.newRequest(http.MethodPut, reqUrl, *usr)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response PutUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %w", err)
	}

	return &response, nil
}

// ListGroups returns a list of groups from the AWS SSO Using the API
func (s *AWSSCIMProvider) ListGroups(ctx context.Context, filter string) (*ListGroupsResponse, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Groups")

	if filter != "" {
		q := reqUrl.Query()
		q.Add("filter", filter)
		reqUrl.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ListGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %w", err)
	}

	return &response, nil
}

// CreateGroup creates a new group in the AWS SSO Using the API
// reference:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/creategroup.html
func (s *AWSSCIMProvider) CreateGroup(ctx context.Context, g *CreateGroupRequest) (*CreateGroupResponse, error) {
	if g == nil {
		return nil, ErrCreateGroupRequestEmpty
	}
	if g.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Groups")

	req, err := s.newRequest(http.MethodPost, reqUrl, *g)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response CreateGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %v", err)
	}

	return &response, nil
}

func (s *AWSSCIMProvider) DeleteGroup(ctx context.Context, id string) error {
	if id == "" {
		return ErrGroupIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Groups/%s", id))

	req, err := s.newRequest(http.MethodPatch, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

func (s *AWSSCIMProvider) PatchGroup(ctx context.Context, pgr *PatchGroupRequest) error {
	if pgr == nil {
		return ErrPatchGroupRequestEmpty
	}
	if pgr.Group.ID == "" {
		return ErrGroupIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Groups/%s", pgr.Group.ID))

	req, err := s.newRequest(http.MethodPatch, reqUrl, pgr.Patch)
	if err != nil {
		return fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// ServiceProviderConfig returns additional information about the AWS SSO SCIM implementation
// references:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/serviceproviderconfig.html
func (s *AWSSCIMProvider) ServiceProviderConfig(ctx context.Context) (*ServiceProviderConfig, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/ServiceProviderConfig")

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ServiceProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws: error decoding response body: %w", err)
	}

	return &response, nil
}
