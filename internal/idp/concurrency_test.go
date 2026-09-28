package idp

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"sync"
	"testing"

	"go.uber.org/mock/gomock"
	admin "google.golang.org/api/admin/directory/v1"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/idp"
)

// groupsMembersWithEmails builds a single-group GroupsMembersResult whose
// members carry the given emails.
func groupsMembersWithEmails(emails ...string) *model.GroupsMembersResult {
	members := make([]*model.Member, 0, len(emails))
	for i, email := range emails {
		members = append(members, &model.Member{IPID: strconv.Itoa(i + 1), Email: email})
	}
	return &model.GroupsMembersResult{
		Items: 1,
		Resources: []*model.GroupMembers{
			{Group: &model.Group{IPID: "g1", Name: "group1"}, Items: len(members), Resources: members},
		},
	}
}

func adminUser(id, email string) *admin.User {
	return &admin.User{
		Id:           id,
		PrimaryEmail: email,
		Name:         &admin.UserName{GivenName: "Given" + id, FamilyName: "Family" + id},
	}
}

// GetUsersByGroupsMembers fans out one GetUser call per unique member email.
// The hand-rolled WaitGroup + semaphore + buffered error channel it used had no
// cancellation: when one call failed, every sibling ran to completion, spending
// Google Directory API quota on results that were then discarded.
//
// With errgroup.WithContext the shared context is cancelled as soon as the first
// error is returned, so callers that honour ctx stop early.
func TestGetUsersByGroupsMembers_cancelsSiblingsOnFirstError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ds := mocks.NewMockGoogleProviderService(ctrl)

	const total = 40
	emails := make([]string, 0, total)
	for i := range total {
		emails = append(emails, "user"+strconv.Itoa(i)+"@example.com")
	}

	var (
		mu        sync.Mutex
		sawCancel bool
	)

	ds.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, userID string) (*admin.User, error) {
			// A cancelled context must be observable by the remaining calls.
			if ctx.Err() != nil {
				mu.Lock()
				sawCancel = true
				mu.Unlock()
				return nil, ctx.Err()
			}
			return nil, errors.New("google: boom for " + userID)
		}).AnyTimes()

	ip, err := NewIdentityProvider(ds)
	if err != nil {
		t.Fatalf("NewIdentityProvider() error = %v", err)
	}

	got, err := ip.GetUsersByGroupsMembers(context.Background(), groupsMembersWithEmails(emails...))
	if err == nil {
		t.Fatal("GetUsersByGroupsMembers() expected an error")
	}
	if got != nil {
		t.Errorf("GetUsersByGroupsMembers() returned %v, want nil on error", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if !sawCancel {
		t.Error("no in-flight GetUser call observed a cancelled context; the fan-out does not cancel siblings on first error")
	}
}

// The result was assembled by ranging a map, so its order changed run to run.
// The hash is order-insensitive (SetHashCode sorts a copy), but the state file
// written to S3 preserves slice order, so every sync produced a large spurious
// diff. Ordering by primary email makes the artifact reviewable.
func TestGetUsersByGroupsMembers_deterministicOrder(t *testing.T) {
	emails := []string{
		"zoe@example.com", "adam@example.com", "mia@example.com",
		"kai@example.com", "ben@example.com", "yara@example.com",
	}

	var first []string
	for run := range 8 {
		ctrl := gomock.NewController(t)

		ds := mocks.NewMockGoogleProviderService(ctrl)
		ds.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, userID string) (*admin.User, error) {
				return adminUser("id-"+userID, userID), nil
			}).Times(len(emails))

		ip, err := NewIdentityProvider(ds)
		if err != nil {
			t.Fatalf("NewIdentityProvider() error = %v", err)
		}

		got, err := ip.GetUsersByGroupsMembers(context.Background(), groupsMembersWithEmails(emails...))
		if err != nil {
			t.Fatalf("GetUsersByGroupsMembers() error = %v", err)
		}

		order := make([]string, 0, len(got.Resources))
		for _, u := range got.Resources {
			order = append(order, u.UserName)
		}

		if run == 0 {
			first = order
			want := slices.Sorted(slices.Values(emails))
			if !slices.Equal(order, want) {
				t.Fatalf("resources are not ordered by primary email:\n got %v\nwant %v", order, want)
			}
		} else if !slices.Equal(first, order) {
			t.Fatalf("run %d produced a different order:\nfirst %v\n  now %v", run, first, order)
		}

		ctrl.Finish()
	}
}
