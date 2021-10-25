package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func OpenWith(host, port, user, pass, database string) (*DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			pass,
			host,
			port,
			database),
	)
	return &DB{db}, err
}
