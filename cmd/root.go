package cmd

import (
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	VersionText = "Minimal Todoist CLI client"

	rootLong = `Todoister is a simple Todoist CLI client written in Go.
`
)

var Version = "DEV"

var ConfigValue util.ConfigType

var RootCmd = &cobra.Command{
	Use:               "todoister command [arguments]",
	Version:           Version,
	Short:             VersionText,
	Long:              rootLong,
	DisableAutoGenTag: true,
}

func initAll() {
	util.SchemaVersion = Version
	util.InitConfig(&ConfigValue)
	util.InitLogger(ConfigValue.Name)
}

func init() {
	cobra.OnInitialize(initAll)
	RootCmd.PersistentFlags().StringVarP(&ConfigValue.Token, "token", "t", "",
		"use <string> as Todoist API token")
	RootCmd.SetVersionTemplate(`{{printf "` + VersionText + ` v%s\n" .Version }}`)
	RootCmd.SetHelpFunc(util.CustomHelpFunc)
}

func Execute() {
	_ = RootCmd.Execute()
}
