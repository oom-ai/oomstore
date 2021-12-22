package types

import "strconv"

type ValueType int

const (
	STRING ValueType = iota + 1
	INT64
	FLOAT64
	BOOL
	TIME
	BYTES
)

func (t ValueType) String() string {
	switch t {
	case STRING:
		return "string"
	case INT64:
		return "int64"
	case FLOAT64:
		return "float64"
	case BOOL:
		return "bool"
	case TIME:
		return "time"
	case BYTES:
		return "bytes"
	default:
		return "Unknown(" + strconv.Itoa(int(t)) + ")"
	}
}
