package bolt

import (
	"errors"
	"os"

	"github.com/boltdb/bolt"
	"github.com/kardianos/osext"
)

const dbname = "sfs.dat"

var db *bolt.DB

func init() {
	exp, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	db, err = bolt.Open(exp+"/"+dbname, os.ModePerm, nil)
	if err != nil {
		panic(err)
	}
}

func Add(key []byte, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})

	return err
}

func Get(key []byte) ([]byte, error) {
	var data []byte
	err := db.Update(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("meta")).Get(key)
		if data != nil {
			return nil
		}
		return errors.New("Not found")
	})

	return data, err
}
