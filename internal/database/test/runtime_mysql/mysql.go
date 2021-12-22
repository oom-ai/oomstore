package runtime_mysql

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
	return dbutil.OpenMysqlDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		"",
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
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", database),
		); err != nil {
			panic(err)
		}
	}
}

func Reset(database string) {
	opt := GetOpt(database)
	if out, err := exec.Command(
		"oomplay", "init", "mysql",
		"--port", opt.Port,
		"--user", opt.User,
		"--password", opt.Password,
		"--database", opt.Database,
	).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func GetOpt(database string) *types.MySQLOpt {
	return &types.MySQLOpt{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "mysql",
		Database: database,
	}
}
