package repository

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/server/models"
)

type DBer interface {
	Close(ctx context.Context)
}
type UserRepositorier interface {
	CreateUser(ctx context.Context, login, passwordHash string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
}

type SecretRepositorier interface {
	SetSecret(ctx context.Context, secret *models.Secret) error
	GetSecret(ctx context.Context, userID, secretID string) (*models.Secret, error)
	DeleteSecret(ctx context.Context, userID, secretID string) error
	ListSecrets(ctx context.Context, userID string) ([]*models.Secret, error)
	ListSecretsWithStatuses(ctx context.Context, userID string) ([]*models.Secret, error)
}
