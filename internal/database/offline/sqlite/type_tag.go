package sqlite

import (
	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TypeTag(dbType string) (string, error) {
	return sqlutil.TypeTag(typeMap, dbType)
}

var (
	typeMap = map[string]string{
		"integer":   types.INT64,
		"float":     types.FLOAT64,
		"blob":      types.BYTES,
		"text":      types.STRING,
		"timestamp": types.TIME,
		"datetime":  types.TIME,
	}
)
