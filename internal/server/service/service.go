package service

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/s-turchinskiy/keeper/internal/server/token"
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

type Service struct {
	TokenManager            token.TokenManager
	usersRepository         repository.UserRepositorier
	secretRepository        repository.SecretRepositorier
	currentConnectionNumber uint64
}

func NewService(tokenManager token.TokenManager,
	usersRepository repository.UserRepositorier,
	secretRepository repository.SecretRepositorier) *Service {

	return &Service{
		TokenManager:     tokenManager,
		usersRepository:  usersRepository,
		secretRepository: secretRepository,
	}
}
