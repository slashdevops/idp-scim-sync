package aws

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mocks "github.com/slashdevops/idp-scim-sync/mocks/aws"
)

// conflictThenEmptyLookup replies 409 to every POST and returns an empty
// Resources list to every GET, which is the state that drives CreateOrGetUser
// and CreateOrGetGroup into their "clear ExternalID and try once more" branch.
//
// The production code called itself with no attempt limit, so this state would
// recurse until the goroutine stack was exhausted. The doTracker counts POSTs
// and hard-stops well before that, so the test reports a bounded number instead
// of crashing the test binary.
type doTracker struct {
	posts int
	limit int
}

func (d *doTracker) do(req *http.Request) (*http.Response, error) {
	body := func(status int, s string) *http.Response {
		return &http.Response{
			StatusCode:    status,
			Header:        http.Header{"Content-Type": []string{"application/json"}},
			Body:          io.NopCloser(strings.NewReader(s)),
			ContentLength: int64(len(s)),
		}
	}

	if req.Method == http.MethodPost {
		d.posts++
		if d.posts > d.limit {
			// Break the loop so the assertion below can report the count.
			return body(http.StatusInternalServerError, `{"detail":"tracker limit reached"}`), nil
		}
		return body(http.StatusConflict, `{"detail":"Duplicate UserName","scimType":"uniqueness","status":"409"}`), nil
	}

	// The follow-up lookup finds nothing, so response.ID stays empty.
	return body(http.StatusOK, `{"totalResults":0,"itemsPerPage":0,"startIndex":1,"Resources":[]}`), nil
}

// maxConflictRetryPosts is the number of POSTs a bounded implementation may
// issue: the original attempt plus exactly one retry with ExternalID cleared.
const maxConflictRetryPosts = 2

func TestCreateOrGetUser_conflictRetryIsBounded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracker := &doTracker{limit: 50}
	mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
	mockHTTPClient.EXPECT().Do(gomock.Any()).DoAndReturn(tracker.do).AnyTimes()

	service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
	assert.NoError(t, err)

	_, err = service.CreateOrGetUser(context.Background(), &CreateUserRequest{
		ID:          "1",
		ExternalID:  "ext-1",
		UserName:    "user.1@mail.com",
		DisplayName: "user 1",
		Name:        &Name{FamilyName: "1", GivenName: "test"},
		Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
		Active:      true,
	})

	assert.Error(t, err, "a conflict that never resolves must surface as an error")
	assert.LessOrEqual(t, tracker.posts, maxConflictRetryPosts,
		"CreateOrGetUser issued %d POSTs; the conflict retry must be bounded to %d",
		tracker.posts, maxConflictRetryPosts)
}

func TestCreateOrGetGroup_conflictRetryIsBounded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tracker := &doTracker{limit: 50}
	mockHTTPClient := mocks.NewMockHTTPClient(mockCtrl)
	mockHTTPClient.EXPECT().Do(gomock.Any()).DoAndReturn(tracker.do).AnyTimes()

	service, err := NewSCIMService(mockHTTPClient, "https://testing.com", "MyToken")
	assert.NoError(t, err)

	_, err = service.CreateOrGetGroup(context.Background(), &CreateGroupRequest{
		DisplayName: "group 1",
		ExternalID:  "ext-g1",
	})

	assert.Error(t, err, "a conflict that never resolves must surface as an error")
	assert.LessOrEqual(t, tracker.posts, maxConflictRetryPosts,
		"CreateOrGetGroup issued %d POSTs; the conflict retry must be bounded to %d",
		tracker.posts, maxConflictRetryPosts)
}
