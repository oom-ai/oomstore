package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateFeatureOpt) (int, error) {
	if err := validateDataType(ctx, sqlxCtx, opt.DBValueType); err != nil {
		return 0, fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	var featureID int
	query := "INSERT INTO feature(name, group_id, db_value_type, value_type, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := sqlxCtx.GetContext(ctx, &featureID, query, opt.FeatureName, opt.GroupID, opt.DBValueType, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
	}
	return featureID, err
}

func updateFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateFeatureOpt) error {
	if opt.NewDescription == nil {
		return fmt.Errorf("invalid option: nothing to update")
	}

	query := "UPDATE feature SET description = $1 WHERE id = $2"
	result, err := sqlxCtx.ExecContext(ctx, query, opt.NewDescription, opt.FeatureID)
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

func validateDataType(ctx context.Context, sqlxCtx metadata.SqlxContext, dataType string) error {
	tmpTable := dbutil.TempTable("validate_data_type")
	stmt := fmt.Sprintf("CREATE TEMPORARY TABLE %s (a %s) ON COMMIT DROP", tmpTable, dataType)
	_, err := sqlxCtx.ExecContext(ctx, stmt)
	return err
}

func getFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Feature, error) {
	var (
		feature types.Feature
		group   *types.Group
		err     error
	)

	query := `SELECT * FROM "feature" WHERE id = $1`
	if err := sqlxCtx.GetContext(ctx, &feature, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature %d not found", id))
		}
		return nil, err
	}

	if group, err = getGroup(ctx, sqlxCtx, feature.GroupID); err != nil {
		return nil, err
	}
	feature.Group = group

	return &feature, nil
}
