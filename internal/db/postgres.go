package db

import (
	"context"
	"fmt"
	"os"

	"github.com/dgrco/TeamActivityTracker-api/internal/environment"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a pool of connections to the database.
func SetupDatabase(env *environment.Environment) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), env.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("successfully connected to database")

	return dbpool
}
