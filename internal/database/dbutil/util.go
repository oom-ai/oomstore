package dbutil

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type RowMap = map[string]interface{}

const CREATE_DATA_TABLE = `CREATE TABLE {{TABLE_NAME}} (
	{{ENTITY_NAME}} VARCHAR({{ENTITY_LENGTH}}) PRIMARY KEY,
	{{COLUMN_DEFS}});
`

func BuildFeatureDataTableSchema(tableName string, entity *types.Entity, features types.FeatureList) string {
	// sort to ensure the schema looks consistent
	sort.Slice(features, func(i, j int) bool {
		return features[i].Name < features[j].Name
	})
	var columnDefs []string
	for _, column := range features {
		columnDef := fmt.Sprintf("%s %s", column.Name, column.DBValueType)
		columnDefs = append(columnDefs, columnDef)
	}

	// fill schema template
	schema := strings.ReplaceAll(CREATE_DATA_TABLE, "{{TABLE_NAME}}", tableName)
	schema = strings.ReplaceAll(schema, "{{ENTITY_NAME}}", entity.Name)
	schema = strings.ReplaceAll(schema, "{{ENTITY_LENGTH}}", strconv.Itoa(entity.Length))
	schema = strings.ReplaceAll(schema, "{{COLUMN_DEFS}}", strings.Join(columnDefs, ",\n"))
	return schema
}
