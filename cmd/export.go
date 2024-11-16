package cmd

import (
	"github.com/spf13/cobra"
	"todoister/util"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export projects in JSON or YAML format",
	Run: func(cmd *cobra.Command, args []string) {
		todoistData, _ := util.GetTodoistData()
		hierarchicalData, _ := util.HierarchicalData(todoistData)
		_ = util.WriteHierarchicalData(hierarchicalData, "todoist.json")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
