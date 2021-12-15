package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/postgres"
	"github.com/ethhte88/oomstore/internal/database/online/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_pg"
)

func prepareStore() (context.Context, online.Store) {
	ctx, db := runtime_pg.PrepareDB()

	_, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", runtime_pg.PostgresDbOpt.Database))
	if err != nil {
		panic(err)
	}
	db.Close()

	store, err := postgres.Open(&runtime_pg.PostgresDbOpt)
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

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
