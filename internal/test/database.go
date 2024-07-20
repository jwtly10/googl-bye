package test

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabaseConfiguration struct {
	RootRelativePath string
}

func NewTestDatabaseWithContainer(config TestDatabaseConfiguration) (*postgres.PostgresContainer, *sql.DB, error) {
	postgresContainer, err := postgres.Run(context.Background(), "postgres:latest",
		testcontainers.WithEnv(map[string]string{
			"DATABASE_USER":     "postgres",
			"DATABASE_NAME":     "googl-bye-db",
			"DATABASE_PORT":     "5432",
			"DATABASE_SSL_MODE": "disable",
		}),
		postgres.WithInitScripts(filepath.Join(config.RootRelativePath, "db/init.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	conn, _ := postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	db, err := sql.Open("postgres", conn)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, db, nil
}
