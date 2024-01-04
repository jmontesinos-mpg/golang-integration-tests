package testcontainers

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq" // add this

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestContainers is a test to demonstrate a simple scenario using TestContainers to start up a Postgres DB.
func TestContainers(t *testing.T) {
	ctx := context.Background()

	// Starts a Postgres container using the host Docker installation.
	postgresContainer, err := postgres.RunContainer(ctx,
		// Default image is a bit obsolete, so it is always advisable to use the one that matches the deployment.
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		// Initializer script, allows to create a DB and all the required tables.
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sh")),
		// By default the container will start when Docker says it is ready,
		// though it could happen with Postgres that the DB itself is not yet ready to accept connections,
		// so it is required to wait for a condition, a log entry in this case, to say the container is ready.
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)))
	require.NoError(t, err)

	t.Cleanup(func() {
		// Clean up the container from Docker, this is quite important to not leave a rather heavy process running in background.
		if err := postgresContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	desc, err := dbClientCall(connStr)
	require.NoError(t, err)

	require.Equal(t, "this is a test", desc)
}

// dbClientCall simulates the business logic of an existing service calling Postgres to retrieve data from the storage layers.
func dbClientCall(connStr string) (string, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	insertRes, err := db.Exec("INSERT INTO testdb (name) VALUES ($1)", "this is a test")
	if err != nil {
		return "", err
	}
	_, err = insertRes.RowsAffected()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("SELECT name FROM testdb")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("No DB rows returned")
	}
	var desc string
	rows.Scan(&desc)

	return desc, nil
}
