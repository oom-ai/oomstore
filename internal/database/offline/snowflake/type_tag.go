package snowflake

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TypeTag(dbType string) (string, error) {
	var s = dbType
	if pos := strings.Index(dbType, "("); pos != -1 {
		fmt.Println(s[:pos])
		s = s[:pos]
	}
	s = strings.TrimSpace(strings.ToLower(s))
	if t, ok := typeMap[s]; !ok {
		return "", fmt.Errorf("unsupported sql type: %s", dbType)
	} else {
		return t, nil
	}
}

// TODO: add type NUMBER, DECIMAL, NUMERIC
var (
	typeMap = map[string]string{
		"boolean": types.BOOL,

		"binary":    types.BYTES,
		"varbinary": types.BYTES,

		"integer":  types.INT64,
		"int":      types.INT64,
		"smallint": types.INT64,
		"bigint":   types.INT64,
		"tinyint":  types.INT64,
		"byteint":  types.INT64,

		"double":           types.FLOAT64,
		"double precision": types.FLOAT64,
		"real":             types.FLOAT64,
		"float":            types.FLOAT64,
		"float4":           types.FLOAT64,
		"float8":           types.FLOAT64,

		"string":    types.STRING,
		"text":      types.STRING,
		"varchar":   types.STRING,
		"char":      types.STRING,
		"character": types.STRING,

		"date":      types.TIME,
		"time":      types.TIME,
		"datetime":  types.TIME,
		"timestamp": types.TIME,
	}
)
