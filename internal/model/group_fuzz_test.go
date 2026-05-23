package model

import (
	"testing"
)

// FuzzGroupsResultUnmarshalBinary fuzzes the gob deserialization of a
// GroupsResult. The decoder loops `Items` times calling dec.Decode, so an
// attacker who controls the encoded blob can ask for an enormous number of
// items and force the process to allocate aggressively before failing. The
// fuzzer should surface panics, hangs, or unbounded allocations.
func FuzzGroupsResultUnmarshalBinary(f *testing.F) {
	// Seed with a valid round-tripped encoding.
	seed := GroupsResult{
		Items: 2,
		Resources: []*Group{
			{IPID: "1", Name: "a", Email: "a@example.com"},
			{IPID: "2", Name: "b", Email: "b@example.com"},
		},
	}
	if data, err := seed.MarshalBinary(); err == nil {
		f.Add(data)
	}
	f.Add([]byte{})
	f.Add([]byte{0x00})
	f.Add([]byte{0xff, 0xff, 0xff, 0xff})

	f.Fuzz(func(t *testing.T, data []byte) {
		var gr GroupsResult
		_ = gr.UnmarshalBinary(data)
	})
}
