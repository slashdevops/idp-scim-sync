package model

import (
	"slices"
	"strings"
)

// Canonical orderings for hash input.
//
// Every hash in this package is a gob encoding, and the hand-written
// MarshalBinary methods walk their Resources slices in slice order. The order
// resources arrive in is therefore part of the hash unless it is normalised
// first — and it must be normalised, because the identity-provider fan-outs and
// the SCIM membership inversion all build their slices from maps, whose
// iteration order is randomised.
//
// These helpers exist so that every SetHashCode normalises the same way. They
// were previously inlined per method, which is how two of them drifted:
// GroupsMembersResult sorted its outer slice but not the members nested inside
// each entry, and State sorted nothing at all.
//
// All of them are SortStableFunc rather than SortFunc. None of the sort keys is
// guaranteed unique — two resources can share a HashCode, an IPID, or an email —
// and an unstable sort leaves tied elements in an unspecified order, which would
// make the hash of a set containing ties non-deterministic.
//
// Callers must pass a copy. Sorting a caller's slice in place would reorder
// live data as a side effect of computing a hash.

// sortGroupsForHash orders groups by hash code.
func sortGroupsForHash(groups []*Group) {
	slices.SortStableFunc(groups, func(a, b *Group) int {
		return strings.Compare(a.HashCode, b.HashCode)
	})
}

// sortUsersForHash orders users by hash code.
func sortUsersForHash(users []*User) {
	slices.SortStableFunc(users, func(a, b *User) int {
		return strings.Compare(a.HashCode, b.HashCode)
	})
}

// sortMembersByEmailForHash orders members by email address, the only member
// field that can never be empty.
func sortMembersByEmailForHash(members []*Member) {
	slices.SortStableFunc(members, func(a, b *Member) int {
		return strings.Compare(a.Email, b.Email)
	})
}

// sortMembersByIPIDForHash orders members by identity-provider id.
//
// MembersResult has always ordered by IPID rather than by email. The two keys
// are kept distinct deliberately: switching this to email would change every
// MembersResult hash for no benefit.
func sortMembersByIPIDForHash(members []*Member) {
	slices.SortStableFunc(members, func(a, b *Member) int {
		return strings.Compare(a.IPID, b.IPID)
	})
}

// sortGroupsMembersForHash orders the group-members entries by hash code, and
// the members nested inside each entry by email.
//
// The nested pass is the part that was missing. GroupMembers.MarshalBinary walks
// its own Resources in slice order, so member order within a group reached the
// bytes of the enclosing GroupsMembersResult hash even though that hash sorted
// its outer slice. Since internal/core compares GroupsMembersResult hashes to
// decide whether membership needs reconciling, a reordering upstream — Google
// returning a group's members differently, or the SCIM inversion appending them
// in goroutine completion order — produced a spurious full membership sync.
//
// entries and every entry's Resources slice must already be copies.
func sortGroupsMembersForHash(entries []*GroupMembers) {
	for _, entry := range entries {
		if entry == nil {
			continue
		}
		sortMembersByEmailForHash(entry.Resources)
	}

	slices.SortStableFunc(entries, func(a, b *GroupMembers) int {
		return strings.Compare(a.HashCode, b.HashCode)
	})
}
