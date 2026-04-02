package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type DatabaseType string

const (
	DatabaseTypeSQLite DatabaseType = "sqlite3"
)

type serverConfig struct {
	Port uint16 `toml:"port"`
}

type databaseSQLite3Config struct {
	Source string `toml:"source"`
}

type databaseConfig struct {
	Type    DatabaseType          `toml:"type"`
	Sqlite3 databaseSQLite3Config `toml:"sqlite3"`
}

type Config struct {
	Server   serverConfig   `toml:"server"`
	Database databaseConfig `toml:"database"`
}

func ParseConfig(config_path string) (*Config, error) {
	var conf Config

	_, err := toml.DecodeFile(config_path, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func NewDefaultConfig() (*Config, error) {
	config_path := os.Getenv("CONFIG_FILE")

	if len(config_path) == 0 {
		config_path = "./configs/doctor/config.toml"
	}

	return ParseConfig(config_path)
}
