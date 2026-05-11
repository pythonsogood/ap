package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
)

type MessageBrokerType string

const (
	MessageBrokerTypeNATS MessageBrokerType = "nats"
)

type messageBrokerNATSConfig struct {
	ConnectionUrl string `toml:"connection_url" env:"MESSAGE_BROKER_NATS_CONNECTION_URL"`
}

type messageBrokerConfig struct {
	Type           MessageBrokerType       `toml:"type" env:"MESSAGE_BROKER_TYPE"`
	LoggedSubjects []string                `toml:"logged_subjects" env:"MESSAGE_BROKER_LOGGED_SUBJECTS"`
	Nats           messageBrokerNATSConfig `toml:"nats"`
}

type logConfig struct {
	File string `toml:"file" env:"LOG_FILE"`
}

type jobQueueConfig struct {
	GatewayUrl     string `toml:"gateway_url" env:"GATEWAY_URL"`
	WorkerPoolSize uint   `toml:"worker_pool_size" env:"WORKER_POOL_SIZE"`
}

type Config struct {
	RedisUrl      string              `toml:"url" env:"REDIS_URL"`
	MessageBroker messageBrokerConfig `toml:"message-broker"`
	Log           logConfig           `toml:"log"`
	JobQueue      jobQueueConfig      `toml:"job_queue"`
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
		config_path = "./configs/notification/config.toml"
	}

	return ParseConfig(config_path)
}
