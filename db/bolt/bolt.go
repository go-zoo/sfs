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

// Add add entry to db.
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

// Del delete the entry to db.
func Del(key []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("meta")).Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get find and return the passed key value from the db.
func Get(key []byte) ([]byte, error) {
	var data []byte
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b != nil {
			data = b.Get(key)
			if data != nil {
				return nil
			}
		}
		return errors.New("Not found")
	})

	return data, err
}
