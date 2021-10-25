package database

import (
	"context"
	"database/sql"
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

type Column struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   string
}

func (db *DB) ColumnInfo(ctx context.Context, table string, column string) (Column, error) {
	var result Column
	err := db.QueryRowContext(ctx, fmt.Sprintf(`show columns from "%s" like "%s"`, table, column)).
		Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf(`column "%s" not found in table "%s"`, column, table)
	}
	return result, err
}

func (db *DB) TableExists(ctx context.Context, table string) (bool, error) {
	var result string
	err := db.GetContext(ctx, &result, fmt.Sprintf(`show tables like "%s"`, table))
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}
