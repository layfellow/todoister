package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

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

	if config.Token == "" {
		if config.Token = viper.GetString("token"); config.Token == "" {
			Die("Missing Todoist token", nil)
		}
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
