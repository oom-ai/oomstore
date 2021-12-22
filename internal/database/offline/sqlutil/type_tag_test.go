package sqlutil_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var typeMap = map[string]string{
	"boolean":  types.BOOL,
	"binary":   types.BYTES,
	"bigint":   types.INT64,
	"double":   types.FLOAT64,
	"varchar":  types.STRING,
	"datetime": types.TIME,
}

func TestTypeTag(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{"boolean", types.BOOL},
		{"binary", types.BYTES},
		{"bigint", types.INT64},
		{"double", types.FLOAT64},
		{"varchar(32)", types.STRING},
		{"VARCHAR(64)", types.STRING},
		{"datetime", types.TIME},
	} {
		actual, err := sqlutil.TypeTag(typeMap, tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
