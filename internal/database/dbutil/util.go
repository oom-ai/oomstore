package dbutil

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type RowMap = map[string]interface{}

const CREATE_DATA_TABLE = `CREATE TABLE "{{TABLE_NAME}}" (
	"{{ENTITY_NAME}}" VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY,
	{{COLUMN_DEFS}});
`

func BuildFeatureDataTableSchema(tableName string, entity *types.Entity, features types.FeatureList) string {
	// sort to ensure the schema looks consistent
	sort.Slice(features, func(i, j int) bool {
		return features[i].Name < features[j].Name
	})
	var columnDefs []string
	for _, column := range features {
		columnDef := fmt.Sprintf(`"%s" %s`, column.Name, column.DBValueType)
		columnDefs = append(columnDefs, columnDef)
	}

	// fill schema template
	schema := strings.ReplaceAll(CREATE_DATA_TABLE, "{{TABLE_NAME}}", tableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", entity.Name)
	schema = strings.ReplaceAll(schema, "{{ENTITY_LENGTH}}", strconv.Itoa(entity.Length))
	schema = strings.ReplaceAll(schema, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))
	return schema
}

func BuildConditions(equal map[string]interface{}, in map[string]interface{}) ([]string, []interface{}, error) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	for key, value := range equal {
		cond = append(cond, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	for key, value := range in {
		s, inArgs, err := sqlx.In(fmt.Sprintf("%s IN (?)", key), value)
		if err != nil {
			return nil, nil, err
		}
		cond = append(cond, s)
		args = append(args, inArgs...)
	}
	return cond, args, nil
}

func Quote(quote string, fields ...string) string {
	var rs []string
	for _, f := range fields {
		rs = append(rs, quote+f+quote)
	}
	return strings.Join(rs, ",")
}

func TempTable(prefix string) string {
	return fmt.Sprintf("tmp_%s_%d", prefix, time.Now().UnixNano())
}

func buildQueryAndArgsForInsertRecords(tableName string, records []interface{}, columns []string) (string, []interface{}, error) {
	if len(records) == 0 {
		return "", nil, nil
	}
	valueFlags := make([]string, 0, len(records))
	for i := 0; i < len(records); i++ {
		valueFlags = append(valueFlags, "(?)")
	}

	return sqlx.In(
		fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES %s`, tableName, Quote(`"`, columns...), strings.Join(valueFlags, ",")),
		records...)
}

func InsertRecordsToTable(db *sqlx.DB, ctx context.Context, tableName string, records []interface{}, columns []string) error {
	query, args, err := buildQueryAndArgsForInsertRecords(tableName, records, columns)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}
	if _, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}

func InsertRecordsToTableTx(tx *sqlx.Tx, ctx context.Context, tableName string, records []interface{}, columns []string) error {
	query, args, err := buildQueryAndArgsForInsertRecords(tableName, records, columns)
	if err != nil {
		return err
	}
	if query == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, tx.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}

// Build MySQL data source name
func GetMysqlDSN(host, port, user, password, database string) string {
	cfg := mysql.NewConfig()
	cfg.Net = fmt.Sprintf("(%s:%s)", host, port)
	cfg.User = user
	cfg.Passwd = password
	cfg.DBName = database
	cfg.ParseTime = true
	return cfg.FormatDSN()
}
