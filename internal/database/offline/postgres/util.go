package postgres

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string, features types.FeatureList) error {
	stmt, err := tx.PreparexContext(ctx, pq.CopyIn(tableName, header...))
	if err != nil {
		return errdefs.WithStack(err)
	}
	defer stmt.Close()

	for {
		record, err := dbutil.ReadLine(source.Reader, source.Delimiter, features, Backend)
		if errdefs.Cause(err) == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != len(header) {
			continue
		}
		if _, err := stmt.ExecContext(ctx, record...); err != nil {
			return errdefs.WithStack(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	return errdefs.WithStack(err)
}
