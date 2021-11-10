package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (db *DB) CreateRevision(ctx context.Context, opt metadatav2.CreateRevisionOpt) (*typesv2.Revision, error) {
	query := "INSERT INTO feature_group_revision(group_name, revision, data_table, anchored, description) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	var revision typesv2.Revision
	if err := db.GetContext(ctx, &revision, query, opt.GroupName, opt.Revision, opt.DataTable, opt.Anchored, opt.Description); err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("revision %v already exist", opt.Revision)
			}
		}
		return nil, err
	}
	return &revision, nil
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadatav2.UpdateRevisionOpt) (int64, error) {
	and := make(map[string]interface{})
	if opt.NewRevision != nil {
		and["revision"] = *opt.NewRevision
	}
	if opt.NewAnchored != nil {
		and["anchored"] = *opt.NewAnchored
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return 0, err
	}
	if len(cond) == 0 {
		return 0, fmt.Errorf("invliad option: nothing to update")
	}
	args = append(args, opt.RevisionID)

	query := fmt.Sprintf("UPDATE feature_group_revision SET %s WHERE id = ?", strings.Join(cond, ","))
	if result, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
