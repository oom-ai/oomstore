package dbutil

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func DBValueType(backend types.BackendType, valueType types.ValueType) (string, error) {
	var mp map[types.ValueType]string
	switch backend {
	case types.BackendPostgres:
		mp = valueTypeToPostgresType
	case types.BackendSQLite:
		mp = valueTypeToSQLiteType
	case types.BackendMySQL:
		mp = valueTypeToMySQLType
	case types.BackendCassandra:
		mp = valueTypeToCassandraType
	case types.BackendSnowflake:
		mp = valueTypeToSnowFlakeType
	case types.BackendDynamoDB:
		mp = valueTypeToDynamoDBType
	case types.BackendRedshift:
		mp = valueTypeToRedshiftType
	case types.BackendBigQuery:
		mp = valueTypeToBigQueryType
	default:
		return "", errdefs.InvalidAttribute(errors.Errorf("unsupported backend: %s", backend))
	}

	t, ok := mp[valueType]
	if !ok {
		return "", errdefs.InvalidAttribute(errors.Errorf("unsupported value type: %s", valueType))
	}
	return t, nil
}

// Used for inferring feature value type from a supported offline db value type
func ValueType(backend types.BackendType, dbValueType string) (types.ValueType, error) {
	var mp map[string]types.ValueType
	switch backend {
	case types.BackendPostgres:
		mp = postgresTypeToValueType
	case types.BackendMySQL:
		mp = mysqlTypeToValueType
	case types.BackendSnowflake:
		mp = snowflakeTypeToValueType
	case types.BackendBigQuery:
		mp = bigQueryTypeToValueType
	case types.BackendSQLite:
		mp = sqliteTypeToValueType
	case types.BackendRedshift:
		mp = redshiftTypeToValueType
	default:
		return 0, errdefs.InvalidAttribute(errors.Errorf("unsupported backend: %s", backend))
	}
	dbValueType = strings.ToLower(dbValueType)
	i := strings.Index(dbValueType, "(")
	if i != -1 {
		dbValueType = dbValueType[0:i]
	}
	t, ok := mp[dbValueType]
	if !ok {
		return 0, errdefs.InvalidAttribute(errors.Errorf("unsupported db value type: %s", dbValueType))
	}

	return t, nil
}

// Mapping feature value type to database data type
var (
	valueTypeToSQLiteType = map[types.ValueType]string{
		types.String:  "text",
		types.Int64:   "integer",
		types.Float64: "float",
		types.Bool:    "integer",
		types.Bytes:   "blob",
		types.Time:    "timestamp",
	}
	valueTypeToMySQLType = map[types.ValueType]string{
		types.String:  "text",
		types.Int64:   "bigint",
		types.Float64: "double",
		types.Bool:    "bool",
		types.Time:    "datetime",
		types.Bytes:   "varbinary",
	}
	valueTypeToPostgresType = map[types.ValueType]string{
		types.String:  "text",
		types.Int64:   "bigint",
		types.Float64: "double precision",
		types.Bool:    "boolean",
		types.Time:    "timestamp",
		types.Bytes:   "bytea",
	}
	valueTypeToSnowFlakeType = map[types.ValueType]string{
		types.String:  "varchar",
		types.Int64:   "bigint",
		types.Float64: "double",
		types.Bool:    "boolean",
		types.Time:    "timestamp",
		types.Bytes:   "varbinary",
	}
	valueTypeToDynamoDBType = map[types.ValueType]string{
		types.String:  "String",
		types.Int64:   "BigInteger",
		types.Float64: "Double",
		types.Bool:    "Boolean",
		types.Time:    "Date",
		types.Bytes:   "ByteBuffer",
	}
	valueTypeToCassandraType = map[types.ValueType]string{
		types.String:  "text",
		types.Int64:   "bigint",
		types.Float64: "double",
		types.Bool:    "boolean",
		types.Time:    "timestamp",
		types.Bytes:   "blob",
	}
	valueTypeToBigQueryType = map[types.ValueType]string{
		types.String:  "string",
		types.Int64:   "bigint",
		types.Float64: "float64",
		types.Bool:    "bool",
		types.Time:    "datetime",
		types.Bytes:   "bytes",
	}
	valueTypeToRedshiftType = valueTypeToPostgresType
)

