package runtime_tidb

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/jmoiron/sqlx"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	if out, err := exec.Command("oomplay", "init", "tidb").CombinedOutput(); err != nil {
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

		if _, err := db.ExecContext(context.Background(),
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", database),
		); err != nil {
			panic(err)
		}
	}
}

func GetOpt(database string) *types.MySQLOpt {
	return &types.MySQLOpt{
		Host:     "127.0.0.1",
		Port:     "24000",
		User:     "oomplay",
		Password: "oomplay",
		Database: database,
	}
}

func rootDB() *sqlx.DB {
	if db, err := dbutil.OpenMysqlDB("127.0.0.1", "24000", "root", "", ""); err != nil {
		panic(err)
	} else {
		return db
	}
}
