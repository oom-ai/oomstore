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
	ctx := context.Background()
	opt := GetOpt(database)
	db, err := dbutil.OpenMysqlDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	return ctx, db
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
