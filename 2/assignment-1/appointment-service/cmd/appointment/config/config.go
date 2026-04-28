package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
)

type DatabaseType string

const (
	DatabaseTypeSQLite DatabaseType = "sqlite3"
)

type doctorServiceConfig struct {
	Address string `toml:"address" env:"SERVICES_DOCTOR_ADDRESS"`
	Timeout uint   `toml:"timeout" env:"SERVICES_DOCTOR_TIMEOUT"`
}

type servicesConfig struct {
	Doctor doctorServiceConfig `toml:"doctor"`
}

type serverConfig struct {
	Port uint16 `toml:"port" env:"SERVER_PORT"`
}

type databaseSQLite3Config struct {
	Source string `toml:"source" env:"DB_SQLITE_SOURCE"`
}

type databaseConfig struct {
	Type    DatabaseType          `toml:"type" env:"DB_TYPE"`
	Sqlite3 databaseSQLite3Config `toml:"sqlite3"`
}

type Config struct {
	Services servicesConfig `toml:"services"`
	Server   serverConfig   `toml:"server"`
	Database databaseConfig `toml:"database"`
}

func ParseConfig(config_path string) (*Config, error) {
	var conf Config

	_, err := toml.DecodeFile(config_path, &conf)

	if err != nil {
		return nil, err
	}

	err = env.Parse(&conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func NewDefaultConfig() (*Config, error) {
	config_path := os.Getenv("CONFIG_FILE")

	if len(config_path) == 0 {
		config_path = "./configs/appointment/config.toml"
	}

	return ParseConfig(config_path)
}
