package sqlutil

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backendType types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	fields := append([]string{opt.EntityName}, opt.Features.Names()...)
	var fieldStr string
	var tableName string
	switch backendType {
	case types.POSTGRES, types.SNOWFLAKE:
		fieldStr = dbutil.Quote(`"`, fields...)
		tableName = dbutil.Quote(`"`, opt.DataTable)
	case types.MYSQL, types.SQLite:
		fieldStr = dbutil.Quote("`", fields...)
		tableName = dbutil.Quote("`", opt.DataTable)
	}
	query := fmt.Sprintf("SELECT %s FROM %s", fieldStr, tableName)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error

	go func() {
		defer close(stream)
		defer close(errs)
		stmt, err := db.Preparex(query)
		if err != nil {
			errs <- err
			return
		}
		defer stmt.Close()
		rows, err := stmt.Queryx()
		if err != nil {
			errs <- err
			return
		}
		defer rows.Close()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				errs <- fmt.Errorf("failed at rows.SliceScan, err=%v", err)
				return
			}
			record[0] = cast.ToString(record[0])
			for i, f := range opt.Features {
				if record[i+1] == nil {
					continue
				}
				if backendType == types.SNOWFLAKE {

					v, err := deserializeByTag(record[i+1], f.ValueType)
					if err != nil {
						errs <- fmt.Errorf("failed at deserializeByTag, err=%v", err)
						return
					}
					record[i+1] = v
				} else {
					if f.ValueType == types.STRING {
						record[i+1] = cast.ToString(record[i+1])
					}
				}
			}
			stream <- record
		}
	}()

	return stream, errs
}

// gosnowflake Scan always produce string when the destination is interface{}
// As a work around, we cast the string to interface{} based on ValueType
// This method is mostly copied from redis.DeserializeByTag, except we use 10 rather than 36 as the base
// In the long run, we should fix the gosnowflake converter with a pr
func deserializeByTag(i interface{}, typeTag string) (interface{}, error) {
	if i == nil {
		return nil, nil
	}

	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("not a string or nil: %v", i)
	}

	switch typeTag {
	case types.STRING:
		return s, nil

	case types.INT64:
		x, err := strconv.ParseInt(s, 10, 64)
		return x, err

	case types.FLOAT64:
		x, err := strconv.ParseFloat(s, 64)
		return x, err

	case types.BOOL:
		if s == "1" {
			return true, nil
		} else if s == "0" {
			return false, nil
		} else {
			return nil, fmt.Errorf("invalid bool value: %s", s)
		}
	case types.TIME:
		x, err := strconv.ParseInt(s, 10, 64)
		return time.UnixMilli(x), err

	case types.BYTES:
		return []byte(s), nil
	default:
		return "", fmt.Errorf("unsupported type tag: %s", typeTag)
	}
}
