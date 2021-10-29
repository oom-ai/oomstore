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
	if err := db.validateDataType(ctx, opt.DBValueType); err != nil {
		return fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	query := "INSERT INTO feature(name, group_name, db_value_type, value_type, description) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.ExecContext(ctx, query, opt.FeatureName, opt.GroupName, opt.DBValueType, opt.ValueType, opt.Description)
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
	var feature types.Feature
	query := `SELECT * FROM "rich_feature" WHERE name = $1`
	if err := db.GetContext(ctx, &feature, query, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func (db *DB) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	features := types.FeatureList{}
	query := `SELECT * FROM "rich_feature"`
	cond, args, err := buildListFeatureCond(opt)
	if err != nil {
		return nil, err
	}
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}
	if err := db.SelectContext(ctx, &features, db.Rebind(query), args...); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	query := "UPDATE feature SET description = $1 WHERE name = $2"
	if result, err := db.ExecContext(ctx, query, opt.NewDescription, opt.FeatureName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func buildListFeatureCond(opt types.ListFeatureOpt) (string, []interface{}, error) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	if opt.EntityName != nil {
		cond = append(cond, "entity_name = ?")
		args = append(args, *opt.EntityName)
	}
	if opt.GroupName != nil {
		cond = append(cond, "group_name = ?")
		args = append(args, *opt.GroupName)
	}
	if opt.FeatureNames != nil {
		if len(opt.FeatureNames) == 0 {
			return "1 = 0", nil, nil
		}
		s, inArgs, err := sqlx.In("name IN (?)", opt.FeatureNames)
		if err != nil {
			return "", nil, err
		}
		cond = append(cond, s)
		args = append(args, inArgs...)
	}
	return strings.Join(cond, " AND "), args, nil
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
