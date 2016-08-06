package cmd

import (
	"sync"

	"github.com/go-zoo/sfs/filesys"

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
			go filesys.ProcessCryptFile("", file, &wg)
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
		go filesys.ProcessDecryptFile("", file, &wg)
	}
	wg.Wait()
}
