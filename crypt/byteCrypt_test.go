package crypt

import (
	"reflect"
	"testing"
)

func TestByteCrypt(t *testing.T) {
	rawData := []byte{0x23, 0x34, 0x45, 0x88, 0x33, 0x76, 0x21, 0x91}
	key := GenerateKey(32)

	cryptData := EncryptByte(key, rawData)
	if cryptData == nil {
		t.Fatal("Encrypted data is nil")
	}
	data := DecryptByte(key, cryptData)
	if data == nil {
		t.Fatal("Decrypted data is nil")
	}

	if !reflect.DeepEqual(data, rawData) {
		t.Logf("Raw: %d", rawData)
		t.Logf("Data: %d", data)
		t.Fatal("Not equal")
	}
}

func TestEncryptByte(t *testing.T) {
	rawData := []byte{0x23, 0x34, 0x45, 0x88, 0x33, 0x76, 0x21, 0x91}
	key := GenerateKey(32)

	cryptData := EncryptByte(key, rawData)
	if cryptData == nil {
		t.Fatal("Encrypted data is nil")
	}
}
