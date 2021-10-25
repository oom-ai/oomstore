package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type RowMap = map[string]interface{}

type DB struct {
	*sqlx.DB
}

type Option struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DbName string
}

func Open(option Option) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.DbName)
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
