package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func OpenWith(host, port, user, password, database string) (*DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			database),
	)
	return &DB{db}, err
}

func Open(opt *types.PostgresOpt) (*DB, error) {
	return OpenWith(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
}
