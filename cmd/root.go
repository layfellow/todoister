package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"todoister/util"
)

func init() {
	cobra.OnInitialize(util.InitConfig)
}

var rootCmd = &cobra.Command{
	Use:   "todoister",
	Short: "Minimal todoist CLI client",
	Long:  "Lorem ipsum dolor sit amet", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: todoister")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
