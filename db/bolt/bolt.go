package bolt

import (
	"errors"
	"os"

	"github.com/boltdb/bolt"
)

const dbname = "sfs.dat"

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open(dbname, os.ModePerm, nil)
	if err != nil {
		panic(err)
	}
}

func Add(key []byte, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("meta")).Put(key, value)
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