// Mapping database data type to feature value type
var (
	postgresTypeToValueType = map[string]types.ValueType{
		"bigint":    types.Int64,
		"int8":      types.Int64,
		"bigserial": types.Int64,
		"serial8":   types.Int64,

		"boolean": types.Bool,
		"bool":    types.Bool,

		"bytea":       types.Bytes,
		"jsonb":       types.Bytes,
		"uuid":        types.Bytes,
		"bit":         types.Bytes,
		"bit varying": types.Bytes,
		"character":   types.Bytes,
		"char":        types.Bytes,
		"json":        types.Bytes,
		"money":       types.Bytes,
		"numeric":     types.Bytes,

		"character varying": types.String,
		"text":              types.String,
		"varchar":           types.String,

		"double precision": types.Float64,
		"float8":           types.Float64,

		"integer": types.Int64,
		"int":     types.Int64,
		"int4":    types.Int64,
		"serial":  types.Int64,
		"serial4": types.Int64,

		"real":   types.Float64,
		"float4": types.Float64,

		"smallint":    types.Int64,
		"int2":        types.Int64,
		"smallserial": types.Int64,
		"serial2":     types.Int64,

		"date":                        types.Time,
		"time":                        types.Time,
		"time without time zone":      types.Time,
		"time with time zone":         types.Time,
		"timetz":                      types.Time,
		"timestamp":                   types.Time,
		"timestamp without time zone": types.Time,
		"timestamp with time zone":    types.Time,
		"timestamptz":                 types.Time,
	}
	bigQueryTypeToValueType = map[string]types.ValueType{
		"bool":     types.Bool,
		"bytes":    types.Bytes,
		"datetime": types.Time,
		"string":   types.String,

		"bigint":   types.Int64,
		"smallint": types.Int64,
		"int64":    types.Int64,
		"integer":  types.Int64,
		"int":      types.Int64,

		"float64": types.Float64,
		"numeric": types.Float64,
		"decimal": types.Float64,
	}
	mysqlTypeToValueType = map[string]types.ValueType{
		"boolean": types.Bool,
		"bool":    types.Bool,

		"binary":    types.Bytes,
		"varbinary": types.Bytes,

		"integer":   types.Int64,
		"int":       types.Int64,
		"smallint":  types.Int64,
		"bigint":    types.Int64,
		"tinyint":   types.Int64,
		"mediumint": types.Int64,

		"double": types.Float64,
		"float":  types.Float64,

		"text":    types.String,
		"varchar": types.String,
		"char":    types.String,

		"date":      types.Time,
		"time":      types.Time,
		"datetime":  types.Time,
		"timestamp": types.Time,
		"year":      types.Time,
	}
	// TODO: add type NUMBER, DECIMAL, NUMERIC
	snowflakeTypeToValueType = map[string]types.ValueType{
		"boolean": types.Bool,

		"binary":    types.Bytes,
		"varbinary": types.Bytes,

		"integer":  types.Int64,
		"int":      types.Int64,
		"smallint": types.Int64,
		"bigint":   types.Int64,
		"tinyint":  types.Int64,
		"byteint":  types.Int64,
		"number":   types.Int64,

		"double":           types.Float64,
		"double precision": types.Float64,
		"real":             types.Float64,
		"float":            types.Float64,
		"float4":           types.Float64,
		"float8":           types.Float64,

		"string":    types.String,
		"text":      types.String,
		"varchar":   types.String,
		"char":      types.String,
		"character": types.String,

		"date":      types.Time,
		"time":      types.Time,
		"datetime":  types.Time,
		"timestamp": types.Time,
	}
	sqliteTypeToValueType = map[string]types.ValueType{
		"integer":   types.Int64,
		"int":       types.Int64,
		"tinyint":   types.Int64,
		"smallint":  types.Int64,
		"bigint":    types.Int64,
		"mediumint": types.Int64,

		"float":            types.Float64,
		"real":             types.Float64,
		"double":           types.Float64,
		"double precision": types.Float64,
		"numeric":          types.Float64,
		"decimal":          types.Float64,

		"blob": types.Bytes,

		"text":              types.String,
		"clob":              types.String,
		"character varying": types.String,
		"varchar":           types.String,
		"nvarchar":          types.String,
		"nchar":             types.String,
		"character":         types.String,
		"native character":  types.String,

		"boolean":  types.Bool,
		"datetime": types.Time,
	}
	// Redshift data type: https://docs.aws.amazon.com/redshift/latest/dg/c_Supported_data_types.html
	redshiftTypeToValueType = map[string]types.ValueType{
		"bigint": types.Int64,
		"int8":   types.Int64,

		"boolean": types.Bool,
		"bool":    types.Bool,

		"character":      types.Bytes,
		"char":           types.Bytes,
		"nchar":          types.Bytes,
		"bpchar":         types.Bytes,
		"varbyte":        types.Bytes,
		"varbinary":      types.Bytes,
		"binary varying": types.Bytes,
		"numeric":        types.Bytes,
		"decimal":        types.Bytes,

		"character varying": types.String,
		"text":              types.String,
		"varchar":           types.String,
		"nvarchar":          types.String,

		"double precision": types.Float64,
		"float8":           types.Float64,
		"float":            types.Float64,

		"integer": types.Int64,
		"int":     types.Int64,
		"int4":    types.Int64,

		"real":   types.Float64,
		"float4": types.Float64,

		"smallint": types.Int64,
		"int2":     types.Int64,

		"date":                        types.Time,
		"time":                        types.Time,
		"time without time zone":      types.Time,
		"time with time zone":         types.Time,
		"timetz":                      types.Time,
		"timestamp":                   types.Time,
		"timestamp without time zone": types.Time,
		"timestamp with time zone":    types.Time,
		"timestamptz":                 types.Time,
	}
)
