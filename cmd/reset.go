package cmd

import (
	"fmt"
	"os"

	"github.com/go-zoo/sfs/db/bolt"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configurations",
	Long:  `Erase the env MasterKey and sfs registery.`,
	Run:   resetRun,
}

func resetRun(cmd *cobra.Command, args []string) {
	exp, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	err = os.Remove(exp + "/sfs.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bolt.Del([]byte("MASTERKEY"))
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] Configurations have been reset!")
}
