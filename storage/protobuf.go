package storage

import (
	"io/ioutil"
	"os"

	pb "github.com/go-zoo/sfs/proto"

	"github.com/golang/protobuf/proto"
)

func ReadProtoFile(filename string) (*pb.Meta, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var metas *pb.Meta
	err = proto.Unmarshal(data, metas)
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func WriteProtoFile(pbm proto.Message, filename string) error {
	data, err := proto.Marshal(pbm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
