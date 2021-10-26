package database

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type RowMap = map[string]interface{}

const CREATE_DATA_TABLE = `CREATE TABLE {{TABLE_NAME}} (
	{{ENTITY_NAME}} VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY,
	{{COLUMN_DEFS}});
`

func BuildFeatureDataTableSchema(tableName string, entity *types.Entity, columns []*types.Feature) string {
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
