package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
)

// Hash returns a sha256 hash of value pass as argument
func Hash(value interface{}) string {
	if value == nil {
		slog.Error("value is nil")
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		slog.Error("error encoding value")
		os.Exit(1)
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
}
