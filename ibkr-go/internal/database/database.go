package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/db"
)

const (
	maxConns          = 25
	minConns          = 5
	maxConnLifetime   = 1 * time.Hour
	maxConnIdleTime   = 30 * time.Minute
	healthCheckPeriod = 1 * time.Minute
)

// DB wraps the database connection pool and queries.
type DB struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

// New creates a new database connection pool.
func New(ctx context.Context, writeDSN, readDSN string) (*DB, error) {
	config, err := pgxpool.ParseConfig(writeDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool.
	config.MaxConns = maxConns
	config.MinConns = minConns
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnIdleTime = maxConnIdleTime
	config.HealthCheckPeriod = healthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		Pool:    pool,
		Queries: db.New(pool),
	}, nil
}

// Close closes the database connection pool.
func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

// Health checks the database connection health.
func (d *DB) Health(ctx context.Context) error {
	if d.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return d.Pool.Ping(ctx)
}
