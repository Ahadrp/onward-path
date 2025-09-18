package usr

import (
	"crypto/rand"
    "crypto/sha256"
    "encoding/hex"
)

// Generate secure random token
func GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Convert a string to sha256
func SHA256(s string) string {
    h := sha256.Sum256([]byte(s))
    return hex.EncodeToString(h[:])
}
