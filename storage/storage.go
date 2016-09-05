package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	pb "github.com/go-zoo/sfs/proto"
	"github.com/go-zoo/sfs/storage/platforms"
)

// Meta store informations about a file
type Meta struct {
	OriginalName string             `json:"orgname"`
	EncodeName   string             `json:"encname"`
	Path         string             `json:"path"`
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
	pathTokens := strings.Split(file.Name(), "/")

	m := pb.Meta{
		OriginalName: fileStat.Name(),
		EncodeName:   fmt.Sprintf("%x", h.Sum(nil)),
		Path:         strings.Join(pathTokens[:len(pathTokens)-1], "/"),
		Length:       fileStat.Size(),
		FileMode:     uint32(fileStat.Mode()),
		Key:          key,
		StorePath:    "",
	}

	fmt.Println(m.String())

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
