package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	schema, err := migrationFS.ReadFile("migrations/001_create_tables.sql")
	if err != nil {
		return fmt.Errorf("unable to read migration file: %w", err)
	}

	_, err = pool.Exec(context.Background(), string(schema))
	if err != nil {
		return fmt.Errorf("unable to run migrations: %w", err)
	}

	return nil
}
