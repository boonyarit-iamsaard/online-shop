package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
}

func Load() (Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.BindEnv("port"); err != nil {
		return Config{}, fmt.Errorf("bind env port: %w", err)
	}
	if err := v.BindEnv("database_url"); err != nil {
		return Config{}, fmt.Errorf("bind env database_url: %w", err)
	}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.Port == "" {
		return Config{}, errors.New("config: PORT is required")
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("config: DATABASE_URL is required")
	}

	return cfg, nil
}
