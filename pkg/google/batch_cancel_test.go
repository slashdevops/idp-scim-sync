package google

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

// ListGroupMembersBatch fans out one members-list call per group. The
// hand-rolled WaitGroup + semaphore + buffered error channel it used had no
// cancellation: when one group's fetch failed, every sibling still ran to
// completion, spending Google Directory API quota on results that were
// discarded, and only the first error off the channel was ever returned.
//
// With errgroup.WithContext the shared context is cancelled on the first error,
// so the google-api-go-client transport aborts the requests still in flight and
// the remaining queued ones never start.
func TestListGroupMembersBatch_cancelsSiblingsOnFirstError(t *testing.T) {
	ctx := context.Background()

	var (
		mu       sync.Mutex
		served   int
		canceled int
	)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		served++
		mu.Unlock()

		// The very first request fails, which must cancel the group context.
		if r.URL.Query().Get("pageToken") == "" && r.Header.Get("X-Fail") == "" {
			// Block until either the client gives up (context cancelled) or a
			// short grace period elapses, so we can observe cancellation.
			select {
			case <-r.Context().Done():
				mu.Lock()
				canceled++
				mu.Unlock()
				return
			default:
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	svc, err := admin.NewService(ctx,
		option.WithHTTPClient(svr.Client()),
		option.WithEndpoint(svr.URL),
		option.WithUserAgent("test"),
	)
	assert.NoError(t, err)

	client, err := NewDirectoryService(svc)
	assert.NoError(t, err)

	const groupCount = 60
	groupIDs := make([]string, 0, groupCount)
	for i := range groupCount {
		groupIDs = append(groupIDs, "group-"+strconv.Itoa(i))
	}

	got, err := client.ListGroupMembersBatch(ctx, groupIDs)
	assert.Error(t, err)
	assert.Nil(t, got)

	mu.Lock()
	defer mu.Unlock()
	// Without cancellation every group is attempted. With it, the fan-out stops
	// early and a meaningful fraction of the groups are never requested.
	assert.Less(t, served, groupCount,
		"all %d groups were requested (%d served); the fan-out does not stop after the first error",
		groupCount, served)
}

// ListGroupMembersBatch must honour a context that is already cancelled rather
// than issuing the whole fan-out anyway.
func TestListGroupMembersBatch_respectsCancelledContext(t *testing.T) {
	var (
		mu     sync.Mutex
		served int
	)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		served++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"kind":"directory#members","members":[]}`))
	}))
	defer svr.Close()

	svc, err := admin.NewService(context.Background(),
		option.WithHTTPClient(svr.Client()),
		option.WithEndpoint(svr.URL),
		option.WithUserAgent("test"),
	)
	assert.NoError(t, err)

	client, err := NewDirectoryService(svc)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := client.ListGroupMembersBatch(ctx, []string{"g1", "g2", "g3"})
	assert.Error(t, err)
	assert.Nil(t, got)

	mu.Lock()
	defer mu.Unlock()
	assert.Zero(t, served, "requests were issued despite an already-cancelled context")
}
