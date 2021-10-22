package metadata

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) validateDataType(ctx context.Context, dataType string) error {
	tmpTableName := fmt.Sprintf("tmp_validate_data_type_%d", rand.Intn(100000))
	return db.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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

func (db *DB) CreateFeature(ctx context.Context, opt types.CreateFeatureOpt) error {
	if err := db.validateDataType(ctx, opt.ValueType); err != nil {
		return fmt.Errorf("err when validating value_type input, details: %s", err.Error())
	}
	query := "INSERT INTO feature(name, group_name, value_type, description) VALUES ($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, query, opt.FeatureName, opt.GroupName, opt.ValueType, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.DuplicateColumn {
				return fmt.Errorf("feature %s already exist", opt.FeatureName)
			}
		}
	}
	return err
}

func (db *DB) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	var feature types.Feature
	query := `SELECT * FROM feature WHERE name = $1`
	if err := db.GetContext(ctx, &feature, query, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func (db *DB) ListFeature(ctx context.Context, groupName *string) ([]*types.Feature, error) {
	query := "SELECT * FROM feature"
	cond, args := buildListFeatureCond(types.ListFeatureOpt{
		GroupName: groupName,
	})
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}

	features := make([]*types.Feature, 0)
	if err := db.SelectContext(ctx, &features, query, args...); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	query := "UPDATE feature SET description = $1 WHERE name = $2"
	_, err := db.ExecContext(ctx, query, opt.NewDescription, opt.FeatureName)
	return err
}

func buildListFeatureCond(opt types.ListFeatureOpt) (string, []interface{}) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	var id int
	if opt.EntityName != nil {
		id++
		cond = append(cond, fmt.Sprintf("entity_name = $%d", id))
		args = append(args, *opt.EntityName)
	}
	if opt.GroupName != nil {
		id++
		cond = append(cond, fmt.Sprintf("group_name = $%d", id))
		args = append(args, *opt.GroupName)
	}
	return strings.Join(cond, " AND "), args
}
