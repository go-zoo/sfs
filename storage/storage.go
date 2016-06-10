package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"local/sfs/crypt"
)

type Meta struct {
	OriginalName string   `json:"orgname"`
	EncodeName   string   `json:"encname"`
	Length       int64    `json:"length"`
	Key          []byte   `json:"key"`
	Platform     Platform `json:"platform"`
	StorePath    string   `json:"store"`
}

var Metas map[string]Meta

func init() {
	Metas = make(map[string]Meta)
	data, err := ioutil.ReadFile("sfs.conf")
	if err != nil {
		d, _ := json.Marshal(make(map[string]interface{}))
		err = ioutil.WriteFile("sfs.conf", d, os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else {
		err = json.Unmarshal(data, &Metas)
		if err != nil {
			panic(err)
		}
	}
}

func NewMeta(key []byte, file *os.File) (Meta, error) {
	fileStat, err := file.Stat()
	if err != nil || fileStat.Size() == 0 {
		return Meta{}, errors.New("Invalid File")
	}
	m := Meta{
		OriginalName: file.Name(),
		EncodeName:   crypt.EncryptStrBase64(key, file.Name()),
		Length:       fileStat.Size(),
		Key:          key,
		Platform:     Platform{},
		StorePath:    "",
	}
	Metas[m.EncodeName] = m
	writeConf()
	return m, nil
}

func FindMeta(encodeName string) (Meta, error) {
	if Metas[encodeName].EncodeName != "" {
		fmt.Println(Metas[encodeName].OriginalName)
		return Metas[encodeName], nil
	}
	return Meta{}, errors.New("Meta not found")
}

func writeConf() {
	data, err := json.MarshalIndent(Metas, "", "  ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("sfs.conf", data, os.ModePerm)
}

type Platform struct{}
