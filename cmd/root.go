package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"todoister/util"
)

var rootCmd = &cobra.Command{
	Use:   "todoister",
	Short: "Minimal todoist CLI client",
	Long:  "Lorem ipsum dolor sit amet", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: todoister")
	},
}

func init() {
	cobra.OnInitialize(util.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&util.TodoistToken, "token", "t", "", "Override Todoist token.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
