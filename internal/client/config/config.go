package config

import (
	"fmt"
	"os"
)

type OptionConfig func(*Config) error

type Config struct {
	DBURL         string
	Login         string
	Password      string
	ServerAddress string
}

func LoadCfg(opts ...OptionConfig) (*Config, error) {

	login := os.Getenv("KEEPER_LOGIN")
	if login == "" {
		return nil, fmt.Errorf("KEEPER_LOGIN environment variable is required")
	}

	password := os.Getenv("KEEPER_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("KEEPER_PASSWORD environment variable is required")
	}

	serverAddress := os.Getenv("KEEPER_SERVER_GRPC_ADDR")
	if serverAddress == "" {
		return nil, fmt.Errorf("KEEPER_SERVER_GRPC_ADDR environment variable is required")
	}

	cfg := &Config{
		Login:         login,
		Password:      password,
		ServerAddress: serverAddress,
	}

	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func WithDB() OptionConfig {

	return func(c *Config) error {

		dbURL := os.Getenv("KEEPER_DB_URL")
		if dbURL == "" {
			return fmt.Errorf("KEEPER_DB_URL environment variable is required")
		}

		c.DBURL = dbURL

		return nil

	}
}
