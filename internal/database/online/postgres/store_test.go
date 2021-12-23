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
	test_impl.TestOpen(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, runtime_pg.DestroyStore(DATABASE))
}
