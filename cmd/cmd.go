package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "name of application",
	Short: "Short description",
	Long: `Longer description..
            feel free to use a few lines here.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
