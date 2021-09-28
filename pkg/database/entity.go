package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func ReadEntityTable(ctx context.Context, db *DB, tableName string, entityKeys, featureNames []string) (*sql.Rows, error) {
	// https://jmoiron.github.io/sqlx/#inQueries
	sql, args, err := sqlx.In(
		fmt.Sprintf("select entity_key, %s from %s where entity_key in (?);", strings.Join(featureNames, ", "), tableName),
		entityKeys,
	)
	if err != nil {
		return nil, err
	}

	sql = db.Rebind(sql)
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed connecting feature store: %v", err)
	}

	return rows, nil
}
