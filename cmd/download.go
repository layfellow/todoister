package cmd

import (
	"github.com/spf13/cobra"
	"todoister/util"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a backup in JSON format from Todoist",
	Run: func(cmd *cobra.Command, args []string) {
		todoistData, _ := util.GetTodoistData()
		hierarchicalData, _ := util.HierarchicalData(todoistData)
		_ = util.WriteHierarchicalData(hierarchicalData, "todoist.json")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
