package dbutil

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Currying
func LoadDataFromSource(backendType types.BackendType, batchSize int) func(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string) error {
	return func(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string) error {
		return loadDataFromSource(tx, ctx, source, tableName, header, backendType, batchSize)
	}
}

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, source *offline.CSVSource, tableName string, header []string, backendType types.BackendType, batchSize int) error {
	records := make([]interface{}, 0, batchSize)
	for {
		record, err := ReadLine(source.Reader, source.Delimiter, backendType)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != len(header) {
			continue
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

func ReadLine(reader *bufio.Reader, delimiter string, backend types.BackendType) ([]interface{}, error) {
	row, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	rowSlice := strings.Split(strings.Trim(row, "\n"), delimiter)
	line := make([]interface{}, 0, len(rowSlice))
	for _, ele := range rowSlice {
		line = append(line, castElement(ele, backend))
	}
	return line, nil
}

func castElement(s string, backend types.BackendType) interface{} {
	if backend != types.BackendMySQL {
		return s
	}
	if s == "true" || s == "TRUE" {
		return 1
	}
	if s == "false" || s == "FALSE" {
		return 0
	}
	return s
}
