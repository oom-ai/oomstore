package types

import (
	"fmt"
	"strconv"
)

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
	}
	return "Unknown(" + strconv.Itoa(int(t)) + ")"
}

func ParseValueType(s string) (ValueType, error) {
	switch s {
	case "string":
		return BYTES, nil
	case "int64":
		return INT64, nil
	case "float64":
		return FLOAT64, nil
	case "bool":
		return BOOL, nil
	case "time":
		return TIME, nil
	case "bytes":
		return BYTES, nil
	}
	return 0, fmt.Errorf("Unknown value type: %s", s)
}
