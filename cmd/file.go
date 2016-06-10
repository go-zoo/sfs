package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"local/sfs/crypt"
	"local/sfs/storage"

	"github.com/spf13/cobra"
)

var cryptCmd = &cobra.Command{
	Use:   "crypt [path to file]",
	Short: "The path to your file",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   cryptRun,
}

func cryptRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	for _, file := range args {
		wg.Add(1)
		go processCryptFile(file, &wg)
	}
	wg.Wait()
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [path to file]",
	Short: "The path to your file",
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

func init() {
	RootCmd.AddCommand(cryptCmd)
	RootCmd.AddCommand(decryptCmd)
}

func processCryptFile(filename string, wg *sync.WaitGroup) {
	file, err := os.Open(filename)
	if err != nil || file == nil {
		panic(err)
	}
	defer file.Close()
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

	cryptoFile := crypt.EncryptByte(key, data)
	if cryptoFile != nil {
		fmt.Printf("[+] %s Encrypted successfuly !\n", meta.OriginalName)
	}
	err = ioutil.WriteFile(meta.EncodeName, cryptoFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
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
	defer file.Close()
	meta, err := storage.FindMeta(filename)
	if err != nil {
		panic(err)
	}

	var data = make([]byte, meta.Length)
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data Size:", len(data))

	restFile := crypt.DecryptByte(meta.Key, data)
	err = ioutil.WriteFile(meta.OriginalName, restFile, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Decyphering successful !")
	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
	wg.Done()
}
