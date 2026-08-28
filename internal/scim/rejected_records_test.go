package scim

import (
	"context"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	mock_scim "github.com/slashdevops/idp-scim-sync/mocks/scim"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"go.uber.org/mock/gomock"
)

// GetGroupsMembers indexed user.Emails[0] directly. buildUser only keeps emails
// flagged primary, so an AWS SCIM user with "emails": [] — or with emails but
// none primary — yields a nil Emails slice and the index panics. AWS IAM
// Identity Center does return that shape; pkg/aws/testdata/ListUserResponse_no_emails.json
// is a checked-in fixture of it.
//
// Because the state file is only written after a successful sync, a panic here
// wedges the sync permanently: the next run re-reads the same user and panics
// again.
func TestProvider_GetGroupsMembers_userWithoutEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_scim.NewMockAWSSCIMProvider(ctrl)
	m.EXPECT().
		ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
		Return(&aws.ListGroupsResponse{
			Resources: []*aws.Group{{ID: "g1", DisplayName: "group1"}},
		}, nil)

	p := &Provider{scim: m, maxMembersPerRequest: 100}

	gr := &model.GroupsResult{Resources: []*model.Group{{SCIMID: "g1", Name: "group1"}}}
	ur := &model.UsersResult{Resources: []*model.User{
		{SCIMID: "u1", UserName: "noemail@example.com", Active: true, Emails: nil},
	}}

	got, err := p.GetGroupsMembers(context.Background(), gr, ur)
	if err != nil {
		t.Fatalf("GetGroupsMembers() error = %v", err)
	}
	if len(got.Resources) != 1 {
		t.Fatalf("GetGroupsMembers() returned %d groups, want 1", len(got.Resources))
	}
	if len(got.Resources[0].Resources) != 1 {
		t.Fatalf("GetGroupsMembers() returned %d members, want 1", len(got.Resources[0].Resources))
	}
	if scimID := got.Resources[0].Resources[0].SCIMID; scimID != "u1" {
		t.Errorf("member SCIMID = %q, want %q", scimID, "u1")
	}
}

// Users whose emails exist but where none is flagged primary hit the same path.
func TestProvider_GetGroupsMembers_userWithNoPrimaryEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_scim.NewMockAWSSCIMProvider(ctrl)
	m.EXPECT().
		ListGroupsWithCursor(gomock.Any(), `members.value eq "u1"`, "").
		Return(&aws.ListGroupsResponse{
			Resources: []*aws.Group{{ID: "g1", DisplayName: "group1"}},
		}, nil)

	p := &Provider{scim: m, maxMembersPerRequest: 100}

	gr := &model.GroupsResult{Resources: []*model.Group{{SCIMID: "g1", Name: "group1"}}}
	ur := &model.UsersResult{Resources: []*model.User{
		{
			SCIMID: "u1",
			Active: true,
			Emails: []model.Email{{Value: "secondary@example.com", Primary: false}},
		},
	}}

	got, err := p.GetGroupsMembers(context.Background(), gr, ur)
	if err != nil {
		t.Fatalf("GetGroupsMembers() error = %v", err)
	}
	if len(got.Resources[0].Resources) != 1 {
		t.Fatalf("GetGroupsMembers() returned %d members, want 1", len(got.Resources[0].Resources))
	}
	// No address is flagged primary, so the first one is used rather than
	// dropping the member's email entirely.
	if email := got.Resources[0].Resources[0].Email; email != "secondary@example.com" {
		t.Errorf("member Email = %q, want %q", email, "secondary@example.com")
	}
}

// scim.buildUser returns nil when the AWS response carries no id. GetUsers
// stored that nil straight into UsersResult.Resources.
func TestProvider_GetUsers_skipsRejectedRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_scim.NewMockAWSSCIMProvider(ctrl)
	m.EXPECT().ListUsers(gomock.Any(), "").Return(&aws.ListUsersResponse{
		Resources: []*aws.User{
			{
				ID:          "u1",
				UserName:    "good@example.com",
				DisplayName: "Good User",
				Name:        &aws.Name{GivenName: "Good", FamilyName: "User"},
				Emails:      []aws.Email{{Value: "good@example.com", Type: "work", Primary: true}},
			},
			// rejected by buildUser: empty id
			{
				UserName:    "noid@example.com",
				DisplayName: "No ID",
				Name:        &aws.Name{GivenName: "No", FamilyName: "ID"},
			},
		},
	}, nil)

	p := &Provider{scim: m, maxMembersPerRequest: 100}

	got, err := p.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}

	if got.Items != 1 {
		t.Errorf("GetUsers() Items = %d, want 1", got.Items)
	}
	for i, u := range got.Resources {
		if u == nil {
			t.Fatalf("GetUsers() Resources[%d] is nil; rejected records must be skipped", i)
		}
	}
}
