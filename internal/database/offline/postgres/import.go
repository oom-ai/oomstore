package postgres

import (
	"context"
	"encoding/csv"
	"io"
	"time"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func loadData(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
	stmt, err := tx.PreparexContext(ctx, pq.CopyIn(tableName, header...))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		args := []interface{}{}
		for _, v := range row {
			args = append(args, v)
		}
		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	return err

}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	var revision int64
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		schema := dbutil.BuildFeatureDataTableSchema(opt.DataTableName, opt.Entity, opt.Features)
		_, err := tx.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		err = loadData(tx, ctx, opt.CsvReader, opt.DataTableName, opt.Header)
		if err != nil {
			return err
		}

		if opt.Revision != nil {
			// use user-defined revision
			revision = *opt.Revision
		} else {
			// generate revision using current timestamp
			revision = time.Now().UnixMilli()
		}
		return nil
	})
	return revision, err
}
