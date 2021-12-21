package runtime_pg

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/orlangure/gnomock"
	mockpg "github.com/orlangure/gnomock/preset/postgres"
)

var PostgresDbOpt types.PostgresOpt

func init() {
	postgresContainer, err := gnomock.Start(
		mockpg.Preset(
			mockpg.WithUser("test", "test"),
			mockpg.WithDatabase("test"),
			mockpg.WithVersion("14.0"),
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

func PrepareDB(t *testing.T) (context.Context, *sqlx.DB) {
	ctx := context.Background()
	db, err := dbutil.OpenPostgresDB(
		PostgresDbOpt.Host,
		PostgresDbOpt.Port,
		PostgresDbOpt.User,
		PostgresDbOpt.Password,
		"test",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(ctx, "drop database if exists oomstore")
	if err != nil {
		t.Fatal(err)
	}

	return ctx, db
}
