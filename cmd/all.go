package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-zoo/sfs/filesys"
	"github.com/spf13/cobra"
)

type fileProcessor func(string, string, *sync.WaitGroup)

func init() {
	encodeCmd.AddCommand(encryptAllCmd)
	decodeCmd.AddCommand(decryptAllCmd)
}

var encryptAllCmd = &cobra.Command{
	Use:   "all files in the current directory will be encrypted.",
	Short: "Encrypt all files in the directory",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   encryptAllRun,
}

func encryptAllRun(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	ProcessFileTree(wd, filesys.ProcessCryptFile)
}

var decryptAllCmd = &cobra.Command{
	Use:   "all files in the current directory will be decrypted.",
	Short: "Decrypt all files in the directory",
	Long:  `Write the path of the file you want to encrypt.`,
	Run:   decryptAllRun,
}

func decryptAllRun(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	ProcessFileTree(wd, filesys.ProcessDecryptFile)
}

// ProcessFileTree look recursively for files in directory
func ProcessFileTree(dirname string, process fileProcessor) {
	var wg sync.WaitGroup
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(files) > 0 {
		for _, file := range files {
			if file.IsDir() {
				ProcessFileTree(dirname+"/"+file.Name(), process)
			} else {
				wg.Add(1)
				go process(dirname, file.Name(), &wg)
			}
		}
		wg.Wait()
	}
}
