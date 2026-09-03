package idp

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
	admin "google.golang.org/api/admin/directory/v1"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/idp"
)

// Go 1.27 prints runtime/pprof goroutine labels in traceback headers, e.g.
//
//	goroutine 42 [running] {sync: "users", user: "zoe@example.com"}:
//
// That turns "something in the user fan-out panicked" into "this record
// panicked", which is the difference between a one-line CloudWatch diagnosis and
// bisecting a directory. The C1/C2 crashes fixed earlier in this branch were
// exactly this shape.
//
// The labels must be set with pprof.SetGoroutineLabels, NOT pprof.Do: pprof.Do
// defers restoring the previous label set, and on a panic those defers run
// during unwinding before the runtime prints the traceback, so the labels are
// already gone. runtime.Stack still shows them, which makes the broken version
// easy to "verify" by accident.
//
// This test captures a traceback from inside a worker and asserts the labels are
// present, so nobody can simplify it back to pprof.Do.
func TestGetUsersByGroupsMembers_labelsGoroutinesForTracebacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ds := mocks.NewMockGoogleProviderService(ctrl)

	var traceback string
	ds.EXPECT().GetUser(gomock.Any(), "zoe@example.com").DoAndReturn(
		func(_ context.Context, userID string) (*admin.User, error) {
			buf := make([]byte, 16<<10)
			traceback = string(buf[:runtime.Stack(buf, false)])
			return &admin.User{
				Id:           "u-1",
				PrimaryEmail: userID,
				Name:         &admin.UserName{GivenName: "Zoe", FamilyName: "Zebra"},
			}, nil
		}).Times(1)

	ip, err := NewIdentityProvider(ds)
	if err != nil {
		t.Fatalf("NewIdentityProvider() error = %v", err)
	}

	gmr := &model.GroupsMembersResult{
		Items: 1,
		Resources: []*model.GroupMembers{
			{
				Group:     &model.Group{IPID: "g-1", Name: "group1"},
				Items:     1,
				Resources: []*model.Member{{IPID: "m-1", Email: "zoe@example.com"}},
			},
		},
	}

	if _, err := ip.GetUsersByGroupsMembers(context.Background(), gmr); err != nil {
		t.Fatalf("GetUsersByGroupsMembers() error = %v", err)
	}

	// Only the first line matters: the goroutine header carries the labels.
	header := traceback
	if i := strings.IndexByte(header, '\n'); i >= 0 {
		header = header[:i]
	}

	for _, want := range []string{`sync`, `users`, `user`, `zoe@example.com`} {
		if !strings.Contains(header, want) {
			t.Errorf("goroutine header is missing %q, so a panic here could not be attributed to a record.\nheader: %s", want, header)
		}
	}
}
