package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
)

func (db *DB) CreateRevision(ctx context.Context, opt metadatav2.CreateRevisionOpt) error {
	query := "INSERT INTO feature_group_revision(group_name, revision, data_table, description) VALUES ($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, query, opt.GroupName, opt.Revision, opt.DataTable, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("revision %v already exist", opt.Revision)
			}
		}
	}
	return err
}
