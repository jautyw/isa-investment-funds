package storage_test

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

const (
	dbname   = ""
	user     = "user"
	password = "pass"
)

func TestAll(t *testing.T) {
	ctx := context.Background()

	// 1. Start the postgres ctr and run any migrations on it
	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	// Run any migrations on the database
	_, _, err = ctr.Exec(ctx, []string{"psql", "-U", user, "-d", dbname, "-c", "CREATE TABLE users (id SERIAL, name TEXT NOT NULL, age INT NOT NULL)"})
	require.NoError(t, err)

	//// 2. Create a snapshot of the database to restore later
	//// tt.options comes the test case, it can be specified as e.g. `postgres.WithSnapshotName("custom-snapshot")` or omitted, to use default name
	//err = ctr.Snapshot(ctx, tt.options...)
	//require.NoError(t, err)

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		conn, err := pgx.Connect(context.Background(), dbURL)
		require.NoError(t, err)
		defer conn.Close(context.Background())

		_, err = conn.Exec(ctx, "INSERT INTO users(name, age) VALUES ($1, $2)", "test", 42)
		require.NoError(t, err)

		var name string
		var age int64
		err = conn.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)
		require.NoError(t, err)

		require.Equal(t, "test", name)
		require.EqualValues(t, 42, age)
	})

	// 4. Run as many tests as you need, they will each get a clean database
	t.Run("Test querying empty DB", func(t *testing.T) {
		t.Cleanup(func() {
			err = ctr.Restore(ctx)
			require.NoError(t, err)
		})

		conn, err := pgx.Connect(context.Background(), dbURL)
		require.NoError(t, err)
		defer conn.Close(context.Background())

		var name string
		var age int64
		err = conn.QueryRow(context.Background(), "SELECT name, age FROM users LIMIT 1").Scan(&name, &age)
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}
