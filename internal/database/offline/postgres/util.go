package postgres

import (
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
)

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, opt dbutil.LoadDataFromSourceOpt) error {
	stmt, err := tx.PreparexContext(ctx, pq.CopyIn(opt.TableName, opt.Header...))
	if err != nil {
		return errdefs.WithStack(err)
	}
	defer stmt.Close()

	for {
		record, err := dbutil.ReadLine(dbutil.ReadLineOpt{
			Source:     opt.Source,
			EntityName: opt.EntityName,
			Header:     opt.Header,
			Features:   opt.Features,
		})
		if errdefs.Cause(err) == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != len(opt.Header) {
			continue
		}
		if _, err := stmt.ExecContext(ctx, record...); err != nil {
			return errdefs.WithStack(err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	return errdefs.WithStack(err)
}
