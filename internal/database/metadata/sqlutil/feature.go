package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

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

	if group, err = GetGroup(ctx, sqlxCtx, feature.GroupID); err != nil {
		return nil, err
	}
	feature.Group = group

	return &feature, nil
}

func getFeatureByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Feature, error) {
	var (
		feature types.Feature
		group   *types.Group
		err     error
	)

	query := `SELECT * FROM "feature" WHERE name = $1`
	if err := sqlxCtx.GetContext(ctx, &feature, query, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature %s not found", name))
		}
		return nil, err
	}

	if group, err = GetGroup(ctx, sqlxCtx, feature.GroupID); err != nil {
		return nil, err
	}
	feature.Group = group

	return &feature, nil
}

func listFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	var (
		features types.FeatureList
		err      error
	)

	query := `SELECT * FROM "feature"`
	cond, args, err := buildListFeatureCond(opt)
	if err != nil {
		return nil, err
	}
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(cond, " AND "))
	}
	if err := sqlxCtx.SelectContext(ctx, &features, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, err
	}

	// enrich group
	groupIDs := features.GroupIDs()
	groups, err := ListGroup(ctx, sqlxCtx, nil, &groupIDs)
	if err != nil {
		return nil, err
	}
	for _, f := range features {
		group := groups.Find(func(g *types.Group) bool { return g.ID == f.GroupID })
		if group == nil {
			return nil, errdefs.InvalidAttribute(fmt.Errorf("no group found for feature %s", f.Name))
		}
		f.Group = group
	}

	// filter by entity
	if opt.EntityID != nil {
		features = features.Filter(func(f *types.Feature) bool {
			return f.Group.EntityID == *opt.EntityID
		})
	}
	return features, nil
}

func buildListFeatureCond(opt metadata.ListFeatureOpt) ([]string, []interface{}, error) {
	in := make(map[string]interface{})
	and := make(map[string]interface{})

	if opt.GroupID != nil {
		and["group_id"] = *opt.GroupID
	}

	if opt.FeatureIDs != nil {
		if len(*opt.FeatureIDs) == 0 {
			return []string{"false"}, nil, nil
		}
		in["id"] = *opt.FeatureIDs
	}

	if opt.FeatureNames != nil {
		if len(*opt.FeatureNames) == 0 {
			return []string{"false"}, nil, nil
		}
		in["name"] = *opt.FeatureNames
	}
	return dbutil.BuildConditions(and, in)
}
