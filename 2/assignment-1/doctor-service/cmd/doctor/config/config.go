package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
)

type CacheType string

const (
	CacheTypeRedis CacheType = "redis"
)

type DatabaseType string

const (
	DatabaseTypeSQLite   DatabaseType = "sqlite3"
	DatabaseTypePostgres DatabaseType = "postgres"
)

type MessageBrokerType string

const (
	MessageBrokerTypeNATS MessageBrokerType = "nats"
)

type serverConfig struct {
	Port uint16 `toml:"port" env:"SERVER_PORT"`
}

type cacheRedisConfig struct {
	Url string `toml:"url" env:"CACHE_REDIS_URL"`
}

type cacheConfig struct {
	Type  CacheType        `toml:"type" env:"CACHE_TYPE"`
	Redis cacheRedisConfig `toml:"redis"`
	Ttl   uint             `toml:"ttl" env:"CACHE_TTL"`
}

type databaseSQLite3Config struct {
	Source string `toml:"source" env:"DB_SQLITE_SOURCE"`
}

type databasePostgresConfig struct {
	Host     string `toml:"host" env:"DB_POSTGRES_HOST"`
	Port     uint   `toml:"port" env:"DB_POSTGRES_PORT"`
	User     string `toml:"user" env:"DB_POSTGRES_USER"`
	Password string `toml:"password" env:"DB_POSTGRES_PASSWORD"`
	Db       string `toml:"db" env:"DB_POSTGRES_DB"`
}

type databaseConfig struct {
	Type     DatabaseType           `toml:"type" env:"DB_TYPE"`
	Sqlite3  databaseSQLite3Config  `toml:"sqlite3"`
	Postgres databasePostgresConfig `toml:"postgres"`
}

type messageBrokerNATSConfig struct {
	ConnectionUrl string `toml:"connection_url" env:"MESSAGE_BROKER_NATS_CONNECTION_URL"`
}

type messageBrokerConfig struct {
	Type MessageBrokerType       `toml:"type" env:"MESSAGE_BROKER_TYPE"`
	Nats messageBrokerNATSConfig `toml:"nats"`
}

type Config struct {
	Server        serverConfig        `toml:"server"`
	Cache         cacheConfig         `toml:"cache"`
	Database      databaseConfig      `toml:"database"`
	MessageBroker messageBrokerConfig `toml:"message-broker"`
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
		config_path = "./configs/doctor/config.toml"
	}

	return ParseConfig(config_path)
}
