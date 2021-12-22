package dbutil

import (
	"fmt"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DBValueType(backend types.BackendType, valueType types.ValueType) (string, error) {
	var mp map[types.ValueType]string
	switch backend {
	case types.POSTGRES:
		mp = postgresType
	case types.SQLite:
		mp = sqliteType
	case types.MYSQL:
		mp = mysqlType
	case types.CASSANDRA:
		mp = cassandraType
	case types.SNOWFLAKE:
		mp = snowFlake
	case types.DYNAMODB:
		mp = dynamoDB
	case types.REDSHIFT:
		mp = redshiftType
	default:
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported backend: %s", backend))
	}

	t, ok := mp[valueType]
	if !ok {
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported value type: %s", valueType))
	}
	return t, nil
}

var (
	sqliteType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "integer",
		types.FLOAT64: "float",
		types.BOOL:    "integer",
		types.BYTES:   "blob",
		types.TIME:    "timestamp",
	}
	mysqlType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "bool",
		types.TIME:    "datetime",
		types.BYTES:   "varbinary",
	}
	postgresType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double precision",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "bytea",
	}
	snowFlake = map[types.ValueType]string{
		types.STRING:  "varchar",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "varbinary",
	}
	dynamoDB = map[types.ValueType]string{
		types.STRING:  "String",
		types.INT64:   "BigInteger",
		types.FLOAT64: "Double",
		types.BOOL:    "Boolean",
		types.TIME:    "Date",
		types.BYTES:   "ByteBuffer",
	}
	cassandraType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "blob",
	}
	redshiftType = postgresType
)
