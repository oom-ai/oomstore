package runtime_pg

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func PrepareDB(t *testing.T, database string) (context.Context, *sqlx.DB) {
	db, err := prepareDB(database)
	if err != nil {
		t.Fatal(err)
	}
	return context.Background(), db
}

func prepareDB(database string) (*sqlx.DB, error) {
	opt := GetOpt(database)
	return dbutil.OpenPostgresDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		// Postgres creates a database with the same name of the user.
		// We need to connect using this database to drop other databases.
		opt.User,
	)
}

func DestroyStore(database string) func() {
	return func() {
		db, err := prepareDB(database)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if _, err := db.ExecContext(context.Background(),
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", database)); err != nil {
			panic(err)
		}
	}
}

func init() {
	// "dummy" db will not actually be used during testing
	opt := GetOpt("dummy")
	if out, err := exec.Command(
		"oomplay", "init", "postgres",
		"--port", opt.Port,
		"--user", opt.User,
		"--password", opt.Password,
		"--database", opt.Database,
	).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func GetOpt(database string) *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "test",
		Password: "test",
		Database: database,
	}
}
