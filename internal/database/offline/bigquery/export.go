package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"google.golang.org/api/iterator"
)

func (db *DB) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	fieldStr := dbutil.Quote("`", append([]string{opt.EntityName}, opt.Features.Names()...)...)
	tableName := dbutil.Quote("`", opt.DataTable)

	q := fmt.Sprintf(`SELECT %s FROM %s.%s`, fieldStr, db.datasetID, tableName)
	if opt.Limit != nil {
		q += fmt.Sprintf(" LIMIT %d", *opt.Limit)
	}
	query := db.Query(q)

	stream := make(chan types.ExportRecord)
	errs := make(chan error, 1) // at most 1 error
	go func() {
		defer close(stream)
		defer close(errs)
		rows, err := query.Read(ctx)
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
				errs <- fmt.Errorf("failed at rows.Next, err=%v", err)
				return
			}
			record := make([]interface{}, 0, len(recordMap))
			record = append(record, recordMap[opt.EntityName])
			for _, feature := range opt.Features {
				record = append(record, recordMap[feature.Name])
			}
			stream <- record
		}
	}()

	return stream, errs

}
