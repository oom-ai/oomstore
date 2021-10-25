package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

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

func Open(option *types.PostgresDbOpt) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.Database)
}

func OpenWith(host string, port, user, pass, dbName string) (*DB, error) {
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

type WalkFunc = func(slice []interface{}) error

func (db *DB) WalkTable(ctx context.Context, table string, fields []string, limit *uint64, walkFunc WalkFunc) error {
	query := fmt.Sprintf("select %s from %s", strings.Join(fields, ","), table)
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return walkRows(rows, walkFunc)
}

func walkRows(rows *sqlx.Rows, walkFunc WalkFunc) error {
	for rows.Next() {
		slice, err := rows.SliceScan()
		if err != nil {
			return err
		}
		if err := walkFunc(slice); err != nil {
			return err
		}
	}
	return nil
}
