package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to our PostgreSQL database.
func Connect() (*pgxpool.Pool, error) {

	// Read the database connection string
	// from the DATABASE_URL environment variable.
	databaseURL := os.Getenv("DATABASE_URL")

	// Make sure the connection string exists.
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Create a connection pool.
	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	// Return the error if the pool could not be created.
	if err != nil {
		return nil, err
	}

	// Test that the database is actually reachable.
	err = pool.Ping(context.Background())

	if err != nil {
		pool.Close()
		return nil, err
	}

	// Return the working connection pool.
	return pool, nil
}
