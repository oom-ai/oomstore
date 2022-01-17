package redshift_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/redshift"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB(t)
	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", DATABASE)); err != nil {
		t.Fatal(err)
	}
	db.Close()

	store, err := redshift.Open(getOpt(DATABASE))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func prepareDB(t *testing.T) (context.Context, *redshift.DB) {
	// open the default db
	db, err := redshift.Open(getOpt("dev"))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	return ctx, db
}

func destroyStore(database string) func() {
	return func() {
		// open the default db
		db, err := redshift.Open(getOpt("dev"))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		ctx := context.Background()

		if _, err = db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", database)); err != nil {
			panic(err)
		}
	}
}

func getOpt(dbname string) *types.PostgresOpt {
	return &types.PostgresOpt{
		Host:     os.Getenv("REDSHIFT_TEST_HOST"),
		User:     os.Getenv("REDSHIFT_TEST_USER"),
		Password: os.Getenv("REDSHIFT_TEST_PASSWORD"),
		Port:     "5439",
		Database: dbname,
	}
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore(DATABASE))
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore, destroyStore(DATABASE))
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore, destroyStore(DATABASE))
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore, destroyStore(DATABASE))
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, destroyStore(DATABASE), func(ctx context.Context) {
		db, err := redshift.Open(getOpt(DATABASE))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		if _, err = db.ExecContext(ctx, `create table "offline_batch_1_1"("user" varchar(16), "age" smallint)`); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, destroyStore(DATABASE))
}
