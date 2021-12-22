package postgres_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/postgres"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
	if err := runtime_pg.Reset(DATABASE); err != nil {
		panic(err)
	}
}

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx, db := runtime_pg.PrepareDB(t, DATABASE)
	opt := runtime_pg.GetOpt(DATABASE)

	_, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", opt.Database))
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	store, err := postgres.Open(opt)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func TestOpen(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}

func TestPing(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestPing(t, prepareStore)
}
