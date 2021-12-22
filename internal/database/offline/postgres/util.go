package postgres

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string) error {
	stmt, err := tx.PreparexContext(ctx, pq.CopyIn(tableName, header...))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for {
		record, err := dbutil.ReadLine(source.Reader, source.Delimiter)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != len(header) {
			continue
		}
		args := make([]interface{}, 0, len(record))
		for _, v := range record {
			args = append(args, v)
		}
		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	return err
}
