package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Hash returns a sha256 hash of value pass as argument
func Hash(value interface{}) string {
	if value == nil {
		log.Fatal("value is nil")
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
}
