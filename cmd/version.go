package cmd

import (
	"fmt"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	versionLong = `Print the Todoister version number.
`
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  versionLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", VersionText, Version)
	},
}

func init() {
	versionCmd.SetHelpFunc(util.CustomHelpFunc)
	RootCmd.AddCommand(versionCmd)
}
