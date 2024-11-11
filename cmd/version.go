package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Minimal Todoist CLI client v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
