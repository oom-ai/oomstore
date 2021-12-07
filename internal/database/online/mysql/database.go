package mysql

import (
	"context"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const BackendType = types.MYSQL

var _ online.Store = &DB{}

type DB struct {
	*sqlx.DB
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.Ping()
}

func OpenWith(host, port, user, password, database string) (*DB, error) {
	db, err := sqlx.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
			user,
			password,
			host,
			port,
			database))
	return &DB{db}, err
}

func Open(opt *types.MySQLOpt) (*DB, error) {
	return OpenWith(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
}
