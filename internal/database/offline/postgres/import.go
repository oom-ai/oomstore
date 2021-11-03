package postgres

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
)

func (db *DB) LoadData(ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
	return dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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
	})
}

func (db *DB) Import(ctx context.Context, opt offline.ImportOpt) (int64, string, error) {
	var revision int64
	var finalTableName string
	err := dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := opt.GroupName + "_" + strconv.Itoa(rand.Int())
		schema := dbutil.BuildFeatureDataTableSchema(tmpTableName, opt.Entity, opt.Features)
		_, err := db.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		err = db.LoadData(ctx, opt.CsvReader, tmpTableName, opt.Header)
		if err != nil {
			return err
		}

		// generate revision using current timestamp
		revision = time.Now().Unix()

		// generate final data table name
		finalTableName = opt.GroupName + "_" + strconv.FormatInt(revision, 10)

		rename := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return revision, finalTableName, err
}
