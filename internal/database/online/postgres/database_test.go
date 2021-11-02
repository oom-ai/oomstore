package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/test"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareStore() (context.Context, online.Store) {
	ctx := context.Background()
	opt := runtime_pg.PostgresDbopt
	store, err := Open(&types.PostgresOpt{
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

	store, err = Open(&opt)
	if err != nil {
		panic(err)
	}
	return ctx, store
}

func TestOpen(t *testing.T) {
	test.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	test.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	test.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}
