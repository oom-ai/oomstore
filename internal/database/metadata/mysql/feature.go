package mysql

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateFeatureOpt) (int, error) {
	if err := opt.ValueType.Validate(); err != nil {
		return 0, err
	}
	query := "INSERT INTO feature(name, group_id, value_type, description) VALUES (?, ?, ?, ?)"
	res, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.FeatureName, opt.GroupID, opt.ValueType, opt.Description)
	if err != nil {
		if er, ok := err.(*mysql.MySQLError); ok {
			if er.Number == ER_DUP_ENTRY {
				return 0, fmt.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
		return 0, err
	}

	featureID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(featureID), err
}

func validateDataType(ctx context.Context, sqlxCtx metadata.SqlxContext, dataType string) error {
	tmpTable := dbutil.TempTable("validate_data_type")
	query := fmt.Sprintf("CREATE TEMPORARY TABLE %s (a %s)", tmpTable, dataType)
	if _, err := sqlxCtx.ExecContext(ctx, query); err != nil {
		return err
	}
	_, err := sqlxCtx.ExecContext(ctx, fmt.Sprintf("DROP TEMPORARY TABLE %s", tmpTable))
	return err
}
