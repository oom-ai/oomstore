package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/mysql"
	"github.com/ethhte88/oomstore/internal/database/offline/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_mysql"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *mysql.DB) {
	ctx := context.Background()
	opt := runtime_mysql.MySQLDbOpt
	store, err := mysql.Open(&types.MySQLOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
	})
	require.NoError(t, err)
	defer store.Close()

	_, err = store.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", opt.Database))
	require.NoError(t, err)

	_, err = store.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", opt.Database))
	require.NoError(t, err)

	db, err := mysql.Open(&runtime_mysql.MySQLDbOpt)
	require.NoError(t, err)
	return ctx, db
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
