package cmd

import (
	"fmt"

	"github.com/go-zoo/sfs/storage"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [filename]",
	Short: "Show infos about a encrypted file",
	Long:  `Write the name of the file you want to get infos of.`,
	Run:   infoRun,
}

func infoRun(cmd *cobra.Command, args []string) {
	m := storage.FindMeta(args[0])
	if err != nil {
		fmt.Println(args[0], "not found")
		return
	}
	m.PrintMeta()
}
