package sqlite

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
		"integer":   types.INT64,
		"float":     types.FLOAT64,
		"blob":      types.BYTES,
		"text":      types.STRING,
		"timestamp": types.TIME,
		"datetime":  types.TIME,
	}
)
