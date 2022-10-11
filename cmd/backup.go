package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Download the latest backup from Todoist",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Config.TOKEN)
	},
}
