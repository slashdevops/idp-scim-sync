// Package model holds the domain types shared by every other package, and the
// hashing scheme that drives change detection.
//
// # Builders are mandatory
//
// Never construct a value in this package directly. Use the builder chain:
//
//	group := model.GroupBuilder().
//	    WithIPID("google-group-id").
//	    WithName("developers").
//	    WithEmail("developers@example.com").
//	    Build()
//
// Build is what sets Items on result types and computes HashCode. A
// hand-constructed struct carries a zero hash and therefore compares unequal to
// everything, which silently turns a no-op sync into a full rewrite.
//
// # Hashing
//
// HashCode is a SHA-256 over a gob encoding of the value, and each type in this
// package hand-writes MarshalBinary to control exactly which fields participate.
// Two consequences follow:
//
//   - SCIMID is excluded on purpose. It is assigned by AWS, not by the identity
//     provider, so including it would make every first-run comparison differ.
//   - Those MarshalBinary methods are a wire format. Changing field order,
//     adding a field or removing one changes every hash, which invalidates every
//     state file already deployed and forces a full re-sync.
//     golden_hash_test.go pins the current values; if it fails, the format
//     moved.
//
// Result types sort a copy of their resources before hashing, so a container
// hash does not depend on the order its elements arrived in — which matters
// because the identity-provider fan-outs build their slices from maps.
//
// # Optional attributes
//
// [SyncFieldSet] selects which optional user attributes take part in a sync. An
// empty set means "all of them", which is the backward-compatible default.
//
// See docs/State-File-example.md for the serialized form and
// docs/Architecture.md for how these types flow through a sync.
package model
