package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"log"
	"time"
)

const schemaName = "keeper"

type PostgreDB struct {
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, addr string) (*PostgreDB, error) {

	db, err := sqlx.Open("pgx", addr)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, errorsutils.WrapError(err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	p := &PostgreDB{db: db, pool: pool}

	_, err = p.db.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schemaName))
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{SchemaName: schemaName})
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/server/repository/postgres/migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, errorsutils.WrapError(err)
	}

	return p, nil

}

func (p *PostgreDB) Close(ctx context.Context) {

	p.pool.Close()

	err := p.db.Close()

	if err != nil {
		log.Println("PostgreSQL stopped with error %w", err)
	} else {
		log.Println("PostgreSQL stopped")
	}

}
