package sqlutil

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func GetValueType(typeMap map[string]types.ValueType, dbType string) (types.ValueType, error) {
	var s = dbType
	if pos := strings.Index(dbType, "("); pos != -1 {
		s = s[:pos]
	}
	s = strings.TrimSpace(strings.ToLower(s))
	if t, ok := typeMap[s]; !ok {
		return 0, fmt.Errorf("unsupported sql type: %s", dbType)
	} else {
		return t, nil
	}
}
