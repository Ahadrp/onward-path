package usr

import (
	"crypto/rand"
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
