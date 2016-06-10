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

var cmdFile = &cobra.Command{
	Use:   "file [path to file]",
	Short: "The path to your file",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   fileRun,
}

func fileRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	for _, file := range args {
		wg.Add(1)
		if d, _ := cmd.Flags().GetBool("d"); d {
			go processDecryptFile(file, &wg)
		} else {
			go processCryptFile(file, &wg)
		}
	}
	wg.Wait()
}

func init() {
	cmdFile.Flags().Bool("d", true, "sfs file -d myfile")
	RootCmd.AddCommand(cmdFile)
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
	ioutil.WriteFile(meta.EncodeName, cryptoFile, os.ModePerm)

	wg.Done()
}

func processDecryptFile(filename string, wg *sync.WaitGroup) {
	var key []byte
	file, err := os.Open(filename)
	if err != nil || file == nil {
		panic(err)
	}
	defer file.Close()
	var data = make([]byte, 256)
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}
	restFile := crypt.DecryptByte(key, data)
	ioutil.WriteFile(crypt.DecryptStrBase64(key, file.Name()), restFile, os.ModePerm)
	if string(data) == string(restFile) {
		fmt.Println("Decyphering successful !")
	}
}
