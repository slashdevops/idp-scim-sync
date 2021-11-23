package hash

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Get return the hash of the value object
// in sha1 format
// https://play.golang.org/p/NAhgOG12YhV
func Get(value interface{}) string {
	if value == nil {
		log.Fatal("hash.Get: value is nil")
	}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
}
