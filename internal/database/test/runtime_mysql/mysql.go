package runtime_mysql

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var MySQLDbOpt = types.MySQLOpt{
	Host:     "127.0.0.1",
	Port:     "3306",
	User:     "root",
	Password: "mysql",
	Database: "oomstore_test",
}

func PrepareDB() (context.Context, *sqlx.DB) {
	ctx := context.Background()
	db, err := dbutil.OpenMysqlDB(
		MySQLDbOpt.Host,
		MySQLDbOpt.Port,
		MySQLDbOpt.User,
		MySQLDbOpt.Password,
		"",
	)
	if err != nil {
		panic(err)
	}

	if out, err := exec.Command(
		"oomplay", "init", "mysql",
		"--port", MySQLDbOpt.Port,
		"--user", MySQLDbOpt.User,
		"--password", MySQLDbOpt.Password,
		"--database", MySQLDbOpt.Database,
	).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
	return ctx, db
}
