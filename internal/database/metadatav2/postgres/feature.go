package postgres

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateFeature(ctx context.Context, opt metadatav2.CreateFeatureOpt) (int16, error) {
	if err := db.validateDataType(ctx, opt.DBValueType); err != nil {
		return 0, fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	var featureId int16
	query := "INSERT INTO feature(name, group_name, db_value_type, value_type, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := db.GetContext(ctx, &featureId, query, opt.FeatureName, opt.GroupName, opt.DBValueType, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
	}
	return featureId, err
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	query := "UPDATE feature SET description = $1 WHERE name = $2"
	if result, err := db.ExecContext(ctx, query, opt.NewDescription, opt.FeatureName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (db *DB) validateDataType(ctx context.Context, dataType string) error {
	tmpTableName := fmt.Sprintf("tmp_validate_data_type_%d", rand.Int())
	return dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", tmpTableName)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE TABLE %s (a %s)", tmpTableName, dataType)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("DROP TABLE %s", tmpTableName)); err != nil {
			return err
		}
		return nil
	})
}
