package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"

	pb "github.com/go-zoo/sfs/proto"
	"github.com/go-zoo/sfs/storage/platforms"
)

// Meta store informations about a file
type Meta struct {
	OriginalName string             `json:"orgname"`
	EncodeName   string             `json:"encname"`
	Length       int64              `json:"length"`
	FileMode     uint32             `json:"filemode"`
	Key          []byte             `json:"key"`
	Platform     platforms.Platform `json:"platform"`
	StorePath    string             `json:"store"`
}

var conf string

// NewMeta generate a Meta struct base on the file
func NewMeta(key []byte, file *os.File) (pb.Meta, error) {
	fileStat, err := file.Stat()
	if err != nil || fileStat.Size() == 0 {
		return pb.Meta{}, errors.New("Invalid File")
	}
	h := md5.New()
	io.WriteString(h, file.Name())

	m := pb.Meta{
		OriginalName: file.Name(),
		EncodeName:   fmt.Sprintf("%x", h.Sum(nil)),
		Length:       fileStat.Size(),
		FileMode:     uint32(fileStat.Mode()),
		Key:          key,
		//Platform:     Platform{},
		StorePath: "",
	}

	return m, WriteToDb(m.EncodeName, &m)
}

// PrintMeta pretty print a file info
// func (m *Meta) PrintMeta() {
// 	fmt.Printf("# Filename : %s\n", m.OriginalName)
// 	fmt.Printf("# Encoded Name : %s\n", m.EncodeName)
// 	fmt.Printf("# Size : %d\n", m.Length)
// }

// FindMeta retrieve Meta info from the sfs.conf file
func FindMeta(encodeName string) (*pb.Meta, error) {
	mess, err := ReadFromDb(encodeName)
	if err != nil {
		return &pb.Meta{}, err
	}
	return mess.(*pb.Meta), nil
}

// DeleteMeta remove old encryption data
func DeleteMeta(encodeName string) error {
	return DelFromDb(encodeName)
}
