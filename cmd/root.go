package cmd

import (
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

var ConfigValue util.ConfigType

var rootCmd = &cobra.Command{
	Use:   "todoister",
	Short: "Minimal todoist CLI client",
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
