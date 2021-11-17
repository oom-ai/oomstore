package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/postgres"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareStore() (context.Context, online.Store) {
	ctx := context.Background()
	opt := runtime_pg.PostgresDbOpt
	store, err := postgres.Open(&types.PostgresOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Password: opt.Password,
		Database: "test",
	})
	if err != nil {
		panic(err)
	}

	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s; ", opt.Database)
	if _, err := store.ExecContext(context.Background(), sql); err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("CREATE DATABASE %s", opt.Database)
	if _, err = store.ExecContext(context.Background(), sql); err != nil {
		panic(err)
	}

	store, err = postgres.Open(&opt)
	if err != nil {
		panic(err)
	}
	return ctx, store
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}
