package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"time"

	"github.com/s-turchinskiy/keeper/internal/server/models"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
)

type SecretRepository struct {
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func NewSecretRepository(postgreDB *PostgreDB) *SecretRepository {
	return &SecretRepository{
		db:   postgreDB.db,
		pool: postgreDB.pool,
	}
}

func (r *SecretRepository) SetSecret(ctx context.Context, secret *models.Secret) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	}(tx)

	query := `
		INSERT INTO keeper.secrets (name, user_id, data, hash)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name, user_id) DO UPDATE SET
			data = $3,
			hash = $4
	`

	_, err = tx.ExecContext(ctx, query,
		secret.ID,
		secret.UserID,
		secret.Data,
		secret.Hash,
	)

	if err != nil {
		return errorsutils.WrapError(err)
	}

	queryStatuses := `
		INSERT INTO keeper.secrets_statuses (name, user_id, last_modified, status)
		VALUES ($1, $2, $3, 'ACTIVE')
		ON CONFLICT (name, user_id) DO UPDATE SET
			last_modified = $3,
			status = 'ACTIVE'
	`

	_, err = tx.ExecContext(ctx, queryStatuses,
		secret.ID,
		secret.UserID,
		secret.LastModified,
	)

	if err != nil {
		return errorsutils.WrapError(err)
	}

	err = tx.Commit()
	if err != nil {
		return errorsutils.WrapError(err)
	}

	return err
}

func (r *SecretRepository) GetSecret(ctx context.Context, userID, secretID string) (*models.Secret, error) {
	query := `
		SELECT s.name, s.user_id, s.data, s.hash, st.last_modified 
		FROM keeper.secrets s
		INNER JOIN keeper.secrets_statuses st
                   ON s.user_id = st.user_id AND s.name = st.name
		WHERE s.user_id = $1 AND s.name = $2
	`

	var secret models.Secret
	err := r.db.QueryRowContext(ctx, query, userID, secretID).Scan(
		&secret.ID,
		&secret.UserID,
		&secret.Data,
		&secret.Hash,
		&secret.LastModified,
	)

	if err != nil {
		return nil, ErrSecretNotFound
	}

	return &secret, nil
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, userID, secretID string) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	}(tx)

	query := `
		DELETE FROM keeper.secrets
		WHERE user_id = $1 AND name = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, secretID)
	if err != nil {
		return err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return ErrSecretNotFound
	}

	queryStatuses := `
		INSERT INTO keeper.secrets_statuses (name, user_id, last_modified, status)
		VALUES ($1, $2, $3, 'DELETED')
		ON CONFLICT (name, user_id) DO UPDATE SET
			last_modified = $3,
			status = 'DELETED'
	`

	_, err = tx.ExecContext(ctx, queryStatuses,
		secretID,
		userID,
		time.Now(),
	)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (r *SecretRepository) ListSecrets(ctx context.Context, userID string) ([]*models.Secret, error) {

	query := `
		SELECT s.name, s.user_id, s.hash, st.last_modified
		FROM keeper.secrets s
		INNER JOIN keeper.secrets_statuses st
                   ON s.user_id = st.user_id AND s.name = st.name
		WHERE s.user_id = $1
		ORDER BY st.last_modified DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var secrets []*models.Secret
	for rows.Next() {
		var secret models.Secret
		err := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.Hash,
			&secret.LastModified,
		)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, &secret)
	}
	return secrets, nil
}

func (r *SecretRepository) ListSecretsWithStatuses(ctx context.Context, userID string) ([]*models.Secret, error) {

	query := `SELECT st.name, st.user_id, st.last_modified, st.status = 'DELETED', s.hash
		FROM keeper.secrets_statuses st
		LEFT JOIN keeper.secrets s
                   ON s.user_id = st.user_id AND s.name = st.name
		WHERE st.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var secrets []*models.Secret
	for rows.Next() {
		var secret models.Secret
		err := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.Hash,
			&secret.Deleted,
			&secret.LastModified,
		)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, &secret)
	}
	return secrets, nil
}
