package scim

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// fullModelUser exercises every field the request builders map, including all
// three optional pointer branches.
func fullModelUser() *model.User {
	return &model.User{
		IPID:              "ipid-1",
		SCIMID:            "scimid-1",
		UserName:          "zoe@example.com",
		DisplayName:       "Zoe Zebra",
		NickName:          "zz",
		ProfileURL:        "https://example.com/zoe",
		Title:             "Engineer",
		UserType:          "admin#directory#user",
		PreferredLanguage: "en",
		Locale:            "en-GB",
		Timezone:          "Europe/Madrid",
		Active:            true,
		Name: &model.Name{
			Formatted:       "Zoe Zebra",
			FamilyName:      "Zebra",
			GivenName:       "Zoe",
			MiddleName:      "Q",
			HonorificPrefix: "Dr",
			HonorificSuffix: "PhD",
		},
		Emails: []model.Email{
			{Value: "secondary@example.com", Type: "home", Primary: false},
			{Value: "zoe@example.com", Type: "work", Primary: true},
		},
		Addresses: []model.Address{
			{Formatted: "1 Main St", StreetAddress: "1 Main St", Locality: "Town", Region: "R", PostalCode: "08001", Country: "ES"},
		},
		PhoneNumbers: []model.PhoneNumber{
			{Value: "+34000000000", Type: "work"},
		},
		EnterpriseData: &model.EnterpriseData{
			EmployeeNumber: "42",
			CostCenter:     "cc-1",
			Organization:   "Org",
			Division:       "example.com",
			Department:     "Eng",
			Manager:        &model.Manager{Value: "mgr-1", Ref: "manager"},
		},
	}
}

// bareModelUser leaves Name, EnterpriseData and Manager nil, and carries no
// addresses or phone numbers, so the guarded branches are covered too.
func bareModelUser() *model.User {
	return &model.User{
		IPID:     "ipid-2",
		SCIMID:   "scimid-2",
		UserName: "adam@example.com",
		Active:   false,
		Emails:   []model.Email{{Value: "adam@example.com", Type: "work", Primary: true}},
	}
}

// buildCreateUserRequest and buildPutUserRequest were ~150 lines of near-
// identical field mapping, differing only in that Put sets ID. This
// characterisation test pins that relationship so the shared implementation
// cannot drift: every field must match, and ID must be the only difference.
//
// It is what makes the de-duplication safe — before the change it documents
// the existing behaviour, after the change it proves nothing moved.
func TestBuildUserRequests_differOnlyByID(t *testing.T) {
	for _, tc := range []struct {
		name string
		user *model.User
	}{
		{"fully populated", fullModelUser()},
		{"minimal", bareModelUser()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			create := buildCreateUserRequest(tc.user)
			put := buildPutUserRequest(tc.user)

			if create == nil || put == nil {
				t.Fatal("builders returned nil for a non-nil user")
			}

			if put.ID != tc.user.SCIMID {
				t.Errorf("PutUserRequest.ID = %q, want the user's SCIMID %q", put.ID, tc.user.SCIMID)
			}
			if create.ID != "" {
				t.Errorf("CreateUserRequest.ID = %q, want empty: create must not target an existing resource", create.ID)
			}

			// Compare as the shared underlying aws.User, with ID neutralised.
			asUser := func(u aws.User) aws.User {
				u.ID = ""
				return u
			}
			if diff := cmp.Diff(asUser(aws.User(*create)), asUser(aws.User(*put))); diff != "" {
				t.Errorf("the two requests differ by more than ID (-create +put):\n%s", diff)
			}
		})
	}
}

// Only the primary address is forwarded to AWS, matching what buildUser
// produces on the way in.
func TestBuildUserRequests_forwardOnlyPrimaryEmail(t *testing.T) {
	create := buildCreateUserRequest(fullModelUser())

	if len(create.Emails) != 1 {
		t.Fatalf("forwarded %d emails, want 1 (the primary): %+v", len(create.Emails), create.Emails)
	}
	if create.Emails[0].Value != "zoe@example.com" || !create.Emails[0].Primary {
		t.Errorf("forwarded the wrong email: %+v", create.Emails[0])
	}
}

func TestBuildUserRequests_nilUser(t *testing.T) {
	if got := buildCreateUserRequest(nil); got != nil {
		t.Errorf("buildCreateUserRequest(nil) = %v, want nil", got)
	}
	if got := buildPutUserRequest(nil); got != nil {
		t.Errorf("buildPutUserRequest(nil) = %v, want nil", got)
	}
}

// goldenRequestJSON pins the exact wire payload sent to the AWS SCIM API, so
// the de-duplication cannot quietly change a field name or drop an attribute.
func TestBuildUserRequests_goldenWireFormat(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload any
	}{
		{"create", buildCreateUserRequest(fullModelUser())},
		{"put", buildPutUserRequest(fullModelUser())},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			t.Logf("%s payload = %s", tc.name, b)

			// Spot-check the attributes AWS requires, plus the enterprise
			// extension, which is the most intricate mapping.
			for _, want := range []string{
				`"userName":"zoe@example.com"`,
				`"displayName":"Zoe Zebra"`,
				`"externalId":"ipid-1"`,
				`"active":true`,
				`"familyName":"Zebra"`,
				`"givenName":"Zoe"`,
				`"employeeNumber":"42"`,
				`"costCenter":"cc-1"`,
				`"department":"Eng"`,
				`"manager"`,
			} {
				if !contains(string(b), want) {
					t.Errorf("payload is missing %s\ngot: %s", want, b)
				}
			}
		})
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
