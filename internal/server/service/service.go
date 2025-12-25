package service

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/server/models"
	redisclient "github.com/s-turchinskiy/keeper/internal/server/redis_client"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/s-turchinskiy/keeper/internal/server/token"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Servicer interface {
	GetNewConnectionNumber(ctx context.Context) uint64
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (string, *models.User, error)

	SetSecret(ctx context.Context, secret *models.Secret) error
	GetSecret(ctx context.Context, userID, secretID string) (*models.Secret, error)
	DeleteSecret(ctx context.Context, userID, secretID string) error
	ListSecrets(ctx context.Context, userID string) ([]*models.Secret, error)
	SyncFromClient(ctx context.Context, userID string, clientSecrets []*models.Secret) (updateInClients []*models.Secret, err error)
}

type OptionService func(*Service)

type Service struct {
	TokenManager            token.TokenManager
	usersRepository         repository.UserRepositorier
	secretRepository        repository.SecretRepositorier
	currentConnectionNumber uint64

	redisClient *redisclient.RedisClient
}

func NewService(tokenManager token.TokenManager,
	usersRepository repository.UserRepositorier,
	secretRepository repository.SecretRepositorier,
	opts ...OptionService) *Service {

	s := &Service{
		TokenManager:     tokenManager,
		usersRepository:  usersRepository,
		secretRepository: secretRepository,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s

}

func WithRedis(rdb *redis.Client, expiration time.Duration) OptionService {

	return func(s *Service) {

		err := rdb.Ping(context.Background()).Err()
		if err != nil {
			log.Printf("couldn't connect to redis, error := %v\n", err)
			return
		}

		s.redisClient = redisclient.NewRedisClient(rdb, expiration)

		fmt.Println("connect to redis success")
	}
}
