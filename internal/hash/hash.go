package hash

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"

	"github.com/mitchellh/hashstructure/v2"
)

// Sha256 returns the sha256 hash of the value object
// thanks to: https://blog.8bitzen.com/posts/22-08-2019-how-to-hash-a-struct-in-go
func Sha256V1(value interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", value)))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Get(value interface{}) string {
	hash, err := hashstructure.Hash(value, hashstructure.FormatV2, &hashstructure.HashOptions{Hasher: fnv.New64()})
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash)
}
