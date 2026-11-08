package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

)

const (
	envDatabaseURL   = "DATABASE_URL"
	envSupabaseDBURL = "SUPABASE_DB_URL"
)

func databaseURL() string {
	url := os.Getenv(envDatabaseURL)
	if url == "" {
		url = os.Getenv(envSupabaseDBURL)
	}
	return url
}

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := databaseURL()
	if url == "" {
		return nil, fmt.Errorf("%s or %s environment variable is required", envDatabaseURL, envSupabaseDBURL)
	}

	config, err := pgxpool.ParseConfig(url)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database config: %w", err)
    }

    config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create pgx pool: %w", err)
    }

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
