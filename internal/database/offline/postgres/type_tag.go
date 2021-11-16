package postgres

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) TypeTag(dbType string) (string, error) {
	return TypeTag(dbType)
}

func TypeTag(dbType string) (string, error) {
	var s = dbType
	if pos := strings.Index(dbType, "("); pos != -1 {
		s = s[:pos]
	}
	s = strings.TrimSpace(strings.ToLower(s))
	if t, ok := typeMap[s]; !ok {
		return "", fmt.Errorf("unsupported sql type: %s", dbType)
	} else {
		return t, nil
	}
}

var (
	typeMap = map[string]string{
		"bigint":    types.INT64,
		"int8":      types.INT64,
		"bigserial": types.INT64,
		"serial8":   types.INT64,

		"boolean": types.BOOL,
		"bool":    types.BOOL,

		"bytea":       types.BYTE_ARRAY,
		"jsonb":       types.BYTE_ARRAY,
		"uuid":        types.BYTE_ARRAY,
		"bit":         types.BYTE_ARRAY,
		"bit varying": types.BYTE_ARRAY,
		"character":   types.BYTE_ARRAY,
		"char":        types.BYTE_ARRAY,
		"json":        types.BYTE_ARRAY,
		"money":       types.BYTE_ARRAY,
		"numeric":     types.BYTE_ARRAY,

		"character varying": types.STRING,
		"text":              types.STRING,
		"varchar":           types.STRING,

		"double precision": types.FLOAT64,
		"float8":           types.FLOAT64,

		"integer": types.INT64,
		"int":     types.INT64,
		"int4":    types.INT64,
		"serial":  types.INT64,
		"serial4": types.INT64,

		"real":   types.FLOAT64,
		"float4": types.FLOAT64,

		"smallint":    types.INT64,
		"int2":        types.INT64,
		"smallserial": types.INT64,
		"serial2":     types.INT64,

		"date":                        types.TIME,
		"time":                        types.TIME,
		"time without time zone":      types.TIME,
		"time with time zone":         types.TIME,
		"timetz":                      types.TIME,
		"timestamp":                   types.TIME,
		"timestamp without time zone": types.TIME,
		"timestamp with time zone":    types.TIME,
		"timestamptz":                 types.TIME,
	}
)
