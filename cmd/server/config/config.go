package config

import (
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"log"
	"os"
	"strconv"
	"time"
)

type OptionConfig func(*Config) error

type Config struct {
	DBURL string //Путь к базе данных

	JWTSecret   string        //Секрет для JWT
	JWTDuration time.Duration //Время жизни JWT

	GrpcAddr string //Адрес GRPC сервера

	RedisAddr       string        //Путь к редису, example "localhost:6379"
	RedisPassword   string        //Пароль от редиса, example ""
	RedisDB         int           //Номер базы данных от редиса, example 0
	RedisExpiration time.Duration //Время жизни кеша в редисе, если 0 то бессрочно

}

func LoadCfg(opts ...OptionConfig) (*Config, error) {

	dbURL := os.Getenv("KEEPER_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("KEEPER_DB_URL environment variable is required")
	}

	grpcAddr := os.Getenv("KEEPER_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, fmt.Errorf("KEEPER_GRPC_ADDR environment variable is required")
	}

	cfg := &Config{
		DBURL:    dbURL,
		GrpcAddr: grpcAddr,
	}

	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func WithJWT() OptionConfig {

	return func(c *Config) error {

		jwtSecret := os.Getenv("KEEPER_JWT_SECRET")
		if jwtSecret == "" {
			return fmt.Errorf("KEEPER_JWT_SECRET environment variable is required")
		}

		c.JWTSecret = jwtSecret

		jwtDuration := 3600
		jwtDurationStr := os.Getenv("KEEPER_JWT_DURATION_SEC")
		if jwtDurationStr == "" {
			return fmt.Errorf("KEEPER_JWT_SECRET environment variable is required")
		}

		c.JWTDuration = time.Duration(jwtDuration) * time.Second

		return nil

	}
}

func WithRedis() OptionConfig {

	return func(c *Config) error {

		addr := os.Getenv("KEEPER_REDIS_ADDR")
		if addr == "" {
			return fmt.Errorf("KEEPER_REDIS_ADDR environment variable is required")
		}

		c.RedisAddr = addr

		if value := os.Getenv("KEEPER_REDIS_PASSWORD"); value != "" {
			c.RedisPassword = value
		}

		if value := os.Getenv("KEEPER_REDIS_DB"); value != "" {
			valTyped, err := strconv.Atoi(value)
			if err == nil {
				c.RedisDB = valTyped
			} else {
				log.Println(errorsutils.WrapError(err))
			}
		}

		if value := os.Getenv("KEEPER_REDIS_EXPIRATION_SEC"); value != "" {
			valTyped, err := strconv.Atoi(value)
			if err == nil {
				c.RedisExpiration = time.Duration(valTyped) * time.Second
			} else {
				log.Println(errorsutils.WrapError(err))
			}
		}

		return nil

	}
}
