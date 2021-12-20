package mysql

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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
