package repository

import (
	"bytes"
	"context"
	"testing"
)

// FuzzDiskRepositoryGetState fuzzes the state-file loading path.
//
// The state file is the primary untrusted-input surface in this project:
// users can persist it to disk, S3, or another external store and re-load
// it later. A malformed state file must never crash, panic, or hang the
// reconciliation loop — at worst it should return a typed error.
func FuzzDiskRepositoryGetState(f *testing.F) {
	// Seed with well-formed and adversarial inputs so the fuzzer has
	// useful starting points to mutate from.
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"schemaVersion":"1.0.0","codeVersion":"v0.0.0","hashCode":"abc","resources":{"groups":{},"users":{},"groupsMembers":{}}}`))
	f.Add([]byte(``))
	f.Add([]byte(`{"resources":{`))
	f.Add([]byte(`{"resources":{"groups":{"items":2147483647,"resources":null}}}`))
	f.Add([]byte("\x00\x01\x02\x03"))

	f.Fuzz(func(t *testing.T, data []byte) {
		buf := bytes.NewBuffer(data)
		dr, err := NewDiskRepository(buf)
		if err != nil {
			return
		}
		// Errors are acceptable; panics, OOMs, and hangs are not.
		_, _ = dr.GetState(context.Background())
	})
}
