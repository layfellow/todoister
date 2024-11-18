package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"todoister/util"
)

var ConfigValue util.ConfigType

var rootCmd = &cobra.Command{
	Use:   "todoister",
	Short: "Minimal todoist CLI client",
	Long:  "Lorem ipsum dolor sit amet", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: todoister")
	},
}

func initAll() {
	util.InitConfig(&ConfigValue)
	util.InitLogger(ConfigValue.Log.Name)
}

func init() {
	cobra.OnInitialize(initAll)
	rootCmd.PersistentFlags().StringVarP(&ConfigValue.Token, "token", "t", "", "Override Todoist token.")
}

func Execute() {
	_ = rootCmd.Execute()
}
