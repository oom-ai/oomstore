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

var (
	createSchemaFuncs = template.FuncMap{
		"entity": func(entity types.Entity, backend types.BackendType) string {
			switch backend {
			case types.CASSANDRA:
				return fmt.Sprintf(`"%s" TEXT PRIMARY KEY`, entity.Name)
			case types.POSTGRES, types.SNOWFLAKE:
				return fmt.Sprintf(`"%s" VARCHAR(%d) PRIMARY KEY`, entity.Name, entity.Length)
			case types.MYSQL:
				return fmt.Sprintf("`%s` VARCHAR(%d) PRIMARY KEY", entity.Name, entity.Length)
			default:
				return fmt.Sprintf("%s VARCHAR(%d) PRIMARY KEY", entity.Name, entity.Length)
			}
		},
		"columnJoin": func(columns []Column, backend types.BackendType) string {
			var format string
			switch backend {
			case types.POSTGRES, types.SNOWFLAKE, types.CASSANDRA:
				format = `"%s" %s`
			case types.MYSQL:
				format = "`%s` %s"
			default:
				format = "%s %s"
			}
			rs := make([]string, 0, len(columns))
			for _, column := range columns {
				rs = append(rs, fmt.Sprintf(format, column.Name, column.DbType))
			}
			return strings.Join(rs, ",\n\t")
		},
	}
)

type Column struct {
	Name   string
	DbType string
}

type CreateSchema struct {
	TableName string
	Entity    types.Entity
	Columns   []Column

	Backend types.BackendType
}

func BuildCreateSchema(tableName string, entity *types.Entity, features types.FeatureList, backend types.BackendType) (string, error) {
	var text string
	switch backend {
	case types.SNOWFLAKE:
		text = snowflakeCreateSchema
	default:
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
		dbType, err := GetDbTypeFrom(backend, feature.ValueType)
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
