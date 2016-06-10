package crypt

import "crypto/rand"

// Generate a random key of the provided length
func GenerateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}
