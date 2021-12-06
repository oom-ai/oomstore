package runtime_pg

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var PostgresDbOpt types.PostgresOpt

func init() {
	postgresContainer, err := gnomock.Start(
		postgres.Preset(
			postgres.WithUser("test", "test"),
			postgres.WithDatabase("test"),
			postgres.WithVersion("14.0"),
		),
		gnomock.WithUseLocalImagesFirst(),
	)
	if err != nil {
		panic(err)
	}

	PostgresDbOpt = types.PostgresOpt{
		Host:     postgresContainer.Host,
		Port:     strconv.Itoa(postgresContainer.DefaultPort()),
		User:     "test",
		Password: "test",
		Database: "oomstore",
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		_ = gnomock.Stop(postgresContainer)
	}()
}
