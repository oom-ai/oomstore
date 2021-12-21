package bigquery

import (
	"bytes"
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"strconv"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) { // Step 1: prepare temporary table entity_rows
	features := types.FeatureList{}
	for _, featureList := range opt.FeatureMap {
		features = append(features, featureList...)
	}
	if len(features) == 0 {
		return nil, nil
	}
	entityRowsTableName, err := createAndImportTableEntityRows(ctx, db, opt.EntityRows, opt.ValueNames, types.MYSQL)
	if err != nil {
		return nil, err
	}

	// Step 2: process features by group, insert result to table joined
	tableNames := make([]string, 0)
	tableToFeatureMap := make(map[string]types.FeatureList)
	for groupName, featureList := range opt.FeatureMap {
		revisionRanges, ok := opt.RevisionRangeMap[groupName]
		if !ok {
			continue
		}
		joinedTableName, err := joinOneGroup(ctx, db, offline.JoinOneGroupOpt{
			GroupName:           groupName,
			Entity:              opt.Entity,
			Features:            featureList,
			RevisionRanges:      revisionRanges,
			EntityRowsTableName: entityRowsTableName,
		}, types.MYSQL)
		if err != nil {
			return nil, err
		}
		if joinedTableName != "" {
			tableNames = append(tableNames, joinedTableName)
			tableToFeatureMap[joinedTableName] = featureList
		}
	}

	//// Step 3: read joined results
	return readJoinedTable(ctx, db, entityRowsTableName, tableNames, tableToFeatureMap, opt.ValueNames, types.MYSQL)
}

func joinOneGroup(ctx context.Context, db *DB, opt offline.JoinOneGroupOpt, backendType types.BackendType) (string, error) {
	if len(opt.Features) == 0 {
		return "", nil
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return "", err
	}
	entityKeyStr := qt("entity_key")
	unixMilliStr := qt("unix_milli")

	// Step 1: create temporary joined table
	joinedTableName, err := createTableJoined(ctx, db, opt.Features, opt.GroupName, opt.ValueNames, backendType)
	if err != nil {
		return "", err
	}

	// Step 2: iterate each table range, join entity_rows table and each data tables
	columns := append(opt.ValueNames, opt.Features.Names()...)
	columnsStr := qt(columns...)
	for _, r := range opt.RevisionRanges {
		q, err := buildJoinQuery(joinSchema{
			TableName:           fmt.Sprintf("%s.%s", db.datasetID, joinedTableName),
			EntityKeyStr:        entityKeyStr,
			EntityName:          opt.Entity.Name,
			UnixMilliStr:        unixMilliStr,
			ColumnsStr:          columnsStr,
			EntityRowsTableName: fmt.Sprintf("%s.%s", db.datasetID, opt.EntityRowsTableName),
			DataTable:           fmt.Sprintf("%s.%s", db.datasetID, r.DataTable),
			Backend:             backendType,
		})
		if err != nil {
			return "", err
		}
		q = strings.Replace(q, "?", strconv.Itoa(int(r.MinRevision)), 1)
		q = strings.Replace(q, "?", strconv.Itoa(int(r.MaxRevision)), 1)
		if _, err = db.Query(q).Read(ctx); err != nil {
			return "", err
		}
	}

	return joinedTableName, nil
}

