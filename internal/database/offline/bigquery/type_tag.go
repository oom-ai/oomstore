package bigquery

import (
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TypeTag(dbType string) (string, error) {
	return sqlutil.TypeTag(typeMap, dbType)
}

var (
	typeMap = map[string]string{
		"bool":     types.BOOL,
		"bytes":    types.BYTES,
		"datetime": types.TIME,
		"string":   types.STRING,

		"bigint":   types.INT64,
		"smallint": types.INT64,
		"int64":    types.INT64,
		"integer":  types.INT64,
		"int":      types.INT64,

		"float64": types.FLOAT64,
		"numeric": types.FLOAT64,
		"decimal": types.FLOAT64,
	}
)
