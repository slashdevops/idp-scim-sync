package hash

import (
	"crypto/sha256"
	"fmt"
)

// Sha256 returns the sha256 hash of the value object
// thanks to: https://blog.8bitzen.com/posts/22-08-2019-how-to-hash-a-struct-in-go
func Sha256(value interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", value)))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
