package sqlutil

import (
	"fmt"
	"strings"
)

func TypeTag(typeMap map[string]string, dbType string) (string, error) {
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
