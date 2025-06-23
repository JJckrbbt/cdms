// internal/database/db.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // NEW: pgx connection pool
	"github.com/jackc/pgx/v5/stdlib"  // NEW: pgx adapter for database/sql
	// REMOVED: _ "github.com/lib/pq"
)

// DBClient now holds both the standard DB handle and the native pgx pool.
// This allows for maximum flexibility.
type DBClient struct {
	DB   *sql.DB
	Pool *pgxpool.Pool // NEW: Native connection pool for high-performance operations
}

// ConnectDB is updated to establish a connection using pgxpool.
func ConnectDB(databaseURL string, logger *slog.Logger) (*DBClient, error) {
	// NEW: pgx has a more robust config structure for creating a connection pool.
	// This is where you set connection pool settings, replacing db.SetMaxOpenConns, etc.
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Set pool configuration here.
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute

	var pool *pgxpool.Pool
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		// Attempt to create the connection pool.
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			// Now, run a Ping on the pool to ensure connectivity.
			if pingErr := pool.Ping(context.Background()); pingErr == nil {
				logger.Info("Successfully connected to database with pgx!")

				// Create a standard library *sql.DB object from the pgxpool.
				// This makes it compatible with sqlc.
				db := stdlib.OpenDBFromPool(pool)

				return &DBClient{DB: db, Pool: pool}, nil
			} else {
				err = pingErr // Update err to the ping error for logging
			}
		}

		logger.Warn("Failed to connect to database, retrying...", "error", err, "attempt", i+1)
		pool.Close() // Close the failed pool before retrying
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, err)
}

// Close gracefully closes the connection pool.
func (c *DBClient) Close() error {
	if c.Pool != nil {
		c.Pool.Close()
	}
	return nil
}

// Ping checks the database connection.
func (c *DBClient) Ping() error {
	return c.Pool.Ping(context.Background())
}
