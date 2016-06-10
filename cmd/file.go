package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"local/sfs/crypt"

	"github.com/spf13/cobra"
)

var cmdFile = &cobra.Command{
	Use:   "file [path to file]",
	Short: "The path to your file",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   fileRun,
}

func fileRun(cmd *cobra.Command, args []string) {
	orgFileName := args[0]
	originalFile, err := ioutil.ReadFile(orgFileName)
	if err != nil {
		panic(err)
	}

	key := crypt.GenerateKey(32)
	cryptName := "encrypt.bin"

	cryptoFile := crypt.EncryptByte(key, originalFile)
	if cryptoFile != nil {
		fmt.Println("Encryption successful !")
	}
	ioutil.WriteFile(cryptName, cryptoFile, os.ModePerm)

	file := crypt.DecryptByte(key, cryptoFile)
	ioutil.WriteFile("decrypt.jpg", file, os.ModePerm)
	if string(file) == string(originalFile) {
		fmt.Println("Decyphering successful !")
	}
}

func init() {
	RootCmd.AddCommand(cmdFile)
}
