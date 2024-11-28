package cmd

import (
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	VersionText = "Minimal Todoist CLI client"
)

var Version = "DEV"

var ConfigValue util.ConfigType

var RootCmd = &cobra.Command{
	Use:               "todoister",
	Version:           Version,
	Short:             VersionText,
	DisableAutoGenTag: true,
}

func initAll() {
	util.InitConfig(&ConfigValue)
	util.InitLogger(ConfigValue.Log.Name)
}

func init() {
	cobra.OnInitialize(initAll)
	RootCmd.PersistentFlags().StringVarP(&ConfigValue.Token, "token", "t", "", "Override Todoist token")
	RootCmd.SetVersionTemplate(`{{printf "` + VersionText + ` v%s\n" .Version }}`)
}

func Execute() {
	_ = RootCmd.Execute()
}
