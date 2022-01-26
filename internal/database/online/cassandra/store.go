package cassandra

import (
	"context"
	"time"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"

	"github.com/gocql/gocql"
	"github.com/oom-ai/oomstore/pkg/errdefs"

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
		return nil, errdefs.WithStack(err)
	}
	return &DB{Session: session}, nil
}

func (db *DB) CreateTable(ctx context.Context, opt online.CreateTableOpt) error {
	dbOpt := dbutil.DBOpt{
		Backend:     Backend,
		CassandraDB: db.Session,
	}
	return sqlutil.CreateTable(ctx, dbOpt, opt)
}
