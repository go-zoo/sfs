package cmd

import "github.com/spf13/cobra"

// RootCmd is the base command of sfs
var RootCmd = &cobra.Command{
	Use:   "sfs [command] files name",
	Short: "Easy file encryption",
	Long: `
    Safe File Storage is a tool to easily
    encrypt and store data
    in cloud platform (Google Drive, One Drive, etc ...)`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
