package cassandra_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/cassandra"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_cassandra"
)

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx, session := runtime_cassandra.PrepareDB(t)

	createKeySpace := fmt.Sprintf(`CREATE KEYSPACE %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : 1
	}`, runtime_cassandra.CassandraDbOpt.KeySpace)

	if err := session.Query(createKeySpace).WithContext(ctx).Exec(); err != nil {
		t.Fatal(err)
	}
	session.Close()

	store, err := cassandra.Open(&runtime_cassandra.CassandraDbOpt)
	if err != nil {
		t.Fatal(err)
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
