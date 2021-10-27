package test

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/oom-ai/oomstore/pkg/onestore/types"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var PostgresDbopt types.PostgresDbOpt

func init() {
	postgresContainer, err := gnomock.Start(postgres.Preset(
		postgres.WithUser("test", "test"),
		postgres.WithDatabase("test"),
		postgres.WithTimezone("Asia/Shanghai"),
	))
	if err != nil {
		panic(err)
	}

	PostgresDbopt = types.PostgresDbOpt{
		Host:     postgresContainer.Host,
		Port:     strconv.Itoa(postgresContainer.DefaultPort()),
		User:     "test",
		Pass:     "test",
		Database: "onestore",
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		_ = gnomock.Stop(postgresContainer)
	}()
}
