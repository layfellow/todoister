package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// Default configuration file name and extension.
const (
	ConfigFile    = "config"
	ConfigFileExt = "toml"
)

type Log struct {
	Name string
}

type Export struct {
	Path   string
	Format string
	Depth  int
}

type ConfigType struct {
	Token string
	Log
	Export
}

// InitConfig initializes the configuration from the configuration file and environment variables.
//   - config: a pointer to a ConfigType struct to store the configuration.
func InitConfig(config *ConfigType) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		Die("Error getting user home directory", err)
	}

	viper.SetEnvPrefix("TODOIST")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName(ConfigFile)
	viper.SetConfigType(ConfigFileExt)

	// If $XDG_CONFIG_HOME is set, use it
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		viper.AddConfigPath(filepath.Join(xdgConfigHome, Prog))
	} else {
		// Otherwise use default XDG-compliant ~/.config
		viper.AddConfigPath(filepath.Join(homeDir, ".config", Prog))
	}

	if err = viper.ReadInConfig(); err != nil {
		// Fall back to more traditional ~/.todoister.toml
		viper.SetConfigFile(filepath.Join(homeDir, fmt.Sprintf(".%s.%s", Prog, ConfigFileExt)))
		_ = viper.ReadInConfig()
	}

	// Config.Token may have been set already by the -t, --token flag.
	if config.Token == "" {
		config.Token = viper.GetString("token")
	}
	if config.Log.Name == "" {
		config.Log.Name, _ = ExpandPath(viper.GetString("log.name"))
	}
	if config.Export.Path == "" {
		config.Export.Path, _ = ExpandPath(viper.GetString("export.path"))
	}
	if config.Export.Format == "" {
		config.Export.Format, _ = ExpandPath(viper.GetString("export.format"))
	}
	if config.Export.Depth == 0 {
		config.Export.Depth = viper.GetInt("export.depth")
	}
}
