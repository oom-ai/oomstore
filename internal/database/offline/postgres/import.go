package postgres

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

const CREATE_DATA_TABLE = `CREATE TABLE {{TABLE_NAME}} (
	{{ENTITY_NAME}} VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY,
	{{COLUMN_DEFS}});
`

func buildFeatureDataTableSchema(tableName string, entity *types.Entity, columns []*types.Feature) string {
	// sort to ensure the schema looks consistent
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Name < columns[j].Name
	})
	var columnDefs []string
	for _, column := range columns {
		columnDef := fmt.Sprintf("%s %s", column.Name, column.ValueType)
		columnDefs = append(columnDefs, columnDef)
	}

	// fill schema template
	schema := strings.ReplaceAll(CREATE_DATA_TABLE, "{{TABLE_NAME}}", tableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", entity.Name)
	schema = strings.ReplaceAll(schema, "{{ENTITY_LENGTH}}", strconv.Itoa(entity.Length))
	schema = strings.ReplaceAll(schema, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))
	return schema
}

func (db *DB) LoadLocalFile(ctx context.Context, filePath, tableName, delimiter string, header []string) error {
	return database.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, pq.CopyIn(tableName, header...))
		if err != nil {
			return err
		}
		defer stmt.Close()

		dataFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer dataFile.Close()

		reader := csv.NewReader(dataFile)
		reader.Comma = []rune(delimiter)[0]

		// skip header
		_, err = reader.Read()
		if err != nil {
			return nil
		}

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			args := []interface{}{}
			for _, v := range row {
				args = append(args, v)
			}
			if _, err := stmt.ExecContext(ctx, args...); err != nil {
				return err
			}
		}

		_, err = stmt.ExecContext(ctx)
		return err
	})
}

func (db *DB) ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt, entity *types.Entity, features []*types.Feature, header []string) (int64, string, error) {
	var revision int64
	var finalTableName string
	err := database.WithTransaction(db.DB, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// create the data table
		tmpTableName := opt.GroupName + "_" + strconv.Itoa(rand.Int())
		schema := buildFeatureDataTableSchema(tmpTableName, entity, features)
		_, err := db.ExecContext(ctx, schema)
		if err != nil {
			return err
		}

		// populate the data table
		err = db.LoadLocalFile(ctx, opt.DataSource.FilePath, tmpTableName, opt.DataSource.Delimiter, header)
		if err != nil {
			return err
		}

		// generate revision using current timestamp
		revision = time.Now().Unix()

		// generate final data table name
		finalTableName = opt.GroupName + "_" + strconv.FormatInt(revision, 10)

		rename := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTableName, finalTableName)
		_, err = tx.ExecContext(ctx, rename)
		return err
	})
	return revision, finalTableName, err
}
