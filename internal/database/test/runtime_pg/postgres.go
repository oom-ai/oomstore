package runtime_pg

import (
	"context"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var PostgresDbOpt = types.PostgresOpt{
	Host:     "127.0.0.1",
	Port:     "5432",
	User:     "test",
	Password: "test",
	Database: "oomstore_test",
}

func PrepareDB(t *testing.T) (context.Context, *sqlx.DB) {
	ctx := context.Background()
	db, err := dbutil.OpenPostgresDB(
		PostgresDbOpt.Host,
		PostgresDbOpt.Port,
		PostgresDbOpt.User,
		PostgresDbOpt.Password,
		// Postgres creates a database with the same name of the user.
		// We need to connect using this database to drop other databases.
		PostgresDbOpt.User,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := exec.Command(
		"oomplay", "init", "postgres",
		"--port", PostgresDbOpt.Port,
		"--user", PostgresDbOpt.User,
		"--password", PostgresDbOpt.Password,
		"--database", PostgresDbOpt.Database,
	).Run(); err != nil {
		t.Fatal(err)
	}

	return ctx, db
}
