package snowflake_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/snowflake"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB()
	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", DATABASE)); err != nil {
		t.Fatal(t)
	}
	return ctx, db
}

func prepareDB() (context.Context, *snowflake.DB) {
	opt := types.SnowflakeOpt{
		Account:  os.Getenv("SNOWFLAKE_TEST_ACCOUNT"),
		User:     os.Getenv("SNOWFLAKE_TEST_USER"),
		Password: os.Getenv("SNOWFLAKE_TEST_PASSWORD"),
	}
	db, err := snowflake.Open(&opt)
	if err != nil {
		panic(err)
	}
	return context.Background(), db
}

func destroyStore(database string) func() {
	return func() {
		ctx, db := prepareDB()
		defer db.Close()

		if _, err := db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", database)); err != nil {
			panic(err)
		}
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

func TestSnapshot(t *testing.T) {
	test_impl.TestSnapshot(t, prepareStore, destroyStore(DATABASE))
}

// We don't use test_impl.TestTableSchema because snowflake cannot be
// accessed by two different sessions.
func TestTableSchema(t *testing.T) {
	t.Cleanup(destroyStore(DATABASE))

	ctx, store := prepareStore(t)
	defer store.Close()
	db := store.(*snowflake.DB)

	if _, err := db.ExecContext(ctx, `create table "offline_batch_1_1"("user" varchar(16), "age" smallint)`); err != nil {
		t.Fatal(err)
	}

	actual, err := store.TableSchema(ctx, "offline_batch_1_1")
	require.NoError(t, err)
	require.Equal(t, 2, len(actual.Fields))

	expected := types.DataTableSchema{
		Fields: []types.DataTableFieldSchema{
			{
				Name:      "user",
				ValueType: types.String,
			},
			{
				Name:      "age",
				ValueType: types.Int64,
			},
		},
	}
	require.ElementsMatch(t, expected.Fields, actual.Fields)
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, destroyStore(DATABASE))
}
