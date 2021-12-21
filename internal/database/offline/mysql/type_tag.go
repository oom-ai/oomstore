package mysql

import (
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TypeTag(dbType string) (string, error) {
	return sqlutil.TypeTag(typeMap, dbType)
}

var (
	typeMap = map[string]string{
		"boolean": types.BOOL,
		"bool":    types.BOOL,

		"binary":    types.BYTES,
		"varbinary": types.BYTES,

		"integer":   types.INT64,
		"int":       types.INT64,
		"smallint":  types.INT64,
		"bigint":    types.INT64,
		"tinyint":   types.INT64,
		"mediumint": types.INT64,

		"double": types.FLOAT64,
		"float":  types.FLOAT64,

		"text":    types.STRING,
		"varchar": types.STRING,
		"char":    types.STRING,

		"date":      types.TIME,
		"time":      types.TIME,
		"datetime":  types.TIME,
		"timestamp": types.TIME,
		"year":      types.TIME,
	}
)
