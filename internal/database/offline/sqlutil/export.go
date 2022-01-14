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

type QueryExportResults func(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOneGroupOpt, query string, args []interface{}) (<-chan types.ExportRecord, <-chan error)

type DoExportOneGroupOpt struct {
	offline.ExportOneGroupOpt
	QueryResults QueryExportResults
}

type DoExportOpt struct {
	offline.ExportOpt
	QueryResults QueryExportResults
}

func ExportOneGroup(ctx context.Context, db *sqlx.DB, opt offline.ExportOneGroupOpt, backend types.BackendType) (<-chan types.ExportRecord, <-chan error) {
	dbOpt := dbutil.DBOpt{
		Backend: backend,
		SqlxDB:  db,
	}
	doJoinOpt := DoExportOneGroupOpt{
		ExportOneGroupOpt: opt,
		QueryResults:      sqlxQueryExportResults,
	}
	return DoExportOneGroup(ctx, dbOpt, doJoinOpt)
}

func DoExportOneGroup(ctx context.Context, dbOpt dbutil.DBOpt, opt DoExportOneGroupOpt) (<-chan types.ExportRecord, <-chan error) {
	var (
		emptyStream = make(chan types.ExportRecord)
		errs        = make(chan error, 1) // at most 1 error
	)
	defer close(emptyStream)
	defer close(errs)

	query, args, err := buildExportOneGroupQuery(dbOpt, opt.ExportOneGroupOpt)
	if err != nil {
		errs <- errdefs.WithStack(err)
		return emptyStream, errs
	}

	return opt.QueryResults(ctx, dbOpt, opt.ExportOneGroupOpt, query, args)
}

func Export(ctx context.Context, db *sqlx.DB, opt offline.ExportOpt, backend types.BackendType) (<-chan types.ExportRecord, <-chan error) {
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

func DoExport(ctx context.Context, dbOpt dbutil.DBOpt, opt DoExportOpt) (<-chan types.ExportRecord, <-chan error) {
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
		errs <- errdefs.WithStack(err)
		return emptyStream, errs
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
				fields = append(fields, fmt.Sprintf("%s.%s AS %s", qt(opt.SnapshotTables[groupID]), qt(f.Name), qt(f.FullName())))
				featureList = append(featureList, f)
			}
		} else {
			for _, f := range features {
				cdc := fmt.Sprintf("%s.%s", opt.CdcTables[groupID]+"_0", qt(f.Name))
				snapshot := fmt.Sprintf("%s.%s", qt(opt.SnapshotTables[groupID]), qt(f.Name))
				fields = append(fields, fmt.Sprintf("(CASE WHEN %s IS NULL THEN %s ELSE %s END) AS %s", cdc, snapshot, cdc, qt(f.FullName())))
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
		errs <- errdefs.WithStack(err)
		return emptyStream, errs
	}
	args := make([]interface{}, 0, len(opt.CdcTables)*2)
	for i := 0; i < len(opt.CdcTables)*2; i++ {
		args = append(args, opt.UnixMilli)
	}
	return queryExportResults(ctx, dbOpt, query, args, featureList)
}

func queryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, query string, args []interface{}, featureList types.FeatureList) (<-chan types.ExportRecord, <-chan error) {
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
			for i, f := range featureList {
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

	return stream, errs
}

func sqlxQueryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOneGroupOpt, query string, args []interface{}) (<-chan types.ExportRecord, <-chan error) {
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
			for i, f := range opt.Features {
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

	return stream, errs
}

func buildExportOneGroupQuery(dbOpt dbutil.DBOpt, opt offline.ExportOneGroupOpt) (string, []interface{}, error) {
	if opt.CdcTable == nil && opt.UnixMilli == nil {
		return buildExportBatchQuery(dbOpt, opt), nil, nil
	}
	if opt.CdcTable != nil && opt.UnixMilli != nil {
		return buildExportStreamQuery(dbOpt, opt)
	}
	return "", nil, fmt.Errorf("invalid option %+v", opt)
}

func buildExportBatchQuery(dbOpt dbutil.DBOpt, opt offline.ExportOneGroupOpt) string {
	fields := append([]string{opt.EntityName}, opt.Features.Names()...)
	qt := dbutil.QuoteFn(dbOpt.Backend)

	tableName := qt(opt.SnapshotTable)
	if dbOpt.Backend == types.BackendBigQuery {
		tableName = fmt.Sprintf("%s.%s", *dbOpt.DatasetID, tableName)
	}
	query := fmt.Sprintf("SELECT %s FROM %s", qt(fields...), tableName)
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query
}

func buildExportStreamQuery(dbOpt dbutil.DBOpt, opt offline.ExportOneGroupOpt) (string, []interface{}, error) {
	query, err := buildAggregateQuery(aggregateQueryParams{
		EntityName:            opt.EntityName,
		UnixMilli:             "unix_milli",
		FeatureNames:          opt.Features.Names(),
		PrevSnapshotTableName: opt.SnapshotTable,
		CurrCdcTableName:      *opt.CdcTable,
		Backend:               dbOpt.Backend,
		DatasetID:             dbOpt.DatasetID,
	})
	if err != nil {
		return "", nil, errdefs.WithStack(err)
	}
	if opt.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	return query, []interface{}{*opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli, *opt.UnixMilli}, nil
}
