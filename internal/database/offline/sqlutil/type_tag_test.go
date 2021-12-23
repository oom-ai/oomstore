package sqlutil_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var typeMap = map[string]types.ValueType{
	"boolean":  types.Bool,
	"binary":   types.Bytes,
	"bigint":   types.Int64,
	"double":   types.Float64,
	"varchar":  types.String,
	"datetime": types.Time,
}

func TestValueType(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected types.ValueType
	}{
		{"boolean", types.Bool},
		{"binary", types.Bytes},
		{"bigint", types.Int64},
		{"double", types.Float64},
		{"varchar(32)", types.String},
		{"VARCHAR(64)", types.String},
		{"datetime", types.Time},
	} {
		actual, err := sqlutil.GetValueType(typeMap, tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
