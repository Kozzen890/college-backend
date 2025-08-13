package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hash returns the SHA256 hash of a string in hex
func SHA256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
