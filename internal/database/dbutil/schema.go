package dbutil

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const createSchema = `
CREATE TABLE {{ .TableName }} (
	{{ entity .Entity .Backend }},
	{{ columnJoin .Columns .Backend }}
)`

const snowflakeCreateSchema = `
CREATE TABLE "{{ .TableName }}" (
	{{ entity .Entity .Backend }},
	{{ columnJoin .Columns .Backend }}
)`

// TODO: Add back `PRIMARY KEY` back when we have functions for
// creating cdc table
var (
	createSchemaFuncs = template.FuncMap{
		"entity": func(entity types.Entity, backend types.BackendType) string {
			entityName := QuoteFn(backend)(entity.Name)
			switch backend {
			case types.BackendCassandra, types.BackendSQLite:
				return fmt.Sprintf(`%s TEXT`, entityName)
			case types.BackendPostgres, types.BackendRedshift, types.BackendSnowflake, types.BackendMySQL:
				return fmt.Sprintf(`%s VARCHAR(%d)`, entityName, entity.Length)
			case types.BackendBigQuery:
				return fmt.Sprintf(`%s STRING`, entityName)
			default:
				panic(fmt.Sprintf("unsupported backend type %s", backend))
			}
		},
		"columnJoin": func(columns []Column, backend types.BackendType) string {
			rs := make([]string, 0, len(columns))
			for _, column := range columns {
				rs = append(rs, fmt.Sprintf("%s %s", QuoteFn(backend)(column.Name), column.DbType))
			}
			return strings.Join(rs, ",\n\t")
		},
	}
)

type Column struct {
	Name      string
	DbType    string
	ValueType types.ValueType
}

type ColumnList []Column

func (c ColumnList) Names() []string {
	names := make([]string, 0, len(c))
	for _, column := range c {
		names = append(names, column.Name)
	}
	return names
}

type CreateSchema struct {
	TableName string
	Entity    types.Entity
	Columns   []Column

	Backend types.BackendType
}

func BuildCreateSchema(tableName string, entity *types.Entity, features types.FeatureList, backend types.BackendType) (string, error) {
	var text string
	if backend == types.BackendSnowflake {
		text = snowflakeCreateSchema
	} else {
		text = createSchema
	}
	buf := bytes.NewBuffer(nil)
	schema, err := newSchema(tableName, *entity, features, backend)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("schema").Funcs(createSchemaFuncs).Parse(text))
	if err = t.Execute(buf, schema); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newSchema(tableName string, entity types.Entity, features types.FeatureList, backend types.BackendType) (CreateSchema, error) {
	columns := make([]Column, 0, len(features))
	for _, feature := range features {
		dbType, err := DBValueType(backend, feature.ValueType)
		if err != nil {
			return CreateSchema{}, err
		}

		columns = append(columns, Column{
			Name:   feature.Name,
			DbType: dbType,
		})
	}

	return CreateSchema{
		TableName: tableName,
		Entity:    entity,
		Columns:   columns,
		Backend:   backend,
	}, nil
}
