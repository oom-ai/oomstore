package postgres

import (
	"context"
	"sort"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func initDB(t *testing.T) {
	opt := runtime_pg.PostgresDbOpt
	store, err := Open(&types.PostgresOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Password: opt.Password,
		Database: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.ExecContext(context.Background(), "drop database if exists oomstore"); err != nil {
		t.Fatal(err)
	}
	store.Close()

	if err := CreateDatabase(context.Background(), runtime_pg.PostgresDbOpt); err != nil {
		t.Fatal(err)
	}
}

func initAndOpenDB(t *testing.T) *DB {
	initDB(t)

	db, err := Open(&runtime_pg.PostgresDbOpt)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCreateDatabase(t *testing.T) {
	ctx := context.Background()
	if err := CreateDatabase(ctx, runtime_pg.PostgresDbOpt); err != nil {
		t.Fatal(err)
	}

	store, err := Open(&runtime_pg.PostgresDbOpt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var tables []string
	if err = store.SelectContext(ctx, &tables,
		`SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;`); err != nil {
		t.Fatal(err)
	}

	var wantTables []string
	for table := range META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range META_VIEW_SCHEMAS {
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
