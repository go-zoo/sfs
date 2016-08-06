package crypt

import (
	"crypto/rand"
	"fmt"
	"os"
)

var MasterKey []byte

func init() {
	mk := os.Getenv("SFSMASTERKEY")
	if mk != "" {
		MasterKey = []byte(mk)
		return
	}
	MasterKey = GenerateKey(32)
	fmt.Println("[!] MasterKey not found !")
	fmt.Println("[!] SFS have generated one for you.")
	fmt.Printf("[!] Add (export SFSMASTERKEY=%s)\n", MasterKey)
}

// Generate a random key of the provided length
func GenerateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

func EncodeWithMaster(key []byte, data []byte) []byte {
	return EncryptByte(MasterKey, EncryptByte(key, data))
}

func DecodeWithMaster(key []byte, data []byte) []byte {
	return DecryptByte(key, DecryptByte(MasterKey, data))
}

// func GenerateWithMasterKey(length int) []byte {
// 	key := GenerateKey(length)
// 	return EncryptByte([]byte(masterKey), key)[:16]
// }
//
// func GetKeyWithMasterKey(key []byte) []byte {
// 	return DecryptByte([]byte(masterKey), key)
// }
