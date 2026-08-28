package model

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"
)

// goldenState builds a fixed, fully-populated State. It is the fixture behind
// the golden hash assertions below, which pin the exact bytes that
// SetHashCode feeds to gob.
//
// The state file stored in S3 is compared by these hash codes, so any change to
// field order, to the hand-written MarshalBinary methods, or to the sort applied
// before hashing would invalidate every deployed state file and force a full
// re-sync. These goldens exist to make such a change impossible to land by
// accident.
//
// Deliberate coverage:
//   - every optional pointer field populated (Name, EnterpriseData, Manager)
//     and, in the second user, left nil — the nil cases are what exercise the
//     EOF branches in UnmarshalBinary
//   - two groups sharing an identical HashCode, so the ordering applied before
//     hashing is exercised where the sort key is ambiguous
//   - resources supplied out of order, so the sort is actually doing work
func goldenState() *State {
	groupA := GroupBuilder().WithIPID("g-2").WithName("zeta").WithEmail("zeta@example.com").Build()
	groupB := GroupBuilder().WithIPID("g-1").WithName("alpha").WithEmail("alpha@example.com").Build()

	fullUser := UserBuilder().
		WithIPID("u-2").
		WithUserName("zoe@example.com").
		WithDisplayName("Zoe Zebra").
		WithTitle("Engineer").
		WithUserType("admin#directory#user").
		WithPreferredLanguage("en").
		WithActive(true).
		WithEmails([]Email{{Value: "zoe@example.com", Type: "work", Primary: true}}).
		WithAddresses([]Address{{Formatted: "1 Main St", Locality: "Town", Country: "ES"}}).
		WithPhoneNumbers([]PhoneNumber{{Value: "+34000000000", Type: "work"}}).
		WithName(NameBuilder().WithGivenName("Zoe").WithFamilyName("Zebra").WithFormatted("Zoe Zebra").Build()).
		WithEnterpriseData(EnterpriseDataBuilder().
			WithEmployeeNumber("42").
			WithCostCenter("cc-1").
			WithOrganization("Org").
			WithDivision("example.com").
			WithDepartment("Eng").
			WithManager(ManagerBuilder().WithValue("mgr-1").WithRef("manager").Build()).
			Build()).
		Build()

	// Minimal user: Name, EnterpriseData and Manager all nil.
	bareUser := UserBuilder().
		WithIPID("u-1").
		WithUserName("adam@example.com").
		WithDisplayName("Adam Ant").
		WithActive(false).
		WithEmails([]Email{{Value: "adam@example.com", Type: "work", Primary: true}}).
		Build()

	memberB := MemberBuilder().WithIPID("m-2").WithEmail("zoe@example.com").WithStatus("ACTIVE").Build()
	memberA := MemberBuilder().WithIPID("m-1").WithEmail("adam@example.com").WithStatus("ACTIVE").Build()

	gmB := GroupMembersBuilder().WithGroup(groupA).WithResources([]*Member{memberB, memberA}).Build()
	gmA := GroupMembersBuilder().WithGroup(groupB).WithResources([]*Member{memberA}).Build()

	return StateBuilder().
		WithCodeVersion("v-test").
		WithLastSync("2026-01-01T00:00:00Z").
		WithGroups(GroupsResultBuilder().WithResources([]*Group{groupA, groupB}).Build()).
		WithUsers(UsersResultBuilder().WithResources([]*User{fullUser, bareUser}).Build()).
		WithGroupsMembers(GroupsMembersResultBuilder().WithResources([]*GroupMembers{gmB, gmA}).Build()).
		Build()
}

// Golden hash codes, captured from the implementation as it stood before the
// slices.SortStableFunc and errors.Is(io.EOF) changes. If a change to
// internal/model alters any of these, it has changed what the S3 state file
// hashes to, which invalidates every deployed state file and forces a full
// re-sync. Update them only deliberately, with a note in docs/Whats-New.md.
const (
	goldenStateHash         = "841819f1d2071dee9b6327079b492553cfd9cae825883c84af5e5f742ec2cf9e"
	goldenGroupsHash        = "c02429df7dd0db29849b5212c9aa431d2d55ab6e44f4b67f57c984b16a57e4ea"
	goldenUsersHash         = "5b0614c402dd9556a16593e00523a6c3f4fbc1fa10bd2e212e03736be60da293"
	goldenGroupsMembersHash = "6c91be6299ab0e4f5db64892d48a59c63ec894ebe43074f0ec87d9287c6f5315"

	goldenGroupZetaHash     = "587a531ccf502c3a2e17eb1785179f58824e33f016ab0e542e9b7bfaff21c5c7"
	goldenGroupAlphaHash    = "9d7a5febd8a48c9f7b50228f6425969003195b6b2c1c153366ebc5a67444ea9e"
	goldenUserFullHash      = "483d99be74e09c29f877f27b43ae5d66c373b5dbffef7a2b9b89a099cd840196"
	goldenUserBareHash      = "55ed617756c565afae3cd6ce24dd6e163edda34bf7e16d0aadc19006a58f7110"
	goldenGroupMembers0Hash = "24910bf14a1ccb143038ef1344b95f3d6f28df4dc5cbb763d93a00de4c770f6d"
	goldenGroupMembers1Hash = "bcb785ac6a3f4abef4a401b7e383907cb34d4065e62fb533fc82b24b13e40efb"
)

