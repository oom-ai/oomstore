package postgres

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
)

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, opt dbutil.LoadDataFromSourceOpt) error {
	stmt, err := tx.PreparexContext(ctx, pq.CopyIn(opt.TableName, opt.Header...))
	if err != nil {
		return errdefs.WithStack(err)
	}
	defer stmt.Close()

	for {
		record, err2 := dbutil.ReadLine(dbutil.ReadLineOpt{
			Source:     opt.Source,
			EntityName: opt.EntityName,
			Header:     opt.Header,
			Features:   opt.Features,
		})
		if err2 != nil {
			if errdefs.Cause(err2) == io.EOF {
				break
			}
			return err2
		}
		if len(record) != len(opt.Header) {
			continue
		}
		if _, err2 := stmt.ExecContext(ctx, record...); err2 != nil {
			return errdefs.WithStack(err2)
		}
	}

	_, err = stmt.ExecContext(ctx)
	return errdefs.WithStack(err)
}
