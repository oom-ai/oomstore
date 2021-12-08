package postgres

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func loadDataFromCSVReader(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
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
