package bigquery

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
