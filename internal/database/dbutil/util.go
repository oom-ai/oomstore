package dbutil

import (
	"fmt"
	"time"
)

type RowMap = map[string]interface{}

func TempTable(prefix string) string {
	return fmt.Sprintf("tmp_%s_%d", prefix, time.Now().UnixNano())
}
