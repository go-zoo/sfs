package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-zoo/sfs/filesys"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(encodeCmd)
	RootCmd.AddCommand(decodeCmd)
}

var encodeCmd = &cobra.Command{
	Use:   "encode [path to file]",
	Short: "Encrypt the provided files",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   encryptRun,
}

func encryptRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) > 0 {
		for _, file := range args {
			wg.Add(1)
			go filesys.ProcessCryptFile(wd, file, &wg)
		}
		wg.Wait()
	}
}

var decodeCmd = &cobra.Command{
	Use:   "decode [path to file]",
	Short: "Decrypt your files",
	Long:  `Write the path of the file you want to decrypt.`,
	Run:   decryptRun,
}

func decryptRun(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range args {
		wg.Add(1)
		go filesys.ProcessDecryptFile(wd, file, &wg)
	}
	wg.Wait()
}
