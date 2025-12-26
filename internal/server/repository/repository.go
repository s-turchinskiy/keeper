package repository

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/server/models"
)

type DBer interface {
	Close(ctx context.Context)
}
type UserRepositorier interface {
	Create(ctx context.Context, login, passwordHash string) (*models.User, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type SecretRepositorier interface {
	CreateUpdate(ctx context.Context, secret *models.Secret) error
	GetByID(ctx context.Context, userID, secretID string) (*models.Secret, error)
	Delete(ctx context.Context, userID, secretID string) error
	GetAll(ctx context.Context, userID string) ([]*models.Secret, error)
	GetAllWithStatuses(ctx context.Context, userID string) ([]*models.Secret, error)
}
