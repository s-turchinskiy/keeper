package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBURL         string
	Login         string
	Password      string
	ServerAddress string
}

func LoadCfg() (*Config, error) {
	dbURL := os.Getenv("KEEPER_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("KEEPER_DB_URL environment variable is required")
	}

	login := os.Getenv("KEEPER_LOGIN")
	if login == "" {
		return nil, fmt.Errorf("KEEPER_LOGIN environment variable is required")
	}

	password := os.Getenv("KEEPER_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("KEEPER_PASSWORD environment variable is required")
	}

	serverAddres := os.Getenv("KEEPER_SERVER_GRPC_ADDR")
	if serverAddres == "" {
		return nil, fmt.Errorf("KEEPER_SERVER_GRPC_ADDR environment variable is required")
	}

	return &Config{
		DBURL:         dbURL,
		Login:         login,
		Password:      password,
		ServerAddress: serverAddres,
	}, nil
}
