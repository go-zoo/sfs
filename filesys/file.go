package filesys

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/go-zoo/sfs/crypt"
	"github.com/go-zoo/sfs/storage"
)

// ProcessCryptFile encrypt file
func ProcessCryptFile(path string, filename string, wg *sync.WaitGroup) {
	filepath := path + "/" + filename
	file, err := os.Open(filepath)
	if err != nil || file == nil {
		panic(err)
	}
	defer func() {
		file.Close()
		err = os.Remove(filepath)
		if err != nil {
			panic(err)
		}
	}()
	key := crypt.GenerateKey(32)
	meta, err := storage.NewMeta(key, file)
	if err != nil {
		panic(err)
	}

	var data = make([]byte, meta.Length)
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	cryptoFile := crypt.EncodeWithMaster(key, data)
	if cryptoFile != nil {
		ss := strings.Split(meta.OriginalName, "/")
		fmt.Printf("[+] %s Encrypted successfuly !\n", ss[len(ss)-1])
	}
	err = ioutil.WriteFile(meta.Path+"/"+meta.EncodeName, cryptoFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	wg.Done()
}

// ProcessDecryptFile decrypt file
func ProcessDecryptFile(path string, name string, wg *sync.WaitGroup) {
	ss := strings.Split(name, "/")
	filename := ss[len(ss)-1]

	meta, err := storage.FindMeta(filename)
	if err != nil {
		fmt.Printf("[!] %s doesn't seems to exists or be encrypted\n", filename)
		wg.Done()
		return
	}

	file, err := os.Open(meta.Path + "/" + filename)
	if err != nil || file == nil {
		panic(err)
	}

	fs, _ := file.Stat()

	defer func() {
		file.Close()
		err = os.Remove(meta.Path + "/" + filename)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	var data = make([]byte, fs.Size())
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	restFile := crypt.DecodeWithMaster(meta.Key, data)
	err = ioutil.WriteFile(meta.Path+"/"+meta.OriginalName, restFile, os.FileMode(meta.FileMode))
	if err != nil {
		panic(err)
	}

	fmt.Printf("[+] Decrypted %s successfuly !\n", meta.OriginalName)

	err = storage.DeleteMeta(filename)
	if err != nil {
		fmt.Println(err)
	}
}
