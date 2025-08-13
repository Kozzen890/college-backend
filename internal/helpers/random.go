package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a secure random string of n bytes (hex encoded)
func GenerateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
