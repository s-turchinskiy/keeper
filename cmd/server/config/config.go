package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DBURL       string        //Путь к базе данных
	JWTSecret   string        //Секрет для JWT
	JWTDuration time.Duration //Время жизни JWT
	GrpcAddr    string        //адрес GRPC сервера
}

func LoadCfg() (*Config, error) {

	dbURL := os.Getenv("KEEPER_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("KEEPER_DB_URL environment variable is required")
	}

	grpcAddr := os.Getenv("KEEPER_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, fmt.Errorf("KEEPER_GRPC_ADDR environment variable is required")
	}

	jwtSecret := os.Getenv("KEEPER_JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("KEEPER_JWT_SECRET environment variable is required")
	}

	jwtDuration := 3600
	jwtDurationStr := os.Getenv("KEEPER_JWT_DURATION_SEC")
	if jwtDurationStr == "" {
		return nil, fmt.Errorf("KEEPER_JWT_SECRET environment variable is required")
	}

	cfg := &Config{
		DBURL:       dbURL,
		JWTSecret:   jwtSecret,
		JWTDuration: time.Duration(jwtDuration) * time.Second,
		GrpcAddr:    grpcAddr,
	}
	return cfg, nil
}
