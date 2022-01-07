package cassandra_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/cassandra"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_cassandra"
)

var KEYSPACE string

func init() {
	KEYSPACE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx, session := runtime_cassandra.PrepareDB()

	createKeySpace := fmt.Sprintf(`CREATE KEYSPACE %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : 1
	}`, KEYSPACE)

	if err := session.Query(createKeySpace).WithContext(ctx).Exec(); err != nil {
		t.Fatal(err)
	}
	session.Close()

	store, err := cassandra.Open(runtime_cassandra.GetOpt(KEYSPACE))
	if err != nil {
		t.Fatal(err)
	}
	return ctx, store
}

func destroyStore(keySpace string) func() {
	return func() {
		ctx, session := runtime_cassandra.PrepareDB()
		query := fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", keySpace)
		if err := session.Query(query).WithContext(ctx).Exec(); err != nil {
			panic(err)
		}
	}
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, destroyStore(KEYSPACE))
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, destroyStore(KEYSPACE))
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, destroyStore(KEYSPACE))
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, destroyStore(KEYSPACE))
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, destroyStore(KEYSPACE))
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, destroyStore(KEYSPACE))
}

func TestPrepareStreamTable(t *testing.T) {
	test_impl.TestPrepareStreamTable(t, prepareStore, destroyStore(KEYSPACE))
}

func TestPush(t *testing.T) {
	test_impl.TestPush(t, prepareStore, destroyStore(KEYSPACE))
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore(KEYSPACE))
}
