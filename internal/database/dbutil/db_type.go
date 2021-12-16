package dbutil

import (
	"fmt"

	"github.com/ethhte88/oomstore/pkg/errdefs"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func GetDbTypeFrom(backend types.BackendType, valueType string) (string, error) {
	var mp map[string]string
	switch backend {
	case types.POSTGRES, types.REDIS, types.MYSQL, types.SNOWFLAKE, types.DYNAMODB:
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported backend: %s", backend))
	case types.CASSANDRA:
		mp = cassandraType
	}

	t, ok := mp[valueType]
	if !ok {
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported value type: %s", valueType))
	}
	return t, nil
}

var (
	cassandraType = map[string]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "text",
	}
)
