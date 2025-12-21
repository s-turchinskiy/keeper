package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"

	"github.com/s-turchinskiy/keeper/internal/server/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository struct {
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func NewUserRepository(postgreDB *PostgreDB) *UserRepository {
	return &UserRepository{
		db:   postgreDB.db,
		pool: postgreDB.pool,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, login, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO keeper.users (login, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, login, password_hash, created_at
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, created_at
		FROM keeper.users
		WHERE login = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, created_at
		FROM keeper.users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}