func readJoinedTable(
	ctx context.Context,
	db *DB,
	entityRowsTableName string,
	tableNames []string,
	featureMap map[string]types.FeatureList,
	valueNames []string,
	backendType types.BackendType,
) (*types.JoinResult, error) {
	if len(tableNames) == 0 {
		return nil, nil
	}
	qt, err := dbutil.QuoteFn(backendType)
	if err != nil {
		return nil, err
	}
	entityKeyStr := qt("entity_key")
	unixMilliStr := qt("unix_milli")

	// Step 1: join temporary tables
	/*
		SELECT
		entity_rows_table.entity_key,
			entity_rows_table.unix_milli,
			joined_table_1.feature_1,
			joined_table_1.feature_2,
			joined_table_2.feature_3
		FROM entity_rows_table
		LEFT JOIN joined_table_1
		ON entity_rows_table.entity_key = joined_table_1.entity_key AND entity_rows_table.unix_milli = joined_table_1.unix_milli
		LEFT JOIN joined_table_2
		ON entity_rows_table.entity_key = joined_table_2.entity_key AND entity_rows_table.unix_milli = joined_table_2.unix_milli;
	*/
	var (
		header         []string
		fields         []string
		joinTablePairs []joinTablePair
	)

	header = append(header, "entity_key", "unix_milli")
	for _, name := range valueNames {
		fields = append(fields, qt(name))
		header = append(header, name)
	}
	for _, name := range tableNames {
		for _, f := range featureMap[name] {
			fields = append(fields, fmt.Sprintf("%s.%s", qt(name), qt(f.Name)))
			header = append(header, f.Name)
		}
	}

	tableNames = append([]string{entityRowsTableName}, tableNames...)
	for i := 0; i < len(tableNames)-1; i++ {
		joinTablePairs = append(joinTablePairs, joinTablePair{
			LeftTable:  tableNames[i],
			RightTable: tableNames[i+1],
		})
	}
	query, err := buildReadJoinResultQuery(readJoinResultQuery{
		EntityRowsTableName: entityRowsTableName,
		EntityKeyStr:        entityKeyStr,
		UnixMilliStr:        unixMilliStr,
		Fields:              fields,
		JoinTables:          joinTablePairs,
		Backend:             backendType,
		DatasetID:           db.datasetID,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: read joined results
	rows, err := db.Query(query).Read(ctx)
	if err != nil {
		return nil, err
	}

	data := make(chan []interface{})
	var scanErr, dropErr error

	//recordMap := make(map[string]bigquery.Value)
	//err = rows.Next(&recordMap)
	//if err == iterator.Done {
	//	return &types.JoinResult{}, nil
	//}
	//if err != nil {
	//	scanErr = err
	//}
	//schema := rows.Schema
	////header := make([]string, 0, len(schema))
	////for _, field := range schema {
	////	header = append(header, field.Name)
	////}
	//record := make([]interface{}, 0, len(recordMap))
	//for _, h := range header {
	//	record = append(record, recordMap[h])
	//}
	go func() {
		defer func() {
			if err := dropTemporaryTables(ctx, db, tableNames); err != nil {
				dropErr = err
			}
			defer close(data)
		}()
		for {
			recordMap := make(map[string]bigquery.Value)
			err = rows.Next(&recordMap)
			if err == iterator.Done {
				break
			}
			if err != nil {
				scanErr = err
			}
			record := make([]interface{}, 0, len(recordMap))
			for _, h := range header {
				record = append(record, recordMap[h])
			}
			data <- record
		}
	}()

	// TODO: return errors through channel
	if scanErr != nil {
		return nil, scanErr
	}
	return &types.JoinResult{
		Header: header,
		Data:   data,
	}, dropErr
}

func dropTemporaryTables(ctx context.Context, db *DB, tableNames []string) error {
	var err error
	for _, tableName := range tableNames {
		if tmpErr := dropTable(ctx, db, tableName); tmpErr != nil {
			err = tmpErr
		}
	}
	return err
}

const JOIN_SCHEMA = `
INSERT INTO {{ qt .TableName }} ( {{ .EntityKeyStr }}, {{ .UnixMilliStr }}, {{.ColumnsStr }})
SELECT
	l.{{ .EntityKeyStr }} AS entity_key,
	l.{{ .UnixMilliStr }} AS unix_milli,
	{{ .ColumnsStr }}
FROM
	{{ qt .EntityRowsTableName }} AS l
LEFT JOIN {{ qt .DataTable }} AS r
ON l.{{ .EntityKeyStr }} = r.{{ qt .EntityName }}
WHERE l.{{ .UnixMilliStr }} >= ? AND l.{{ .UnixMilliStr }} < ?
`

type joinSchema struct {
	TableName           string
	EntityKeyStr        string
	EntityName          string
	UnixMilliStr        string
	ColumnsStr          string
	EntityRowsTableName string
	DataTable           string
	Backend             types.BackendType
}

func buildJoinQuery(schema joinSchema) (string, error) {
	qt, err := dbutil.QuoteFn(schema.Backend)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("join").Funcs(template.FuncMap{
		"qt": qt,
	}).Parse(JOIN_SCHEMA))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, schema); err != nil {
		return "", err
	}
	return buf.String(), nil
}
