package dbutil

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

const (
	CREATE_DATA_TABLE_POSTGRES = `CREATE TABLE "{{TABLE_NAME}}" (
	"{{ENTITY_NAME}}" VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY,
	{{COLUMN_DEFS}});
`
	CREATE_DATA_TABLE_MYSQL = "CREATE TABLE `{{TABLE_NAME}}` ( " +
		"`{{ENTITY_NAME}}` VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY," +
		"{{COLUMN_DEFS}}); "
)

func BuildFeatureDataTableSchema(tableName string, entity *types.Entity, features types.FeatureList, backendType types.BackendType) (string, error) {
	var columnFormat, tableSchema string
	switch backendType {
	case types.POSTGRES:
		columnFormat = `"%s" %s`
		tableSchema = CREATE_DATA_TABLE_POSTGRES
	case types.MYSQL:
		columnFormat = "`%s` %s"
		tableSchema = CREATE_DATA_TABLE_MYSQL
	default:
		return "", fmt.Errorf("unsupported backend type %s", backendType)
	}

	// sort to ensure the schema looks consistent
	sort.Slice(features, func(i, j int) bool {
		return features[i].Name < features[j].Name
	})
	var columnDefs []string
	for _, column := range features {
		def := fmt.Sprintf(columnFormat, column.Name, column.DBValueType)
		columnDefs = append(columnDefs, def)
	}

	// fill schema template
	schema := strings.ReplaceAll(tableSchema, "{{TABLE_NAME}}", tableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", entity.Name)
	schema = strings.ReplaceAll(schema, "{{ENTITY_LENGTH}}", strconv.Itoa(entity.Length))
	schema = strings.ReplaceAll(schema, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))

	return schema, nil
}
