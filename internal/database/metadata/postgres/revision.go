package postgres

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) ListRevision(ctx context.Context, opt metadata.ListRevisionOpt) ([]*types.Revision, error) {
	var revisions []*types.Revision
	query := "SELECT * FROM feature_group_revision"
	cond, args, err := buildListRevisionCond(opt)
	if err != nil {
		return nil, err
	}
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(cond, " AND "))
	}
	if err := db.SelectContext(ctx, &revisions, db.Rebind(query), args...); err != nil {
		return nil, err
	}
	return revisions, nil
}

func buildListRevisionCond(opt metadata.ListRevisionOpt) ([]string, []interface{}, error) {
	and := make(map[string]interface{})
	in := make(map[string]interface{})

	if opt.GroupName != nil {
		and["group_name"] = *opt.GroupName
	}
	if opt.DataTables != nil {
		if len(opt.DataTables) == 0 {
			return []string{"false"}, nil, nil
		}
		in["data_table"] = opt.DataTables
	}
	return dbutil.BuildConditions(and, in)
}

func (db *DB) GetRevision(ctx context.Context, opt metadata.GetRevisionOpt) (*types.Revision, error) {
	and := make(map[string]interface{})
	if opt.GroupName != nil {
		and["group_name"] = *opt.GroupName
	}
	if opt.Revision != nil {
		and["revision"] = *opt.Revision
	}
	if opt.RevisionId != nil {
		and["id"] = *opt.RevisionId
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM feature_group_revision"
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(cond, " AND "))
	}

	var rs types.Revision
	if err := db.GetContext(ctx, &rs, db.Rebind(query), args...); err != nil {
		return nil, err
	}
	return &rs, nil
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (*types.Revision, error) {
	query := "INSERT INTO feature_group_revision(group_name, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	var revision types.Revision

	err := db.GetContext(ctx, &revision, query, opt.GroupName, opt.Revision, opt.DataTable, opt.Anchored, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("revision %v already exist", opt.Revision)
			}
		}
		return nil, err
	}
	return &revision, nil
}

func (db *DB) GetLatestRevision(ctx context.Context, groupName string) (*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision WHERE group_name = $1 ORDER BY create_time DESC LIMIT 1"
	var revision types.Revision
	if err := db.GetContext(ctx, &revision, query, groupName); err != nil {
		return nil, err
	}
	return &revision, nil
}

func (db *DB) BuildRevisionRanges(ctx context.Context, groupName string) ([]*types.RevisionRange, error) {
	query := fmt.Sprintf(`
		SELECT
			revision AS min_revision,
			LEAD(revision, 1, %d) OVER (ORDER BY revision) AS max_revision,
			data_table
		FROM feature_group_revision
		WHERE group_name = $1
	`, math.MaxInt64)

	var ranges []*types.RevisionRange
	if err := db.SelectContext(ctx, &ranges, query, groupName); err != nil {
		return nil, err
	}
	return ranges, nil
}
