package encryption

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hash(value, key string) string {
	// Concatenate the password and secret key
	data := []byte(value + key)

	// Calculate the SHA-256 hash
	hash := sha256.New()
	hash.Write(data)
	hashed := hash.Sum(nil)

	// Convert the hash to a hex-encoded string
	return hex.EncodeToString(hashed)
}
