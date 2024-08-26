package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Configuration struct {
	TOKEN string
}

var cfgFile string
var Config Configuration

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

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".todo")
	}

	viper.SetEnvPrefix("TODOIST")
	viper.AutomaticEnv()
	Config.TOKEN = viper.GetString("TOKEN")

	if Config.TOKEN == "" {
		if err := viper.ReadInConfig(); err == nil {
			// fmt.Println("Config file:", viper.ConfigFileUsed())
			err := viper.Unmarshal(&Config)
			if err != nil {
				fmt.Printf("Error reading config file, %v\n", err)
				os.Exit(1)
			}
		}
	}

	if Config.TOKEN == "" {
		fmt.Println("Missing TOKEN, use --help")
		os.Exit(1)
	}
}
