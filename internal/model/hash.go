package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"slices"
)

// Hash returns a sha256 hash of value pass as argument.
// It panics if value is nil or cannot be gob-encoded, since these
// conditions indicate a programming error in the caller.
func Hash(value any) string {
	if value == nil {
		panic("model: Hash called with nil value")
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		panic(fmt.Sprintf("model: Hash encoding error: %v", err))
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
}

// compactNilPointers returns s with its nil elements removed, preserving the
// order of the remaining ones. When s holds no nil elements it is returned
// unchanged, so hashes computed over well-formed data stay byte-identical to
// those produced by previous releases.
//
// The *Result.SetHashCode methods hash a sorted copy of their Resources slice.
// Both the sort comparators and gob dereference every element, and gob panics
// outright on a nil pointer ("gob: cannot encode nil pointer of type ..."), so
// a single nil entry would abort a sync mid-run. Nil entries can reach these
// slices because idp.buildUser and scim.buildUser return nil to reject a
// malformed record. Such an entry carries no identity, so it is dropped from
// the hash input rather than ordered within it.
func compactNilPointers[T any](s []*T) []*T {
	if !slices.Contains(s, nil) {
		return s
	}

	out := make([]*T, 0, len(s))
	for _, v := range s {
		if v != nil {
			out = append(out, v)
		}
	}

	return out
}
