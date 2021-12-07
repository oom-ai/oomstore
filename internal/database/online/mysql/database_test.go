package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/internal/database/online/mysql"
	"github.com/ethhte88/oomstore/internal/database/online/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_mysql"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func prepareStore() (context.Context, online.Store) {
	ctx := context.Background()
	opt := runtime_mysql.MySQLDbOpt
	store, err := mysql.Open(&types.MySQLOpt{
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

	store, err = mysql.Open(&opt)
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
