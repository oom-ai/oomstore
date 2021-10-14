package onestore

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

const CREATE_DATA_TABLE = "CREATE TABLE {{TABLE_NAME}} (\n" +
	"`{{ENTITY_NAME}}` VARCHAR(32) COMMENT 'entity key' PRIMARY KEY,\n" +
	"{{COLUMN_DEFS}});"

func buildFeatureGroupSchema(columns []*types.Feature) string {
	// sort to ensure the schema looks consistent
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Name < columns[j].Name
	})
	var columnDefs []string
	for _, column := range columns {
		columnDef := fmt.Sprintf("`%s` %s COMMENT '%s'", column.Name, column.ValueType, column.Description)
		columnDefs = append(columnDefs, columnDef)
	}
	return strings.ReplaceAll(CREATE_DATA_TABLE, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))
}

func getCsvHeader(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	return header, nil
}

func hasDup(a []string) bool {
	s := make(map[string]bool)
	for _, e := range a {
		if s[e] {
			return true
		}
		s[e] = true
	}
	return false
}

func stringSliceEqual(a, b []string) bool {
	ma := make(map[string]bool)
	mb := make(map[string]bool)
	for _, e := range a {
		ma[e] = true
	}
	for _, e := range b {
		mb[e] = true
	}
	if len(ma) != len(mb) {
		return false
	}
	for k := range mb {
		if _, ok := ma[k]; !ok {
			return false
		}
	}
	return true
}

func (s *OneStore) ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt) error {
	// get columns of the group
	columns, err := s.db.ListFeature(ctx, &opt.GroupName)
	if err != nil {
		return err
	}

	// make sure csv data source has all defined columns
	header, err := getCsvHeader(opt.DataSource.FilePath)
	if err != nil {
		return err
	}
	if hasDup(header) {
		return fmt.Errorf("csv data source has duplicated columns: %v", header)
	}
	var columnNames []string
	for _, column := range columns {
		columnNames = append(columnNames, column.Name)
	}
	if !stringSliceEqual(header, columnNames) {
		return fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v", header, columnNames)
	}

	// create the data table
	schema := buildFeatureGroupSchema(columns)
	group, err := s.db.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	tmpTableName := opt.GroupName + "_" + strconv.Itoa(rand.Intn(100000))
	schema = strings.ReplaceAll(schema, "{{TABLE_NAME}}", tmpTableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", group.EntityName)
	_, err = s.db.ExecContext(ctx, schema)
	if err != nil {
		return err
	}

	// populate the data table
	err = s.db.LoadLocalFile(ctx, opt.DataSource.FilePath, tmpTableName, opt.DataSource.Separator, opt.DataSource.Delimiter, header)
	if err != nil {
		return err
	}

	// now get a timestamp
	ts := time.Now().Unix()

	// in a txn, rename the data table, insert into feature_group_revision, update feature_group
	err = database.WithTransaction(s.db, func(tx *sql.Tx) error {
		finalTableName := opt.GroupName + "_" + strconv.FormatInt(ts, 10)
		rename := fmt.Sprintf("RENAME `%s` TO `%s`", tmpTableName, finalTableName)
		if _, err = tx.ExecContext(ctx, rename); err != nil {
			return err
		}

		if err = database.InsertRevision(ctx, tx, opt.GroupName, ts, finalTableName, opt.Description); err != nil {
			return err
		}

		if err = database.UpdateFeatureGroup(ctx, tx, opt.GroupName, ts, finalTableName); err != nil {
			return err
		}

		return nil
	})

	return err
}
