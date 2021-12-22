package types

import (
	"fmt"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

type ValueType int

const (
	INVALID ValueType = iota
	STRING
	INT64
	FLOAT64
	BOOL
	TIME
	BYTES
)

var allValueTypes = map[ValueType]string{
	STRING:  "string",
	INT64:   "int64",
	FLOAT64: "float64",
	BOOL:    "bool",
	TIME:    "time",
	BYTES:   "bytes",
}

func (t ValueType) String() string {
	if s, ok := allValueTypes[t]; ok {
		return s
	} else {
		return "Unknown(" + strconv.Itoa(int(t)) + ")"
	}
}

func (v ValueType) Validate() error {
	if _, ok := allValueTypes[v]; ok {
		return nil
	} else {
		return errdefs.InvalidAttribute(fmt.Errorf("Invalid value type: %d", v))
	}
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
	return INVALID, errdefs.InvalidAttribute(fmt.Errorf("Unknown value type: %s", s))
}
