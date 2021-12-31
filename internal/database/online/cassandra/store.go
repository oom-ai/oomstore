package cassandra

import (
	"context"
	"time"

	"github.com/gocql/gocql"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Backend   = types.BackendCassandra
	BatchSize = 1000
)

var _ online.Store = &DB{}

type DB struct {
	*gocql.Session
}

func (db *DB) Ping(ctx context.Context) error {
	return nil
}

func (db *DB) Close() error {
	db.Session.Close()
	return nil
}

func Open(option *types.CassandraOpt) (*DB, error) {
	cluster := gocql.NewCluster(option.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: option.User,
		Password: option.Password,
	}
	cluster.Keyspace = option.KeySpace
	if option.Timeout != 0 {
		cluster.Timeout = option.Timeout
	} else {
		cluster.Timeout = time.Second * 5
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &DB{Session: session}, nil
}
