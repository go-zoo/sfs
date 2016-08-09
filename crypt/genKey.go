package crypt

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/go-zoo/sfs/db/bolt"
)

// MasterKey is the used to encrypt everything ...
var MasterKey []byte

func init() {
	if !checkMasterKey() {
		var password string
		fmt.Println("[!] MasterKey not found !")
		fmt.Printf("[+] Enter a password : ")
		fmt.Scan(&password)
		err := setPassword(password)
		if err != nil {
			panic(err)
		}
		fmt.Println("[!] SFS have generated one for you.")
	}
}

func checkMasterKey() bool {
	mk, err := bolt.Get([]byte("MASTERKEY"))
	if err != nil || mk == nil {
		MasterKey = GenerateKey(32)
		err = bolt.Add([]byte("MASTERKEY"), MasterKey)
		if err != nil {
			panic(err)
		}
		return false
	}
	MasterKey = mk
	return true
}

func setPassword(password string) error {
	h := md5.New()
	io.WriteString(h, password)

	return bolt.Add([]byte("password"), h.Sum(nil))
}

// GenerateKey a random key of the provided length
func GenerateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

// EncodeWithMaster encode a key with MasterKey encoding it.
func EncodeWithMaster(key []byte, data []byte) []byte {
	return EncryptByte(MasterKey, EncryptByte(key, data))
}

// DecodeWithMaster decode the key with the MasterKey
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
