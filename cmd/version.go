package cmd

import (
	"fmt"
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the Todoister version number.\n",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", VersionText, Version)
	},
}

func init() {
	versionCmd.SetHelpFunc(util.CustomHelpFunc)
	RootCmd.AddCommand(versionCmd)
}
