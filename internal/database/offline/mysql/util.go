package mysql

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

const (
	MySQLBatchSize = 20
)

func loadDataFromCSVReader(tx *sqlx.Tx, ctx context.Context, csvReader *csv.Reader, tableName string, header []string) error {
	records := make([]interface{}, 0, MySQLBatchSize)
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

		if len(records) == MySQLBatchSize {
			if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, header, types.MYSQL); err != nil {
				return err
			}
			records = make([]interface{}, 0, MySQLBatchSize)
		}
	}
	if err := dbutil.InsertRecordsToTableTx(tx, ctx, tableName, records, header, types.MYSQL); err != nil {
		return err
	}
	return nil

}
