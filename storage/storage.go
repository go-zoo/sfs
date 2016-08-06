package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-zoo/sfs/crypt"

	"github.com/kardianos/osext"
)

// Meta store informations about a file
type Meta struct {
	OriginalName string   `json:"orgname"`
	EncodeName   string   `json:"encname"`
	Length       int64    `json:"length"`
	Key          []byte   `json:"key"`
	Platform     Platform `json:"platform"`
	StorePath    string   `json:"store"`
}

var conf string

var metas = struct {
	sync.RWMutex
	m map[string]Meta
}{
	m: make(map[string]Meta),
}

func init() {
	exp, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	conf = exp + "/sfs.conf"

	data, err := ioutil.ReadFile(conf)
	if err != nil {
		d, _ := json.Marshal(make(map[string]interface{}))
		err = ioutil.WriteFile(conf, d, os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else {
		err = json.Unmarshal(data, &metas.m)
		if err != nil {
			panic(err)
		}
	}
}

// NewMeta generate a Meta struct base on the file
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
	metas.Lock()
	metas.m[m.EncodeName] = m
	metas.Unlock()
	writeConf()
	return m, nil
}

// PrintMeta pretty print a file info
func (m *Meta) PrintMeta() {
	fmt.Printf("# Filename : %s\n", m.OriginalName)
	fmt.Printf("# Encoded Name : %s\n", m.EncodeName)
	fmt.Printf("# Size : %d\n", m.Length)
	//fmt.Println("# Platform : %s", "")
}

// FindMeta retrieve Meta info from the sfs.conf file
func FindMeta(encodeName string) (Meta, error) {
	if metas.m[encodeName].Key != nil {
		return metas.m[encodeName], nil
	}
	return Meta{}, errors.New("Meta not found")
}

// DeleteMeta remove old encryption data
func DeleteMeta(encodeName string) error {
	metas.Lock()
	if metas.m[encodeName].EncodeName != "" {
		delete(metas.m, metas.m[encodeName].EncodeName)
		metas.Unlock()
		writeConf()
		return nil
	}
	return errors.New("Meta not found")
}

func writeConf() {
	metas.Lock()
	data, err := json.MarshalIndent(metas.m, "", "  ")
	metas.Unlock()
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(conf, data, os.ModePerm)
}

// Platform is not implemented
type Platform struct{}
