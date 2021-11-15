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
	"strings"

	"github.com/pkg/errors"
)

// Consume http methods
// implement scim.AWSSCIMProvider interface

// AWS SSO SCIM API
// reference: https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html

var (
	// ErrURLEmpty is returned when the URL is empty.
	ErrURLEmpty = errors.Errorf("aws: url may not be empty")

	// ErrDisplayNameEmpty is returned when the display name is empty.
	ErrDisplayNameEmpty = errors.Errorf("aws: display name may not be empty")

	// ErrGivenNameEmpty is returned when the given name is empty.
	ErrGivenNameEmpty = errors.Errorf("aws: given name may not be empty")

	// ErrFamilyNameEmpty is returned when the family name is empty.
	ErrFamilyNameEmpty = errors.Errorf("aws: family name may not be empty")

	// ErrEmailsTooMany is returned when the emails has more than one entity.
	ErrEmailsTooMany = errors.Errorf("aws: emails may not be more than 1")

	// ErrCreateGroupRequestEmpty is returned when the create group request is empty.
	ErrCreateGroupRequestEmpty = errors.Errorf("aws: create group request may not be empty")

	// ErrCreateUserRequestEmpty is returned when the create user request is empty.
	ErrCreateUserRequestEmpty = errors.Errorf("aws: create user request may not be empty")

	// ErrPatchGroupRequestEmpty is returned when the patch group request is empty.
	ErrPatchGroupRequestEmpty = errors.Errorf("aws: patch group request may not be empty")

	// ErrGroupIDEmpty is returned when the group id is empty.
	ErrGroupIDEmpty = errors.Errorf("aws: group id may not be empty")

	// ErrUserIDEmpty is returned when the user id is empty.
	ErrUserIDEmpty = errors.Errorf("aws: user id may not be empty")

	// ErrPatchUserRequestEmpty is returned when the patch user request is empty.
	ErrPatchUserRequestEmpty = errors.Errorf("aws: patch user request may not be empty")

	// ErrPutUserRequestEmpty is returned when the put user request is empty.
	ErrPutUserRequestEmpty = errors.Errorf("aws: put user request may not be empty")

	// ErrUserNameEmpty is returned when the user name is empty.
	ErrUserNameEmpty = errors.Errorf("aws: user name may not be empty")

	// ErrUserUserNameEmpty is returned when the userName is empty.
	ErrUserUserNameEmpty = errors.Errorf("aws: userName may not be empty")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/aws/scim_mocks.go -source=scim.go HTTPClient

// HTTPClient is an interface for sending HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AWSSCIMService is an AWS SCIM Service.
type AWSSCIMService struct {
	httpClient  HTTPClient
	url         *url.URL
	UserAgent   string
	bearerToken string
}

// NewSCIMService creates a new AWS SCIM Service.
func NewSCIMService(httpClient HTTPClient, urlStr string, token string) (*AWSSCIMService, error) {
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

	return &AWSSCIMService{
		httpClient:  httpClient,
		url:         url,
		bearerToken: token,
	}, nil
}

// newRequest creates an http.Request with the given method, URL, and (optionally) body.
func (s *AWSSCIMService) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
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
		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Type", "application/scim+json")
	}

	req.Header.Set("Accept", "application/json")

	if s.UserAgent != "" {
		req.Header.Set("User-Agent", s.UserAgent)
	}
	return req, nil
}

// do sends an HTTP request and returns an HTTP response, following policy (e.g. redirects, cookies, auth) as configured on the client.
func (s *AWSSCIMService) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	// Set bearer token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("aws do: error sending request: %w", err)
	}

	if err := s.checkHTTPResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// checkHTTPResponse checks the status code of the HTTP response.
func (s *AWSSCIMService) checkHTTPResponse(r *http.Response) error {
	if r.StatusCode < http.StatusOK || r.StatusCode >= http.StatusBadRequest {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		for {
			var errResp APIErrorResponse
			if err := dec.Decode(&errResp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("aws checkHTTPResponse: error decoding error response: %w", err)
			}

		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("aws checkHTTPResponse: error reading response body: %w", err)
		}
		return fmt.Errorf("aws checkHTTPResponse: error code: '%s', body: '%s'", r.Status, string(b))
	}

	return nil
}

