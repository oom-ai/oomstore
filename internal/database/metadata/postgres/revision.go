package postgres

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/internal/database/metadata"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision"
	var cond []interface{}
	if groupName != nil {
		query += " WHERE group_name = $1"
		cond = append(cond, *groupName)
	}
	revisions := make([]*types.Revision, 0)

	if err := db.SelectContext(ctx, &revisions, query, cond...); err != nil {
		return nil, err
	}
	return revisions, nil
}

func (db *DB) GetRevision(ctx context.Context, groupName string, revision int64) (*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision WHERE group_name = $1 and revision = $2"
	var rs types.Revision
	if err := db.GetContext(ctx, rs, query, groupName, revision); err != nil {
		return nil, err
	}
	return &rs, nil
}

func (db *DB) GetRevisionsByDataTables(ctx context.Context, dataTables []string) ([]*types.Revision, error) {
	query := "SELECT * FROM feature_group_revision WHERE data_table IN (?)"
	sql, args, err := sqlx.In(query, dataTables)
	if err != nil {
		return nil, err
	}

	revisions := make([]*types.Revision, 0)
	err = db.SelectContext(ctx, &revisions, db.Rebind(sql), args...)
	if err != nil {
		return nil, err
	}
	return revisions, nil
}

func (db *DB) InsertRevision(ctx context.Context, opt metadata.InsertRevisionOpt) error {
	query := "INSERT INTO feature_group_revision(group_name, revision, data_table, description) VALUES ($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, query, opt.GroupName, opt.Revision, opt.DataTable, opt.Description)
	return err
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
