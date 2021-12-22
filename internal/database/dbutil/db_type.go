package dbutil

import (
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DBValueType(backend types.BackendType, valueType types.ValueType) (string, error) {
	var mp map[types.ValueType]string
	switch backend {
	case types.POSTGRES:
		mp = valueTypeToPostgresType
	case types.SQLite:
		mp = valueTypeToSQLiteType
	case types.MYSQL:
		mp = valueTypeToMySQLType
	case types.CASSANDRA:
		mp = valueTypeToCassandraType
	case types.SNOWFLAKE:
		mp = valueTypeToSnowFlake
	case types.DYNAMODB:
		mp = valueTypeToDynamoDB
	case types.REDSHIFT:
		mp = valueTypeToRedshiftType
	default:
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported backend: %s", backend))
	}

	t, ok := mp[valueType]
	if !ok {
		return "", errdefs.InvalidAttribute(fmt.Errorf("unsupported value type: %s", valueType))
	}
	return t, nil
}

func ValueType(backend types.BackendType, dbValueType string) (types.ValueType, error) {
	var mp map[string]types.ValueType
	switch backend {
	case types.POSTGRES:
		mp = postgresTypeToValueType
	case types.MYSQL:
		mp = mySQLTypeToValueType
	case types.SNOWFLAKE:
		mp = snowflakeTypeToValueType
	case types.BIGQUERY:
		mp = bigQueryTypeToValueType
	case types.SQLite:
		mp = sqliteTypeToValueType
	case types.REDSHIFT:
		mp = redshiftTypeToValueType
	default:
		return 0, errdefs.InvalidAttribute(fmt.Errorf("unsupported backend: %s", backend))
	}

	t, ok := mp[strings.ToLower(dbValueType)]
	if !ok {
		return 0, errdefs.InvalidAttribute(fmt.Errorf("unsupported db value type: %s", dbValueType))
	}
	return t, nil
}

// Mapping feature value type to database data type
var (
	valueTypeToSQLiteType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "integer",
		types.FLOAT64: "float",
		types.BOOL:    "integer",
		types.BYTES:   "blob",
		types.TIME:    "timestamp",
	}
	valueTypeToMySQLType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "bool",
		types.TIME:    "datetime",
		types.BYTES:   "varbinary",
	}
	valueTypeToPostgresType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double precision",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "bytea",
	}
	valueTypeToSnowFlake = map[types.ValueType]string{
		types.STRING:  "varchar",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "varbinary",
	}
	valueTypeToDynamoDB = map[types.ValueType]string{
		types.STRING:  "String",
		types.INT64:   "BigInteger",
		types.FLOAT64: "Double",
		types.BOOL:    "Boolean",
		types.TIME:    "Date",
		types.BYTES:   "ByteBuffer",
	}
	valueTypeToCassandraType = map[types.ValueType]string{
		types.STRING:  "text",
		types.INT64:   "bigint",
		types.FLOAT64: "double",
		types.BOOL:    "boolean",
		types.TIME:    "timestamp",
		types.BYTES:   "blob",
	}
	valueTypeToRedshiftType = valueTypeToPostgresType
)

// Mapping database data type to feature value type
var (
	postgresTypeToValueType = map[string]types.ValueType{
		"bigint":    types.INT64,
		"int8":      types.INT64,
		"bigserial": types.INT64,
		"serial8":   types.INT64,

		"boolean": types.BOOL,
		"bool":    types.BOOL,

		"bytea":       types.BYTES,
		"jsonb":       types.BYTES,
		"uuid":        types.BYTES,
		"bit":         types.BYTES,
		"bit varying": types.BYTES,
		"character":   types.BYTES,
		"char":        types.BYTES,
		"json":        types.BYTES,
		"money":       types.BYTES,
		"numeric":     types.BYTES,

		"character varying": types.STRING,
		"text":              types.STRING,
		"varchar":           types.STRING,

		"double precision": types.FLOAT64,
		"float8":           types.FLOAT64,

		"integer": types.INT64,
		"int":     types.INT64,
		"int4":    types.INT64,
		"serial":  types.INT64,
		"serial4": types.INT64,

		"real":   types.FLOAT64,
		"float4": types.FLOAT64,

		"smallint":    types.INT64,
		"int2":        types.INT64,
		"smallserial": types.INT64,
		"serial2":     types.INT64,

		"date":                        types.TIME,
		"time":                        types.TIME,
		"time without time zone":      types.TIME,
		"time with time zone":         types.TIME,
		"timetz":                      types.TIME,
		"timestamp":                   types.TIME,
		"timestamp without time zone": types.TIME,
		"timestamp with time zone":    types.TIME,
		"timestamptz":                 types.TIME,
	}
	bigQueryTypeToValueType = map[string]types.ValueType{
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
	mySQLTypeToValueType = map[string]types.ValueType{
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
	// TODO: add type NUMBER, DECIMAL, NUMERIC
	snowflakeTypeToValueType = map[string]types.ValueType{
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
	sqliteTypeToValueType = map[string]types.ValueType{
		"integer":   types.INT64,
		"float":     types.FLOAT64,
		"blob":      types.BYTES,
		"text":      types.STRING,
		"timestamp": types.TIME,
		"datetime":  types.TIME,
	}
	// Redshift data type: https://docs.aws.amazon.com/redshift/latest/dg/c_Supported_data_types.html
	redshiftTypeToValueType = map[string]types.ValueType{
		"bigint": types.INT64,
		"int8":   types.INT64,

		"boolean": types.BOOL,
		"bool":    types.BOOL,

		"character":      types.BYTES,
		"char":           types.BYTES,
		"nchar":          types.BYTES,
		"bpchar":         types.BYTES,
		"varbyte":        types.BYTES,
		"varbinary":      types.BYTES,
		"binary varying": types.BYTES,
		"numeric":        types.BYTES,
		"decimal":        types.BYTES,

		"character varying": types.STRING,
		"text":              types.STRING,
		"varchar":           types.STRING,
		"nvarchar":          types.STRING,

		"double precision": types.FLOAT64,
		"float8":           types.FLOAT64,
		"float":            types.FLOAT64,

		"integer": types.INT64,
		"int":     types.INT64,
		"int4":    types.INT64,

		"real":   types.FLOAT64,
		"float4": types.FLOAT64,

		"smallint": types.INT64,
		"int2":     types.INT64,

		"date":                        types.TIME,
		"time":                        types.TIME,
		"time without time zone":      types.TIME,
		"time with time zone":         types.TIME,
		"timetz":                      types.TIME,
		"timestamp":                   types.TIME,
		"timestamp without time zone": types.TIME,
		"timestamp with time zone":    types.TIME,
		"timestamptz":                 types.TIME,
	}
)
