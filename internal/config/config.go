package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

// Http represents the HTTP server configuration.
type Http struct {
	Address      string        `mapstructure:"address"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

// Postgres represents the database configuration.
type Postgres struct {
	Name         string        `yaml:"dbname"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	SSLMode      string        `yaml:"sslmode"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
}

// Log
type Log struct {
	LogPath string `mapstructure:"logPath"`
}

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
}

// LoadConfig loads the configuration from the config file config.yaml.
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}
