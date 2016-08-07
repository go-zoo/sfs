package storage

import (
	"io/ioutil"
	"os"

	"github.com/go-zoo/sfs/db/bolt"
	pb "github.com/go-zoo/sfs/proto"

	"github.com/golang/protobuf/proto"
)

func ReadFromFile(filename string) (*pb.Meta, error) {
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

func ReadFromDb(key string) (proto.Message, error) {
	data, err := bolt.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	//data = crypt.DecryptByte(crypt.MasterKey, data)

	meta := &pb.Meta{}
	err = proto.Unmarshal(data, meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func WriteToFile(pbm proto.Message, filename string) error {
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

func WriteToDb(key string, pbm proto.Message) error {
	data, err := proto.Marshal(pbm)
	if err != nil {
		return err
	}
	//data = crypt.EncryptByte(crypt.MasterKey, data)

	return bolt.Add([]byte(key), data)
}

func DelFromDb(key string) error {
	return bolt.Del([]byte(key))
}
