package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-zoo/sfs/db/bolt"
	"github.com/go-zoo/sfs/filesys"

	"github.com/spf13/cobra"
)

var recursive bool
var iteration int

func init() {
	encodeCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "-r")
	encodeCmd.PersistentFlags().IntVarP(&iteration, "iteration", "i", 1, "-i [number of iteration]")

	decodeCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "-r")
	decodeCmd.PersistentFlags().IntVarP(&iteration, "iteration", "i", 1, "-i [number of iteration]")

	RootCmd.AddCommand(encodeCmd)
	RootCmd.AddCommand(decodeCmd)
}

func setPassword() {

}

func checkPassword() bool {
	var password string
	fmt.Printf("[+] Enter your password : ")
	fmt.Scan(&password)

	d, err := bolt.Get([]byte("password"))
	if err != nil {
		return false
	}

	encoPassword := fmt.Sprintf("%x", d)

	h := md5.New()
	io.WriteString(h, password)
	password = fmt.Sprintf("%x", h.Sum(nil))

	if encoPassword != password {
		return false
	}

	return true
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
		for i := 0; i < iteration; i++ {
			for _, file := range args {
				wg.Add(1)
				go filesys.ProcessCryptFile(wd, file, &wg)
			}
			wg.Wait()
		}
	}
}

var decodeCmd = &cobra.Command{
	Use:   "decode [path to file]",
	Short: "Decrypt your files",
	Long:  `Write the path of the file you want to decrypt.`,
	Run:   decryptRun,
}

func decryptRun(cmd *cobra.Command, args []string) {
	if !checkPassword() {
		fmt.Println("[!] Invalid Password")
		return
	}
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
