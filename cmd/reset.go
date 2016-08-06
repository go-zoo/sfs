package cmd

import (
	"fmt"
	"os"

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

	err = os.Unsetenv("SFSMASTERKEY") // not working ...
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[+] Configurations have been reset!")
}
