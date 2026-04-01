package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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
