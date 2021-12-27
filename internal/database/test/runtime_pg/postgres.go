package runtime_pg

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	if out, err := exec.Command("oomplay", "init", "postgres").CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func CreateDatabase(database string) {
	db := rootDB()
	defer db.Close()

	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
		panic(err)
	}
}

func DestroyStore(database string) func() {
	return func() {
		db := rootDB()
		defer db.Close()

		if _, err := db.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", database)); err != nil {
			panic(err)
		}
	}
}

func GetOpt(database string) *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     "127.0.0.1",
		Port:     "25432",
		User:     "oomplay",
		Password: "oomplay",
		Database: database,
	}
}

func RootOpt(database string) *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     "127.0.0.1",
		Port:     "25432",
		User:     "postgres",
		Password: "postgres",
		Database: database,
	}
}

func rootDB() *sqlx.DB {
	opt := RootOpt("postgres")
	if db, err := dbutil.OpenPostgresDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		opt.Database,
	); err != nil {
		panic(err)
	} else {
		return db
	}
}
