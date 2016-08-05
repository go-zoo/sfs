package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-zoo/sfs/crypt"
	"github.com/go-zoo/sfs/storage"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(encryptCmd)
	RootCmd.AddCommand(decryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [path to file]",
	Short: "Encrypt the provided files",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   encryptRun,
}

func encryptRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	if len(args) > 0 {
		for _, file := range args {
			wg.Add(1)
			go processCryptFile(file, &wg)
		}
		wg.Wait()
	}
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [path to file]",
	Short: "Decrypt your files",
	Long:  `Write the path of the file you want to decrypt.`,
	Run:   decryptRun,
}

func decryptRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	for _, file := range args {
		wg.Add(1)
		go processDecryptFile(file, &wg)
	}
	wg.Wait()
}

func processCryptFile(filename string, wg *sync.WaitGroup) {
	file, err := os.Open(filename)
	if err != nil || file == nil {
		panic(err)
	}
	defer func() {
		file.Close()
		err = os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()
	key := crypt.GenerateKey(16)
	meta, err := storage.NewMeta(key, file)
	if err != nil {
		panic(err)
	}

	var data = make([]byte, meta.Length)
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	cryptoFile := crypt.EncryptByte(key, data)
	if cryptoFile != nil {
		fmt.Printf("[+] %s Encrypted successfuly !\n", meta.OriginalName)
	}
	err = ioutil.WriteFile(meta.EncodeName, cryptoFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	wg.Done()
}

func processDecryptFile(filename string, wg *sync.WaitGroup) {
	file, err := os.Open(filename)
	if err != nil || file == nil {
		panic(err)
	}
	fs, _ := file.Stat()
	defer func() {
		file.Close()
		err = os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()

	meta, err := storage.FindMeta(filename)
	if err != nil {
		panic(err)
	}

	var data = make([]byte, fs.Size())
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	restFile := crypt.DecryptByte(meta.Key, data)
	err = ioutil.WriteFile(meta.OriginalName, restFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] Decyphering successful !")

	err = storage.DeleteMeta(filename)
	if err != nil {
		fmt.Println(err)
	}

	wg.Done()
}
