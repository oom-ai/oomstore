package dbutil

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ethhte88/oomstore/pkg/errdefs"
)

const cassandraSchema = `
CREATE TABLE {{ .TableName }} (
    {{ .EntityName }}  TEXT PRIMARY KEY,
{{- range $i, $column := .Columns }}
    {{- if $i -}} , {{- end }}
    {{ $column.Name }}  {{$column.DbType}}
{{- end }}
)
`

type SchemaType string

const (
	Cassandra SchemaType = "cassandra"
)

type Column struct {
	Name   string
	DbType string
}

type Schema struct {
	TableName  string
	EntityName string
	Columns    []Column
}

func BuildSchema(schema Schema, schemaType SchemaType) (string, error) {
	var text string
	switch schemaType {
	case Cassandra:
		text = cassandraSchema
	default:
		return "", errdefs.InvalidAttribute(fmt.Errorf("schema type %s is not supported", schemaType))
	}

	buf := bytes.NewBuffer(nil)

	t := template.Must(template.New("schema").Parse(text))
	if err := t.Execute(buf, schema); err != nil {
		return "", err
	}

	return buf.String(), nil
}
