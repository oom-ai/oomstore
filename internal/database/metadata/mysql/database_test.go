package mysql_test

import (
	"context"
	"sort"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mysql"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, metadata.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *mysql.DB) {
	ctx := context.Background()
	opt := runtime_mysql.MySQLDbOpt
	db, err := mysql.OpenDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		opt.Database,
	)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "drop database if exists test")
	require.NoError(t, err)
	db.Close()

	err = mysql.CreateDatabase(ctx, runtime_mysql.MySQLDbOpt)
	require.NoError(t, err)

	mysqlDB, err := mysql.Open(ctx, &runtime_mysql.MySQLDbOpt)
	require.NoError(t, err)

	return ctx, mysqlDB
}

func TestCreateDatabase(t *testing.T) {
	ctx, store := prepareDB(t)
	defer store.Close()

	var tables []string
	err := store.SelectContext(ctx, &tables,
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'test'
			ORDER BY table_name;`)
	require.NoError(t, err)

	var wantTables []string
	for table := range mysql.META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range mysql.META_VIEW_SCHEMAS {
		wantTables = append(wantTables, table)
	}

	sort.Slice(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	sort.Slice(wantTables, func(i, j int) bool {
		return wantTables[i] < wantTables[j]
	})
	assert.Equal(t, wantTables, tables)
}

//func TestPing(t *testing.T) {
//	test_impl.TestPing(t, prepareStore)
//}
