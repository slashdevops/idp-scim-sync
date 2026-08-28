package idp

import (
	"context"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/idp"
	"go.uber.org/mock/gomock"
	admin "google.golang.org/api/admin/directory/v1"
)

// buildUser returns nil to reject a record the AWS SCIM API would refuse
// (missing name, given name, family name or primary email), logging a warning
// as it does so. Callers must honour that rejection instead of storing the nil:
// a nil element in UsersResult.Resources is dropped from the hash, but it also
// travels on into the reconciliation and SCIM layers where it dereferences.
//
// Google Workspace can legitimately return such records — accounts created via
// the Directory API without a family name, for example.

func TestGetUsers_skipsRejectedRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ds := mocks.NewMockGoogleProviderService(ctrl)
	ctx := context.Background()

	ds.EXPECT().ListUsers(ctx, gomock.Eq([]string{""})).Return([]*admin.User{
		{
			Id:           "1",
			PrimaryEmail: "good@example.com",
			Name:         &admin.UserName{GivenName: "Good", FamilyName: "User"},
		},
		// rejected: no family name
		{
			Id:           "2",
			PrimaryEmail: "nofamily@example.com",
			Name:         &admin.UserName{GivenName: "NoFamily"},
		},
		// rejected: no name at all
		{
			Id:           "3",
			PrimaryEmail: "noname@example.com",
		},
	}, nil).Times(1)

	ip, err := NewIdentityProvider(ds)
	if err != nil {
		t.Fatalf("NewIdentityProvider() error = %v", err)
	}

	got, err := ip.GetUsers(ctx, []string{""})
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}

	if got.Items != 1 {
		t.Errorf("GetUsers() Items = %d, want 1", got.Items)
	}
	if len(got.Resources) != 1 {
		t.Fatalf("GetUsers() returned %d resources, want 1", len(got.Resources))
	}
	for i, u := range got.Resources {
		if u == nil {
			t.Fatalf("GetUsers() Resources[%d] is nil; rejected records must be skipped, not stored", i)
		}
	}
	if got.Resources[0].UserName != "good@example.com" {
		t.Errorf("GetUsers() kept the wrong user: %q", got.Resources[0].UserName)
	}
}

func TestGetUsersByGroupsMembers_skipsRejectedRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ds := mocks.NewMockGoogleProviderService(ctrl)
	ctx := context.Background()

	ds.EXPECT().GetUser(gomock.Any(), "good@example.com").Return(&admin.User{
		Id:           "1",
		PrimaryEmail: "good@example.com",
		Name:         &admin.UserName{GivenName: "Good", FamilyName: "User"},
	}, nil).Times(1)

	// rejected by buildUser: no family name
	ds.EXPECT().GetUser(gomock.Any(), "nofamily@example.com").Return(&admin.User{
		Id:           "2",
		PrimaryEmail: "nofamily@example.com",
		Name:         &admin.UserName{GivenName: "NoFamily"},
	}, nil).Times(1)

	ip, err := NewIdentityProvider(ds)
	if err != nil {
		t.Fatalf("NewIdentityProvider() error = %v", err)
	}

	gmr := &model.GroupsMembersResult{
		Items: 1,
		Resources: []*model.GroupMembers{
			{
				Group: &model.Group{IPID: "g1", Name: "group1"},
				Items: 2,
				Resources: []*model.Member{
					{IPID: "1", Email: "good@example.com"},
					{IPID: "2", Email: "nofamily@example.com"},
				},
			},
		},
	}

	got, err := ip.GetUsersByGroupsMembers(ctx, gmr)
	if err != nil {
		t.Fatalf("GetUsersByGroupsMembers() error = %v", err)
	}

	if got.Items != 1 {
		t.Errorf("GetUsersByGroupsMembers() Items = %d, want 1", got.Items)
	}
	for i, u := range got.Resources {
		if u == nil {
			t.Fatalf("GetUsersByGroupsMembers() Resources[%d] is nil; rejected records must be skipped", i)
		}
	}
}
