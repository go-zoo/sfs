package crypt

import (
	"crypto/rand"
	"os"
)

var MasterKey []byte

func init() {
	MasterKey = []byte(os.Getenv("SFSMASTERKEY"))
}

// Generate a random key of the provided length
func GenerateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

// func GenerateWithMasterKey(length int) []byte {
// 	key := GenerateKey(length)
// 	return EncryptByte([]byte(masterKey), key)[:16]
// }
//
// func GetKeyWithMasterKey(key []byte) []byte {
// 	return DecryptByte([]byte(masterKey), key)
// }
