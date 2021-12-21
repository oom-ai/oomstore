package redshift_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/redshift"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB(t)
	if _, err := db.ExecContext(ctx, "CREATE DATABASE test"); err != nil {
		t.Fatal(err)
	}

	store, err := redshift.Open(getOpt("test"))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

// Check if DB 'test' exists in the redshift cluster
// Redshift does not support CREATE DATABASE IF NOT EXISTS or DROP DATABASE IF EXISTS
func existsTestDB(ctx context.Context, db *redshift.DB) (bool, error) {
	row := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'test')")
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func prepareDB(t *testing.T) (context.Context, *redshift.DB) {
	// open the default db
	db, err := redshift.Open(getOpt("dev"))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	exists, err := existsTestDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		if _, err = db.ExecContext(ctx, "DROP DATABASE test"); err != nil {
			t.Fatal(err)
		}
	}

	return ctx, db
}

func getOpt(dbname string) *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     os.Getenv("REDSHIFT_TEST_HOST"),
		User:     os.Getenv("REDSHIFT_TEST_USER"),
		Password: os.Getenv("REDSHIFT_TEST_PASSWORD"),
		Port:     "5439",
		Database: dbname,
	}
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore)
}
