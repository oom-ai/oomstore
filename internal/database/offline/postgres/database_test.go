package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/internal/database/test"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func initDB(t *testing.T) {
	opt := test.PostgresDbopt
	store, err := postgres.Open(&types.PostgresOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Password: opt.Password,
		Database: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if _, err := store.ExecContext(context.Background(), fmt.Sprintf("drop database if exists %s", opt.Database)); err != nil {
		t.Fatal(err)
	}

	if _, err = store.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", opt.Database)); err != nil {
		t.Fatal(err)
	}
}

func initAndOpenDB(t *testing.T) *postgres.DB {
	initDB(t)

	db, err := postgres.Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
