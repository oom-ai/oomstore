package types

import (
	"fmt"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/errdefs"
)

type ValueType int

const (
	Invalid ValueType = iota
	String
	Int64
	Float64
	Bool
	Time
	Bytes
)

var allValueTypes = map[ValueType]string{
	String:  "string",
	Int64:   "int64",
	Float64: "float64",
	Bool:    "bool",
	Time:    "time",
	Bytes:   "bytes",
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
		return Bytes, nil
	case "int64":
		return Int64, nil
	case "float64":
		return Float64, nil
	case "bool":
		return Bool, nil
	case "time":
		return Time, nil
	case "bytes":
		return Bytes, nil
	}
	return Invalid, errdefs.InvalidAttribute(fmt.Errorf("Unknown value type: %s", s))
}
