package util

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func InitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Preferred location is ~/.config/todoister/config.json
	configPath := filepath.Join(homeDir, ".config", "todoister", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Fallback location is ~/.todoister.json
		configPath = filepath.Join(homeDir, ".todoister.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("Configuration file not found in either location")
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	TodoistToken = viper.GetString("token")
	if TodoistToken == "" {
		log.Fatalf("Token not found in configuration file")
	}
}
