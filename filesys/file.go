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
	err = ioutil.WriteFile(path+"/"+meta.EncodeName, cryptoFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	wg.Done()
}

// ProcessDecryptFile decrypt file
func ProcessDecryptFile(path string, filename string, wg *sync.WaitGroup) {
	filepath := path + "/" + filename

	file, err := os.Open(filepath)
	if err != nil || file == nil {
		panic(err)
	}

	fs, _ := file.Stat()
	meta, err := storage.FindMeta(filename)
	if err != nil {
		fmt.Printf("[!] %s doesn't seems to exists or be encrypted\n", fs.Name())
		wg.Done()
		return
	}

	defer func() {
		file.Close()
		err = os.Remove(filepath)
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
	err = ioutil.WriteFile(meta.OriginalName, restFile, os.FileMode(meta.FileMode))
	if err != nil {
		panic(err)
	}
	ss := strings.Split(meta.OriginalName, "/")
	fmt.Printf("[+] Decrypted %s successfuly !\n", ss[len(ss)-1])

	err = storage.DeleteMeta(filename)
	if err != nil {
		fmt.Println(err)
	}
}
