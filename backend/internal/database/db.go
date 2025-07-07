package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type DBClient struct {
	DB   *sql.DB
	Pool *pgxpool.Pool
}

func ConnectDB(databaseURL string, logger *slog.Logger) (*DBClient, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute

	var pool *pgxpool.Pool
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			if pingErr := pool.Ping(context.Background()); pingErr == nil {
				logger.Info("Successfully connected to database with pgx!")

				db := stdlib.OpenDBFromPool(pool)

				return &DBClient{DB: db, Pool: pool}, nil
			} else {
				err = pingErr
			}
		}

		logger.Warn("Failed to connect to database, retrying...", "error", err, "attempt", i+1)
		pool.Close()
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, err)
}

func (c *DBClient) Close() error {
	if c.Pool != nil {
		c.Pool.Close()
	}
	return nil
}

func (c *DBClient) Ping() error {
	return c.Pool.Ping(context.Background())
}
