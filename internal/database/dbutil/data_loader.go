package dbutil

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

// Currying
func LoadDataFromCSVReader(backendType types.BackendType, batchSize int) func(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
	return func(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
		return loadDataFromCSVReader(tx, ctx, csvReader, tableName, header, backendType, batchSize)
	}
}

func loadDataFromCSVReader(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string, backendType types.BackendType, batchSize int) error {
	records := make([]interface{}, 0, batchSize)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		record := make([]interface{}, 0, len(row))
		for _, v := range row {
			record = append(record, v)
		}

		records = append(records, record)

		if len(records) == batchSize {
			if err := InsertRecordsToTableTx(tx, ctx, tableName, records, header, backendType); err != nil {
				return err
			}
			records = make([]interface{}, 0, batchSize)
		}
	}
	if err := InsertRecordsToTableTx(tx, ctx, tableName, records, header, backendType); err != nil {
		return err
	}
	return nil
}
