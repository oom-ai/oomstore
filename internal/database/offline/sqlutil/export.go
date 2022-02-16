package sqlutil

import (
	"context"
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type QueryExportResults func(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}, features types.FeatureList) (*types.ExportResult, error)

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backend types.BackendType) (*types.ExportResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend: backend,
		SqlxDB:  db,
	}
	doJoinOpt := DoExportOpt{
		ExportOpt:    opt,
		QueryResults: sqlxQueryExportResults,
	}
	return DoExport(ctx, dbOpt, doJoinOpt)
}

type DoExportOpt struct {
	offline.ExportOpt
	QueryResults QueryExportResults
}

func DoExport(ctx context.Context, dbOpt dbutil.DBOpt, opt DoExportOpt) (*types.ExportResult, error) {
	// Step 0: prepare variables
	var (
		snapshotTables = make([]string, 0, len(opt.SnapshotTables))
		cdcTables      = make([]string, 0, len(opt.CdcTables))
		groupIDs       = make([]int, 0, len(opt.Features))
	)
	for groupID := range opt.Features {
		groupIDs = append(groupIDs, groupID)
	}
	sort.Slice(groupIDs, func(i, j int) bool {
		return groupIDs[i] < groupIDs[j]
	})
	for _, groupID := range groupIDs {
		snapshotTables = append(snapshotTables, buildTableName(dbOpt, opt.SnapshotTables[groupID]))
		if _, ok := opt.CdcTables[groupID]; ok {
			cdcTables = append(cdcTables, buildTableName(dbOpt, opt.CdcTables[groupID]))
		}
	}

	// Step 1: prepare export_entity table, which contains all entity keys from source tables
	tableName, err := prepareEntityTable(ctx, dbOpt, opt.ExportOpt, snapshotTables, cdcTables)
	if err != nil {
		return nil, err
	}

	// Step 2: join export_entity table, snapshot tables and cdc tables
	qt := dbutil.QuoteFn(dbOpt.Backend)
	var (
		j           int
		fields      []string
		featureList types.FeatureList
	)
	for i, groupID := range groupIDs {
		features := opt.Features[groupID]
		if features[0].Group.Category == types.CategoryBatch {
			for _, f := range features {
				fields = append(fields, fmt.Sprintf("%s.%s AS %s", snapshotTables[i], qt(f.Name), qt(f.DBFullName(dbOpt.Backend))))
				featureList = append(featureList, f)
			}
		} else {
			cdc := cdcTables[j] + "_0"
			j++
			for _, f := range features {
				fields = append(fields, fmt.Sprintf("(CASE WHEN %s.%s IS NULL THEN %s.%s ELSE %s.%s END) AS %s", qt(cdc), qt(f.Name), qt(snapshotTables[i]), qt(f.Name), qt(cdc), qt(f.Name), qt(f.DBFullName(dbOpt.Backend))))
				featureList = append(featureList, f)
			}
		}
	}
	query, err := buildExportQuery(exportQueryParams{
		EntityTableName: tableName,
		EntityName:      opt.EntityName,
		UnixMilli:       "unix_milli",
		SnapshotTables:  snapshotTables,
		CdcTables:       cdcTables,
		Fields:          fields,
		Backend:         dbOpt.Backend,
		DatasetID:       dbOpt.DatasetID,
	})
	if err != nil {
		return nil, err
	}
	if opt.Limit != nil {
		query = fmt.Sprintf("%s LIMIT %d", query, *opt.Limit)
	}
	args := make([]interface{}, 0, len(opt.CdcTables)*2)
	for i := 0; i < len(opt.CdcTables)*2; i++ {
		args = append(args, opt.UnixMilli)
	}
	return opt.QueryResults(ctx, dbOpt, opt.ExportOpt, query, args, featureList)
}

func sqlxQueryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}, features types.FeatureList) (*types.ExportResult, error) {
	stream := make(chan types.ExportRecord)

	go func() {
		defer close(stream)
		stmt, err := dbOpt.SqlxDB.Preparex(dbOpt.SqlxDB.Rebind(query))
		if err != nil {
			select {
			case stream <- types.ExportRecord{Error: errdefs.WithStack(err)}:
				return
			case <-ctx.Done():
				return
			}
		}
		defer stmt.Close()

		rows, err := stmt.Queryx(args...)
		if err != nil {
			select {
			case stream <- types.ExportRecord{Error: errdefs.WithStack(err)}:
				return
			case <-ctx.Done():
				return
			}
		}
		defer rows.Close()

		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				select {
				case stream <- types.ExportRecord{Error: errdefs.Errorf("failed at rows.SliceScan, err=%v", err)}:
					return
				case <-ctx.Done():
					return
				}
			}
			record[0] = cast.ToString(record[0])
			for i, f := range features {
				if record[i+1] == nil {
					continue
				}
				deserializedValue, err := dbutil.DeserializeByValueType(record[i+1], f.ValueType, dbOpt.Backend)
				if err != nil {
					select {
					case stream <- types.ExportRecord{Error: err}:
						return
					case <-ctx.Done():
						return
					}
				}
				record[i+1] = deserializedValue
			}
			select {
			case stream <- types.ExportRecord{Record: record, Error: nil}:
				// nothing to do
			case <-ctx.Done():
				return
			}
		}
	}()
	header := append([]string{opt.EntityName}, features.FullNames()...)
	return types.NewExportResult(header, stream), nil
}
