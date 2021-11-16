package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createFeatureTx(ctx context.Context, tx *sqlx.Tx, opt metadata.CreateFeatureOpt) (int16, error) {
	if err := validateDataType(ctx, tx, opt.DBValueType); err != nil {
		return 0, fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	var featureId int16
	query := "INSERT INTO feature(name, group_id, db_value_type, value_type, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := tx.GetContext(ctx, &featureId, query, opt.Name, opt.GroupID, opt.DBValueType, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("feature %s already exists", opt.Name)
			}
		}
	}
	return featureId, err
}

func updateFeatureTx(ctx context.Context, tx *sqlx.Tx, opt metadata.UpdateFeatureOpt) error {
	query := "UPDATE feature SET description = $1 WHERE id = $2"
	result, err := tx.ExecContext(ctx, query, opt.NewDescription, opt.FeatureID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("failed to update feature %d: feature not found", opt.FeatureID)
	}
	return nil
}

func validateDataType(ctx context.Context, tx *sqlx.Tx, dataType string) error {
	tmpTable := dbutil.TempTable("validate_data_type")
	stmt := fmt.Sprintf("CREATE TEMPORARY TABLE %s (a %s) ON COMMIT DROP", tmpTable, dataType)
	_, err := tx.ExecContext(ctx, stmt)
	return err
}
