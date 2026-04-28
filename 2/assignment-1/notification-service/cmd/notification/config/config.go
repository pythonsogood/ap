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
	ConnectionUrl string `toml:"connection_url" env:"NATS_CONNECTION_URL"`
}

type messageBrokerConfig struct {
	Type MessageBrokerType       `toml:"type"`
	Nats messageBrokerNATSConfig `toml:"nats"`
}

type Config struct {
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
		config_path = "./configs/notification/config.toml"
	}

	return ParseConfig(config_path)
}
