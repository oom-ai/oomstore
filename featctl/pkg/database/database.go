package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Option struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DbName string
}

func Open(option *Option) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.DbName)
}

func OpenWith(host, port, user, pass, dbName string) (*DB, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			user,
			pass,
			host,
			port,
			dbName),
	)
	return &DB{db}, err
}

func (db *DB) TableExists(ctx context.Context, table string) (bool, error) {
	var result string
	err := db.QueryRowContext(ctx, `show tables like ?`, table).Scan(&result)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
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
	err := db.QueryRowContext(ctx, fmt.Sprintf("show columns from `%s` like ?", table), column).
		Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf("column '%s' not found in table '%s'", column, table)
	}
	return result, err
}
