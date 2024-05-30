package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Set globals defaults
	workingPath, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting working path", "err", err)
	}

	defaultConfigs()

	// Default config file in executable route ./config/config.yml
	defPath := filepath.Join(workingPath, "config", "config.yml")
	viper.SetDefault("config.file", defPath)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ltpapi")
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigFile(viper.GetString("config.file"))
	slog.Info("Reading config", "filename", viper.GetString("config.file"))
	err = viper.ReadInConfig()

	if err != nil {
		if defPath == viper.GetString("config.file") {
			slog.Info("Default config file not found. Using default values")
		} else {
			slog.Warn("Configured config file not found. Using default values")
		}
	}
}

func defaultConfigs() {
	viper.SetDefault("server.port", ":8081")
	viper.SetDefault("ticker.enabled", false)
	viper.SetDefault("ticker.timeout", "50s")
	viper.SetDefault("available_pairs", []string{"BTC/USD", "BTC/CHF", "BTC/EUR"})
}
