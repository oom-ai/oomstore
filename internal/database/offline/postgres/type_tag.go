package postgres

import (
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TypeTag(dbType string) (string, error) {
	return sqlutil.TypeTag(typeMap, dbType)
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
