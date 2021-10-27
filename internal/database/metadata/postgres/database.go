package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/onestore-ai/onestore/internal/database/metadata"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

var _ metadata.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func Open(option *types.PostgresDbOpt) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.Database)
}

func OpenWith(host, port, user, pass, dbName string) (*DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			pass,
			host,
			port,
			dbName),
	)
	return &DB{db}, err
}