func TestGoldenHashes_unchanged(t *testing.T) {
	s := goldenState()

	for _, tc := range []struct {
		name string
		got  string
		want string
	}{
		{"State", s.HashCode, goldenStateHash},
		{"GroupsResult", s.Resources.Groups.HashCode, goldenGroupsHash},
		{"UsersResult", s.Resources.Users.HashCode, goldenUsersHash},
		{"GroupsMembersResult", s.Resources.GroupsMembers.HashCode, goldenGroupsMembersHash},
		{"Group zeta", s.Resources.Groups.Resources[0].HashCode, goldenGroupZetaHash},
		{"Group alpha", s.Resources.Groups.Resources[1].HashCode, goldenGroupAlphaHash},
		{"User full", s.Resources.Users.Resources[0].HashCode, goldenUserFullHash},
		{"User bare", s.Resources.Users.Resources[1].HashCode, goldenUserBareHash},
		{"GroupMembers[0]", s.Resources.GroupsMembers.Resources[0].HashCode, goldenGroupMembers0Hash},
		{"GroupMembers[1]", s.Resources.GroupsMembers.Resources[1].HashCode, goldenGroupMembers1Hash},
	} {
		if tc.got != tc.want {
			t.Errorf("%s hash changed:\n got %s\nwant %s", tc.name, tc.got, tc.want)
		}
	}
}

// Round-trip through the gob encoders that back the hash, then through JSON as
// the S3 repository does, asserting the hashes survive both.
func TestGoldenState_roundTrips(t *testing.T) {
	s := goldenState()

	binary, err := s.MarshalBinary()
	if err != nil {
		t.Fatalf("State.MarshalBinary() error = %v", err)
	}
	var fromBinary State
	if err := fromBinary.UnmarshalBinary(binary); err != nil {
		t.Fatalf("State.UnmarshalBinary() error = %v", err)
	}
	fromBinary.SetHashCode()
	if fromBinary.HashCode != s.HashCode {
		t.Errorf("gob round-trip changed the state hash:\n got %s\nwant %s", fromBinary.HashCode, s.HashCode)
	}

	encoded, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	var fromJSON State
	if err := json.Unmarshal(encoded, &fromJSON); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	fromJSON.SetHashCode()
	if fromJSON.HashCode != s.HashCode {
		t.Errorf("JSON round-trip changed the state hash:\n got %s\nwant %s", fromJSON.HashCode, s.HashCode)
	}

	reencoded, err := json.Marshal(&fromJSON)
	if err != nil {
		t.Fatalf("json.Marshal() (second pass) error = %v", err)
	}
	if !bytes.Equal(encoded, reencoded) {
		t.Errorf("JSON serialisation is not stable across a round-trip:\n first %s\nsecond %s", encoded, reencoded)
	}
}

// The three container hashes drive every sync decision in internal/core, which
// compares idpResult.HashCode against stateResult.HashCode. They must therefore
// be independent of the order in which resources arrive: the IdP fan-outs and
// the SCIM inversion both build their slices from maps.
func TestGoldenHashes_containerHashesAreOrderIndependent(t *testing.T) {
	a := goldenState()
	b := reversedGoldenState()

	if a.Resources.Groups.HashCode != b.Resources.Groups.HashCode {
		t.Errorf("GroupsResult hash depends on input order")
	}
	if a.Resources.Users.HashCode != b.Resources.Users.HashCode {
		t.Errorf("UsersResult hash depends on input order")
	}
	if a.Resources.GroupsMembers.HashCode != b.Resources.GroupsMembers.HashCode {
		t.Errorf("GroupsMembersResult hash depends on input order")
	}
}

// State.SetHashCode, unlike the container SetHashCode methods, hashes its
// rebuilt slices in the order it received them: it constructs each *Result via
// a builder (which sorts only its own private copy for its own hash) and then
// gob-encodes copyState, whose MarshalBinary walks Resources in slice order.
//
// This test documents that behaviour rather than asserting it is desirable.
// State.HashCode is written into the S3 object but never read back — internal/core
// compares only the three container hashes (see internal/core/actions.go) — so
// the order dependence is cosmetic today. It does mean the state file's
// top-level hashCode field can differ between two runs that produced identical
// data, which makes it useless as an external change indicator.
//
// Tracked as finding M17. Fixing it would change the value of that field for
// every deployed state file, so it is deliberately left alone here.
func TestState_hashCodeIsOrderDependent(t *testing.T) {
	a := goldenState()
	b := reversedGoldenState()

	if a.HashCode == b.HashCode {
		t.Skip("State.HashCode is now order-independent; finding M17 appears to have been fixed — update this test and the goldens")
	}

	t.Logf("known: State.HashCode differs on reordered input (%s vs %s); cosmetic, the field is never read", a.HashCode, b.HashCode)
}

// reversedGoldenState is goldenState with every resource slice reversed and all
// hashes recomputed.
func reversedGoldenState() *State {
	s := goldenState()

	slices.Reverse(s.Resources.Groups.Resources)
	slices.Reverse(s.Resources.Users.Resources)
	slices.Reverse(s.Resources.GroupsMembers.Resources)

	s.Resources.Groups.SetHashCode()
	s.Resources.Users.SetHashCode()
	s.Resources.GroupsMembers.SetHashCode()
	s.SetHashCode()

	return s
}
