package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func Open(option Option) (*DB, error) {
	return OpenWith(option.Host, option.Port, option.User, option.Pass, option.DbName)
}

func OpenWith(host, port, user, pass, dbName string) (*DB, error) {
	db, err := sqlx.Open(
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
	err := db.QueryRowContext(ctx, fmt.Sprintf("show columns from `%s` like '%s'", table, column)).
		Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf("column '%s' not found in table '%s'", column, table)
	}
	return result, err
}

func (db *DB) TableExists(ctx context.Context, table string) (bool, error) {
	var result string
	err := db.GetContext(ctx, &result, fmt.Sprintf(`show tables like '%s'`, table))
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

type WalkFunc = func(values []interface{}) error

func (db *DB) WalkTable(ctx context.Context, table string, fields []string, limit *uint64, walkFunc WalkFunc) error {
	marks := []string{}
	for range fields {
		marks = append(marks, "?")
	}
	query := fmt.Sprintf("select %s from %s", strings.Join(marks, ","), table)
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := db.QueryxContext(ctx, query, fields)
	if err != nil {
		return err
	}
	defer rows.Close()

	return walkRows(rows, walkFunc)
}

func walkRows(rows *sqlx.Rows, walkFunc WalkFunc) error {
	for rows.Next() {
		record, err := rows.SliceScan()
		if err != nil {
			return err
		}
		if err := walkFunc(record); err != nil {
			return err
		}
	}
	return nil
}
