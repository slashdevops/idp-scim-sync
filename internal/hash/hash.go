package hash

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"hash/fnv"

	"github.com/mitchellh/hashstructure/v2"
)

// Get return the hash of the value object
func GetV2(value interface{}) string {
	hash, err := hashstructure.Hash(value, hashstructure.FormatV2, &hashstructure.HashOptions{Hasher: fnv.New64()})
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash)
}

// Get return the hash of the value object
// in sha1 format
func Get(value interface{}) string {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
}
