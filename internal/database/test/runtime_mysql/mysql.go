package runtime_mysql

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mysql"
)

var MySQLDbOpt types.MySQLOpt

func init() {
	mysqlContainer, err := gnomock.Start(
		mysql.Preset(
			mysql.WithUser("test", "test"),
			mysql.WithDatabase("test"),
			mysql.WithVersion("8.0"),
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
