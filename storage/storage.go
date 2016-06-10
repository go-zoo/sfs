package storage

import (
	"errors"
	"os"

	"local/sfs/crypt"
)

type Meta struct {
	OriginalName string
	EncodeName   string
	Length       int64
	Key          []byte
	Platform     Platform
	StorePath    string
}

func NewMeta(key []byte, file *os.File) (Meta, error) {
	fileStat, err := file.Stat()
	if err != nil || fileStat.Size() == 0 {
		return Meta{}, errors.New("Invalid File")
	}
	return Meta{
		OriginalName: file.Name(),
		EncodeName:   crypt.EncryptStrBase64(key, file.Name()),
		Length:       fileStat.Size(),
		Key:          key,
		Platform:     Platform{},
		StorePath:    "",
	}, nil
}

type Platform struct{}
