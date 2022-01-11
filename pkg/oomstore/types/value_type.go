package types

import (
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

var allValueTypeStrings = map[string]ValueType{
	"string":  String,
	"int64":   Int64,
	"float64": Float64,
	"bool":    Bool,
	"time":    Time,
	"bytes":   Bytes,
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
		return errdefs.InvalidAttribute(errdefs.Errorf("Invalid value type: %d", v))
	}
}

func ParseValueType(s string) (ValueType, error) {
	if v, ok := allValueTypeStrings[s]; ok {
		return v, nil
	} else {
		return Invalid, errdefs.InvalidAttribute(errdefs.Errorf("Unknown value type: %s", s))
	}
}
