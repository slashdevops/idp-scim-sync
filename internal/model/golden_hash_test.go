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

// Golden hash codes. If a change to internal/model alters any of these, it has
// changed what the S3 state file hashes to. Update them only deliberately, with
// a note in docs/Whats-New.md.
//
// Two were updated when the hash input was given a canonical ordering at every
// level (findings M17 and M23):
//
//   - goldenGroupsMembersHash changed because GroupsMembersResult previously
//     sorted its outer slice but not the members nested inside each entry, so
//     member order within a group leaked into the hash. internal/core reads this
//     hash to decide whether membership needs reconciling, so the change costs
//     one membership reconciliation pass on the first run after deploying.
//   - goldenStateHash changed because State.SetHashCode sorted nothing at all.
//     Nothing reads State.HashCode, so this one has no runtime effect.
//
// goldenGroupsHash, goldenUsersHash and every individual resource hash are
// deliberately UNCHANGED: groups and users do not re-sync.
const (
	goldenStateHash         = "e7a1cd954b805df554106ca5a3a4cac837f11e50af4bf3b0ea5542a775bcc2af"
	goldenGroupsHash        = "c02429df7dd0db29849b5212c9aa431d2d55ab6e44f4b67f57c984b16a57e4ea"
	goldenUsersHash         = "5b0614c402dd9556a16593e00523a6c3f4fbc1fa10bd2e212e03736be60da293"
	goldenGroupsMembersHash = "f391383ec27336152cafeb3c2796f4303ad5f77bd6446d734bd58316a905a204"

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

// State.SetHashCode must be order-independent for the same reason the container
// SetHashCode methods are: the identity-provider fan-outs build their slices
// from maps, so two runs over identical data can present resources in different
// orders.
//
// It was not. State.SetHashCode rebuilt each *Result through a builder and then
// gob-encoded copyState, whose MarshalBinary chain walks Resources in slice
// order. The builder's Build() sorts only its own private copy, to compute that
// container's own HashCode field — the sort never reached the bytes the State
// hash is built from.
//
// State.HashCode is written into the S3 object but never read back
// (internal/core/actions.go compares only the three container hashes), so this
// was never a correctness bug. It did make the state file's top-level hashCode
// field useless as an external change indicator, which is a trap worth closing.
func TestGoldenHashes_stateHashIsOrderIndependent(t *testing.T) {
	a := goldenState()
	b := reversedGoldenState()

	if a.HashCode != b.HashCode {
		t.Errorf("State hash depends on input order:\n  in-order %s\n  reversed %s", a.HashCode, b.HashCode)
	}
}

// reversedGoldenState is goldenState with every resource slice reversed and all
// hashes recomputed.
func reversedGoldenState() *State {
	s := goldenState()

	slices.Reverse(s.Resources.Groups.Resources)
	slices.Reverse(s.Resources.Users.Resources)
	slices.Reverse(s.Resources.GroupsMembers.Resources)

	// Reverse the members inside each group too: GroupMembers.MarshalBinary
	// walks them in slice order, so nested ordering reaches the hash bytes just
	// as the outer ordering does.
	for _, gm := range s.Resources.GroupsMembers.Resources {
		slices.Reverse(gm.Resources)
	}

	s.Resources.Groups.SetHashCode()
	s.Resources.Users.SetHashCode()
	s.Resources.GroupsMembers.SetHashCode()
	s.SetHashCode()

	return s
}

// Member order within a group must not affect any hash.
//
// GroupMembers.MarshalBinary walks its Resources in slice order, so before
// findings M17/M23 the order of members inside a group reached the bytes of both
// the enclosing GroupsMembersResult hash and the State hash — while
// GroupMembers' own hash was stable, because it sorted its own copy.
//
// That mattered in production because internal/core compares
// GroupsMembersResult hashes to decide whether membership needs reconciling, and
// two upstream sources present members in an unstable order: the Google
// Directory API makes no ordering guarantee, and internal/scim's membership
// inversion appends members in goroutine completion order.
func TestGroupsMembersResult_hashIgnoresWithinGroupMemberOrder(t *testing.T) {
	build := func(reverse bool) (gm, gmr, state string) {
		group := GroupBuilder().WithIPID("g-1").WithName("alpha").Build()
		members := []*Member{
			MemberBuilder().WithIPID("m-1").WithEmail("adam@example.com").WithStatus("ACTIVE").Build(),
			MemberBuilder().WithIPID("m-2").WithEmail("zoe@example.com").WithStatus("ACTIVE").Build(),
		}
		if reverse {
			slices.Reverse(members)
		}

		entry := GroupMembersBuilder().WithGroup(group).WithResources(members).Build()
		result := GroupsMembersResultBuilder().WithResources([]*GroupMembers{entry}).Build()
		st := StateBuilder().WithGroupsMembers(result).Build()

		return entry.HashCode, result.HashCode, st.HashCode
	}

	aGM, aGMR, aState := build(false)
	bGM, bGMR, bState := build(true)

	if aGM != bGM {
		t.Errorf("GroupMembers hash depends on member order:\n got %s\nwant %s", bGM, aGM)
	}
	if aGMR != bGMR {
		t.Errorf("GroupsMembersResult hash depends on member order (this one is read by internal/core):\n got %s\nwant %s", bGMR, aGMR)
	}
	if aState != bState {
		t.Errorf("State hash depends on member order:\n got %s\nwant %s", bState, aState)
	}
}

// Computing a hash must not reorder the caller's data. The copies inside
// SetHashCode exist for exactly this reason, and the nested member slices need
// their own copy because copying a GroupMembers struct shares its slice backing
// array.
func TestSetHashCode_doesNotMutateCallerOrdering(t *testing.T) {
	group := GroupBuilder().WithIPID("g-1").WithName("alpha").Build()
	zoe := MemberBuilder().WithIPID("m-2").WithEmail("zoe@example.com").Build()
	adam := MemberBuilder().WithIPID("m-1").WithEmail("adam@example.com").Build()

	// Deliberately not in email order.
	entry := GroupMembersBuilder().WithGroup(group).WithResources([]*Member{zoe, adam}).Build()
	result := GroupsMembersResultBuilder().WithResources([]*GroupMembers{entry}).Build()

	result.SetHashCode()

	if got := result.Resources[0].Resources[0].Email; got != "zoe@example.com" {
		t.Errorf("SetHashCode reordered the caller's members: first is now %q, want %q", got, "zoe@example.com")
	}
}
