package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) error {
	return createFeature(ctx, db, opt)
}

func (tx *Tx) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) error {
	return createFeature(ctx, tx, opt)
}

func createFeature(ctx context.Context, ext metadata.ExtContext, opt metadata.CreateFeatureOpt) error {
	if err := validateDataType(ctx, ext, opt.DBValueType); err != nil {
		return fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	query := "INSERT INTO feature(name, group_name, db_value_type, value_type, description) VALUES ($1, $2, $3, $4, $5)"
	_, err := ext.ExecContext(ctx, query, opt.FeatureName, opt.GroupName, opt.DBValueType, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("feature %s already exists", opt.FeatureName)
			}
		}
	}
	return err
}

func (db *DB) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	return getFeature(ctx, db, featureName)
}

func (tx *Tx) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	return getFeature(ctx, tx, featureName)
}

func getFeature(ctx context.Context, ext metadata.ExtContext, featureName string) (*types.Feature, error) {
	var feature types.Feature
	query := `SELECT * FROM "rich_feature" WHERE name = $1`
	if err := ext.GetContext(ctx, &feature, query, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func (db *DB) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	return listFeature(ctx, db, opt)
}

func (tx *Tx) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	return listFeature(ctx, tx, opt)
}

func listFeature(ctx context.Context, ext metadata.ExtContext, opt types.ListFeatureOpt) (types.FeatureList, error) {
	features := types.FeatureList{}
	query := `SELECT * FROM "rich_feature"`
	cond, args, err := buildListFeatureCond(opt)
	if err != nil {
		return nil, err
	}
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(cond, " AND "))
	}
	if err := ext.SelectContext(ctx, &features, ext.Rebind(query), args...); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	return updateFeature(ctx, db, opt)
}

func (tx *Tx) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	return updateFeature(ctx, tx, opt)
}

func updateFeature(ctx context.Context, ext metadata.ExtContext, opt types.UpdateFeatureOpt) (int64, error) {
	query := "UPDATE feature SET description = $1 WHERE name = $2"
	if result, err := ext.ExecContext(ctx, query, opt.NewDescription, opt.FeatureName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func buildListFeatureCond(opt types.ListFeatureOpt) ([]string, []interface{}, error) {
	and := make(map[string]interface{})
	in := make(map[string]interface{})
	if opt.EntityName != nil {
		and["entity_name"] = *opt.EntityName
	}
	if opt.GroupName != nil {
		and["group_name"] = *opt.GroupName
	}
	if opt.FeatureNames != nil {
		if len(opt.FeatureNames) == 0 {
			return []string{"1 = 0"}, nil, nil
		}
		in["name"] = opt.FeatureNames
	}
	return dbutil.BuildConditions(and, in)
}

func validateDataType(ctx context.Context, ext metadata.ExtContext, dataType string) error {
	tmpTableName := fmt.Sprintf("tmp_validate_data_type_%d", rand.Int())
	return ext.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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
