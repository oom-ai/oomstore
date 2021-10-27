package postgres

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/internal/database/offline"
)

func (db *DB) LoadLocalFile(ctx context.Context, filePath, tableName, delimiter string, header []string) error {
	return dbutil.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, pq.CopyIn(tableName, header...))
		if err != nil {
			return err
		}
		defer stmt.Close()

		dataFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer dataFile.Close()

		reader := csv.NewReader(dataFile)
		reader.Comma = []rune(delimiter)[0]

		// skip header
		_, err = reader.Read()
		if err != nil {
			return nil
		}

		for {
			row, err := reader.Read()
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

func (db *DB) ImportBatchFeatures(ctx context.Context, opt offline.ImportBatchFeaturesOpt) (int64, string, error) {
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
		err = db.LoadLocalFile(ctx, opt.DataSource.FilePath, tmpTableName, opt.DataSource.Delimiter, opt.Header)
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
