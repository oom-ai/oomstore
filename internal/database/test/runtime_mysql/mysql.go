package runtime_mysql

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
	"github.com/orlangure/gnomock"
	mockmysql "github.com/orlangure/gnomock/preset/mysql"
)

var MySQLDbOpt types.MySQLOpt

func init() {
	mysqlContainer, err := gnomock.Start(
		mockmysql.Preset(
			mockmysql.WithUser("test", "test"),
			mockmysql.WithDatabase("test"),
			mockmysql.WithVersion("8.0"),
		),
		gnomock.WithUseLocalImagesFirst(),
	)
	if err != nil {
		panic(err)
	}

	MySQLDbOpt = types.MySQLOpt{
		Host:     mysqlContainer.Host,
		Port:     strconv.Itoa(mysqlContainer.DefaultPort()),
		User:     "test",
		Password: "test",
		Database: "test",
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		_ = gnomock.Stop(mysqlContainer)
	}()
}

func PrepareDB() (context.Context, *sqlx.DB) {
	ctx := context.Background()
	db, err := dbutil.OpenMysqlDB(
		MySQLDbOpt.Host,
		MySQLDbOpt.Port,
		MySQLDbOpt.User,
		MySQLDbOpt.Password,
		MySQLDbOpt.Database,
	)
	if err != nil {
		panic(err)
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("drop database if exists %s", MySQLDbOpt.Database))
	if err != nil {
		panic(err)
	}
	return ctx, db
}
