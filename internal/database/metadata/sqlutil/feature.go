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

func UpdateFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateFeatureOpt) error {
	if opt.NewDescription == nil {
		return errdefs.Errorf("invalid option: nothing to update")
	}

	query := "UPDATE feature SET description = ? WHERE id = ?"
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.NewDescription, opt.FeatureID)
	if err != nil {
		return errdefs.WithStack(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errdefs.WithStack(err)
	}
	if rowsAffected != 1 {
		return errdefs.Errorf("failed to update feature %d: feature not found", opt.FeatureID)
	}
	return nil
}

func GetFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Feature, error) {
	var (
		feature types.Feature
		group   *types.Group
		err     error
	)

	query := `SELECT * FROM feature WHERE id = ?`
	if err := sqlxCtx.GetContext(ctx, &feature, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("feature %d not found", id))
		}
		return nil, errdefs.WithStack(err)
	}

	if group, err = GetGroup(ctx, sqlxCtx, feature.GroupID); err != nil {
		return nil, errdefs.WithStack(err)
	}
	feature.Group = group

	return &feature, nil
}

func GetFeatureByName(ctx context.Context, sqlxCtx metadata.SqlxContext, groupName string, featureName string) (*types.Feature, error) {
	var (
		feature types.Feature
		group   *types.Group
		err     error
	)

	if group, err = GetGroupByName(ctx, sqlxCtx, groupName); err != nil {
		return nil, errdefs.WithStack(err)
	}
	feature.Group = group

	query := `SELECT * FROM feature WHERE name = ? AND group_id = ?`
	if err := sqlxCtx.GetContext(ctx, &feature, sqlxCtx.Rebind(query), featureName, group.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("feature %s not found in group %s", featureName, groupName))
		}
		return nil, errdefs.WithStack(err)
	}

	return &feature, nil
}

func ListFeature(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	var (
		features types.FeatureList
		err      error
	)

	query := `SELECT * FROM feature`
	cond, args, err := buildListFeatureCond(opt)
	if err != nil {
		return nil, err
	}
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(cond, " AND "))
	}
	query = fmt.Sprintf("%s ORDER BY id ASC", query)
	if err := sqlxCtx.SelectContext(ctx, &features, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, errdefs.WithStack(err)
	}

	// enrich group
	groupIDs := features.GroupIDs()
	groups, err := ListGroup(ctx, sqlxCtx, metadata.ListGroupOpt{
		GroupIDs: &groupIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, f := range features {
		group := groups.Find(func(g *types.Group) bool { return g.ID == f.GroupID })
		if group == nil {
			return nil, errdefs.InvalidAttribute(errdefs.Errorf("no group found for feature %s", f.Name))
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

	return dbutil.BuildConditions(and, in)
}