// CreateUser creates a new user in the AWS SSO Using the API.
// references:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/createuser.html
func (s *AWSSCIMService) CreateUser(ctx context.Context, usr *CreateUserRequest) (*CreateUserResponse, error) {
	if usr == nil {
		return nil, ErrCreateUserRequestEmpty
	}
	// Check constraints in the reference document
	if usr.UserName == "" {
		return nil, ErrUserNameEmpty
	}
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
		return nil, fmt.Errorf("aws CreateUser: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Users")

	// log.WithFields(log.Fields{
	// 	"group": usr.ID,
	// }).Debugf("creating user, payload: %s", utils.ToJSON(usr))

	req, err := s.newRequest(http.MethodPost, reqUrl, *usr)
	if err != nil {
		return nil, fmt.Errorf("aws CreateUser: error creating request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws CreateUser: error sending request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws CreateUser: error decoding response body: %w", err)
	}

	return &response, nil
}

// DeleteUser deletes a user in the AWS SSO Using the API.
func (s *AWSSCIMService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrUserIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Users/%s", id))

	req, err := s.newRequest(http.MethodDelete, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("aws: error creating request, http method: %s, url: %v, error: %w", http.MethodDelete, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodDelete, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// GetUser returns an user from the AWS SSO Using the API
func (s *AWSSCIMService) GetUserByUserName(ctx context.Context, userName string) (*GetUserResponse, error) {
	if userName == "" {
		return nil, ErrUserUserNameEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws GetUserByUserName: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Users")

	filter := fmt.Sprintf("userName eq \"%s\"", userName)
	q := reqUrl.Query()
	q.Add("filter", filter)
	reqUrl.RawQuery = q.Encode()

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws GetUserByUserName: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws GetUserByUserName: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var lur ListUsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&lur); err != nil {
		return nil, fmt.Errorf("aws GetUserByUserName: error decoding response body: %w", err)
	}

	var response GetUserResponse

	dataJSON := fmt.Sprintf("%v", lur.Resources[0])
	if err != nil {
		return nil, fmt.Errorf("aws GetUserByUserName: error decoding response body: %w", err)
	}

	data := strings.NewReader(dataJSON)
	if len(lur.Resources) > 0 {
		if err = json.NewDecoder(data).Decode(&response); err != nil {
			return nil, fmt.Errorf("aws GetUserByUserName: error decoding response body: %w", err)
		}
	}

	return &response, nil
}

// GetUser returns an user from the AWS SSO Using the API
func (s *AWSSCIMService) GetUser(ctx context.Context, userID string) (*GetUserResponse, error) {
	if userID == "" {
		return nil, ErrUserIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws GetUser: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Users/%s", userID))

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws GetUser: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws GetUser: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response GetUserResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws GetUser: error decoding response body: %w", err)
	}

	return &response, nil
}

// ListUsers returns a list of users from the AWS SSO Using the API
func (s *AWSSCIMService) ListUsers(ctx context.Context, filter string) (*ListUsersResponse, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws ListUsers: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Users")

	if filter != "" {
		q := reqUrl.Query()
		q.Add("filter", filter)
		reqUrl.RawQuery = q.Encode()
	}

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws ListUsers: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws ListUsers: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ListUsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws ListUsers: error decoding response body: %w", err)
	}

	return &response, nil
}

// PatchUser updates a user in the AWS SSO Using the API
func (s *AWSSCIMService) PatchUser(ctx context.Context, pur *PatchUserRequest) error {
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

	resp, err := s.do(ctx, req)
	if err != nil {
		return fmt.Errorf("aws: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// PutUser creates a new user in the AWS SSO Using the API.
func (s *AWSSCIMService) PutUser(ctx context.Context, usr *PutUserRequest) (*PutUserResponse, error) {
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
		return nil, fmt.Errorf("aws PutUser: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Users/%s", usr.ID))

	// log.WithFields(log.Fields{
	// 	"group": usr.ID,
	// }).Debugf("patching group, payload: %s", utils.ToJSON(usr))

	req, err := s.newRequest(http.MethodPut, reqUrl, *usr)
	if err != nil {
		return nil, fmt.Errorf("aws PutUser: error creating request, http method: %s, url: %v, error: %w", http.MethodPut, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws PutUser: error sending request, http method: %s, url: %v, error: %w", http.MethodPut, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response PutUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws PutUser: error decoding response body: %w", err)
	}

	return &response, nil
}

// ListGroups returns a list of groups from the AWS SSO Using the API
func (s *AWSSCIMService) ListGroups(ctx context.Context, filter string) (*ListGroupsResponse, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws ListGroups: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Groups")

	if filter != "" {
		q := reqUrl.Query()
		q.Add("filter", filter)
		reqUrl.RawQuery = q.Encode()
	}

	// log.Debugf("req: %s", reqUrl.String())

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws ListGroups: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws ListGroups: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ListGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws ListGroups: error decoding response body: %w", err)
	}

	return &response, nil
}

// CreateGroup creates a new group in the AWS SSO Using the API
// reference:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/creategroup.html
func (s *AWSSCIMService) CreateGroup(ctx context.Context, g *CreateGroupRequest) (*CreateGroupResponse, error) {
	if g == nil {
		return nil, ErrCreateGroupRequestEmpty
	}
	if g.DisplayName == "" {
		return nil, ErrDisplayNameEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws CreateGroup: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/Groups")

	req, err := s.newRequest(http.MethodPost, reqUrl, *g)
	if err != nil {
		return nil, fmt.Errorf("aws CreateGroup: error creating request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws CreateGroup: error sending request, http method: %s, url: %v, error: %w", http.MethodPost, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response CreateGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws CreateGroup: error decoding response body: %v", err)
	}

	return &response, nil
}

// DeleteGroup deletes a group from the AWS SSO Using the API
func (s *AWSSCIMService) DeleteGroup(ctx context.Context, id string) error {
	if id == "" {
		return ErrGroupIDEmpty
	}

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws DeleteGroup: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Groups/%s", id))

	req, err := s.newRequest(http.MethodDelete, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("aws DeleteGroup: error creating request, http method: %s, url: %v, error: %w", http.MethodDelete, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return fmt.Errorf("aws DeleteGroup: error sending request, http method: %s, url: %v, error: %w", http.MethodDelete, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// PatchGroup updates a group in the AWS SSO Using the API
func (s *AWSSCIMService) PatchGroup(ctx context.Context, pgr *PatchGroupRequest) error {
	if pgr == nil {
		return ErrPatchGroupRequestEmpty
	}
	if pgr.Group.ID == "" {
		return ErrGroupIDEmpty
	}

	// log.WithFields(log.Fields{
	// 	"group": pgr.Group.ID,
	// }).Debugf("patching group, payload: %s", utils.ToJSON(pgr))

	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return fmt.Errorf("aws PatchGroup: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, fmt.Sprintf("/Groups/%s", pgr.Group.ID))

	req, err := s.newRequest(http.MethodPatch, reqUrl, pgr.Patch)
	if err != nil {
		return fmt.Errorf("aws PatchGroup: error creating request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return fmt.Errorf("aws PatchGroup: error sending request, http method: %s, url: %v, error: %w", http.MethodPatch, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	return nil
}

// ServiceProviderConfig returns additional information about the AWS SSO SCIM implementation
// references:
// + https://docs.aws.amazon.com/singlesignon/latest/developerguide/serviceproviderconfig.html
func (s *AWSSCIMService) ServiceProviderConfig(ctx context.Context) (*ServiceProviderConfig, error) {
	reqUrl, err := url.Parse(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("aws ServiceProviderConfig: error parsing url: %w", err)
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/ServiceProviderConfig")

	req, err := s.newRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("aws ServiceProviderConfig: error creating request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("aws ServiceProviderConfig: error sending request, http method: %s, url: %v, error: %w", http.MethodGet, reqUrl.String(), err)
	}
	defer resp.Body.Close()

	var response ServiceProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("aws ServiceProviderConfig: error decoding response body: %w", err)
	}

	return &response, nil
}
