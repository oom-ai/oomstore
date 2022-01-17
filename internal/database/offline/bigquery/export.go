package bigquery

import (
	"context"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"
	"google.golang.org/api/iterator"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (*types.ExportResult, error) {
	dbOpt := dbutil.DBOpt{
		Backend:    types.BackendBigQuery,
		BigQueryDB: db.Client,
		DatasetID:  &db.datasetID,
	}
	doExportOpt := sqlutil.DoExportOpt{
		ExportOpt:    opt,
		QueryResults: bigqueryQueryExportResults,
	}
	return sqlutil.DoExport(ctx, dbOpt, doExportOpt)
}

func bigqueryQueryExportResults(ctx context.Context, dbOpt dbutil.DBOpt, opt offline.ExportOpt, query string, args []interface{}, features types.FeatureList) (*types.ExportResult, error) {
	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error
	for _, arg := range args {
		query = strings.Replace(query, "?", cast.ToString(arg), 1)
	}

	go func() {
		defer close(stream)
		defer close(errs)
		rows, err := dbOpt.BigQueryDB.Query(query).Read(ctx)
		if err != nil {
			errs <- err
			return
		}
		for {
			recordMap := make(map[string]bigquery.Value)
			err = rows.Next(&recordMap)
			if err == iterator.Done {
				break
			}
			if err != nil {
				errs <- errdefs.Errorf("failed at rows.Next, err=%v", err)
				return
			}
			record := make([]interface{}, 0, len(recordMap))
			record = append(record, recordMap[opt.EntityName])
			for _, feature := range features {
				record = append(record, recordMap[feature.DBFullName(Backend)])
			}
			stream <- record
		}
	}()
	header := append([]string{opt.EntityName}, features.FullNames()...)
	return types.NewExportResult(header, stream, errs), nil
}
