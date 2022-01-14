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

type DoExportOpt struct {
	offline.ExportOpt
	QueryResults QueryExportResults
}

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backend types.BackendType) (*types.ExportResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend: backend,
		SqlxDB:  db,
	}
	doJoinOpt := DoExportOpt{
		ExportOpt:    opt,
		QueryResults: queryExportResults,
	}
	return DoExport(ctx, dbOpt, doJoinOpt)
}

func DoExport(ctx context.Context, dbOpt dbutil.DBOpt, opt DoExportOpt) (*types.ExportResult, error) {
	var (
		emptyStream = make(chan types.ExportRecord)
		errs        = make(chan error, 1) // at most 1 error
	)
	defer close(emptyStream)
	defer close(errs)

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
		snapshotTables = append(snapshotTables, opt.SnapshotTables[groupID])
		if _, ok := opt.CdcTables[groupID]; ok {
			cdcTables = append(cdcTables, opt.CdcTables[groupID])
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
		fields      []string
		featureList types.FeatureList
	)
	for _, groupID := range groupIDs {
		features := opt.Features[groupID]
		if features[0].Group.Category == types.CategoryBatch {
			for _, f := range features {
				fields = append(fields, fmt.Sprintf("%s.%s AS %s", qt(opt.SnapshotTables[groupID]), qt(f.Name), qt(f.FullName)))
				featureList = append(featureList, f)
			}
		} else {
			for _, f := range features {
				cdc := fmt.Sprintf("%s.%s", opt.CdcTables[groupID]+"_0", qt(f.Name))
				snapshot := fmt.Sprintf("%s.%s", qt(opt.SnapshotTables[groupID]), qt(f.Name))
				fields = append(fields, fmt.Sprintf("(CASE WHEN %s IS NULL THEN %s ELSE %s END) AS %s", cdc, snapshot, cdc, qt(f.FullName)))
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
	})
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, 0, len(opt.CdcTables)*2)
	for i := 0; i < len(opt.CdcTables)*2; i++ {
		args = append(args, opt.UnixMilli)
	}
	return queryExportResults(ctx, dbOpt, opt.ExportOpt, query, args, featureList)
}

func queryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}, features types.FeatureList) (*types.ExportResult, error) {
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error

	go func() {
		defer close(stream)
		defer close(errs)
		stmt, err := dbOpt.SqlxDB.Preparex(dbOpt.SqlxDB.Rebind(query))
		if err != nil {
			errs <- errdefs.WithStack(err)
			return
		}
		defer stmt.Close()
		rows, err := stmt.Queryx(args...)
		if err != nil {
			errs <- errdefs.WithStack(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			record, err := rows.SliceScan()
			if err != nil {
				errs <- errdefs.Errorf("failed at rows.SliceScan, err=%v", err)
				return
			}
			record[0] = cast.ToString(record[0])
			for i, f := range features {
				if record[i+1] == nil {
					continue
				}
				deserializedValue, err := dbutil.DeserializeByValueType(record[i+1], f.ValueType, dbOpt.Backend)
				if err != nil {
					errs <- err
					return
				}
				record[i+1] = deserializedValue
			}
			stream <- record
		}
	}()
	header := append([]string{opt.EntityName}, features.FullNames()...)
	return types.NewExportResult(header, stream, errs), nil
}
