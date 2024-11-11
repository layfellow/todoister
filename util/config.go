package util

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func InitConfig() {
	if TodoistToken != "" { // Already set at command line via --token option.
		return
	}

	if err := viper.BindEnv("token", "TODOIST_TOKEN"); err != nil {
		log.Fatalf("Error binding environment variable: %v", err)
	}
	if TodoistToken = viper.GetString("token"); TodoistToken != "" {
		return
	}

	viper.SetConfigType("json")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	preferredConfigPath := filepath.Join(homeDir, ".config", "todoister", "config.json")
	alternativeConfigPath := filepath.Join(homeDir, ".todoister.json")

	if _, err := os.Stat(preferredConfigPath); err == nil {
		viper.SetConfigFile(preferredConfigPath)
	} else if _, err := os.Stat(alternativeConfigPath); err == nil {
		viper.SetConfigFile(alternativeConfigPath)
	} else {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
	if TodoistToken = viper.GetString("token"); TodoistToken == "" {
		log.Fatal("Token must be set via config file, environment variable, or command line argument")
	}
}
