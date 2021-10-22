package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	// entity
	CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error
	GetEntity(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) ([]*types.Entity, error)
	UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error

	// feature
	CreateFeature(ctx context.Context, opt types.CreateFeatureOpt) error
	GetFeature(ctx context.Context, featureName string) (*types.Feature, error)
	ListFeature(ctx context.Context, groupName *string) ([]*types.Feature, error)
	UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error

	// rich feature
	GetRichFeature(ctx context.Context, featureName string) (*types.RichFeature, error)
	GetRichFeatures(ctx context.Context, featureNames []string) ([]*types.RichFeature, error)
	ListRichFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.RichFeature, error)

	// feature group
	CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt, category string) error
	GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error)
	ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error)
	UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) error

	// revision
	ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error)
}

var _ Store = &DB{}

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
