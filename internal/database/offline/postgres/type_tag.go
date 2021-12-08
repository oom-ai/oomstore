package postgres

import (
	"fmt"
	"strings"

	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

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

		"bytea":       types.BYTES,
		"jsonb":       types.BYTES,
		"uuid":        types.BYTES,
		"bit":         types.BYTES,
		"bit varying": types.BYTES,
		"character":   types.BYTES,
		"char":        types.BYTES,
		"json":        types.BYTES,
		"money":       types.BYTES,
		"numeric":     types.BYTES,

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
