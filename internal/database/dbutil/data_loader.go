package dbutil

import (
	"context"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type LoadDataFromSourceOpt struct {
	Source     *offline.CSVSource
	EntityName string
	TableName  string
	Header     []string
	Features   types.FeatureList
	Backend    types.BackendType
}

// Currying
func LoadDataFromSource(backend types.BackendType, batchSize int) func(tx *sqlx.Tx, ctx context.Context, opt LoadDataFromSourceOpt) error {
	return func(tx *sqlx.Tx, ctx context.Context, opt LoadDataFromSourceOpt) error {
		return loadDataFromSource(tx, ctx, opt, batchSize)
	}
}

func loadDataFromSource(tx *sqlx.Tx, ctx context.Context, opt LoadDataFromSourceOpt, batchSize int) error {
	records := make([]interface{}, 0, batchSize)
	for {
		record, err := ReadLine(ReadLineOpt{
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
		records = append(records, record)
		if len(records) == batchSize {
			if err := InsertRecordsToTableTx(tx, ctx, opt.TableName, records, opt.Header, opt.Backend); err != nil {
				return err
			}
			records = make([]interface{}, 0, batchSize)
		}
	}
	if err := InsertRecordsToTableTx(tx, ctx, opt.TableName, records, opt.Header, opt.Backend); err != nil {
		return err
	}
	return nil
}

type ReadLineOpt struct {
	Source     *offline.CSVSource
	EntityName string
	Header     []string
	Features   types.FeatureList
}

func ReadLine(opt ReadLineOpt) ([]interface{}, error) {
	row, err := opt.Source.Reader.ReadString('\n')
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	rowSlice := strings.Split(strings.Trim(row, "\n"), opt.Source.Delimiter)
	line := make([]interface{}, 0, len(rowSlice))
	for i, ele := range rowSlice {
		if len(opt.Header) == 0 || len(opt.Features) == 0 || opt.Header[i] == opt.EntityName {
			// entity_key doesn't need to change type
			line = append(line, ele)
		} else if opt.Header[i] == "unix_milli" {
			line = append(line, castElement(ele, types.Int64))
		} else {
			feature := opt.Features.Find(func(f *types.Feature) bool {
				return f.Name == opt.Header[i]
			})
			line = append(line, castElement(ele, feature.ValueType))
		}
	}
	return line, nil
}

func castElement(s string, valueType types.ValueType) interface{} {
	if valueType != types.Bool {
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
