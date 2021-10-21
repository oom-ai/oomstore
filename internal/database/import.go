package database

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (db *DB) LoadLocalFile(ctx context.Context, filePath, tableName, delimiter string, header []string) error {
	return db.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
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
