package sqlite_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/sqlite"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestValueTypeTag(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{
			"blob",
			types.BYTES,
		},
		{
			"integer",
			types.INT64,
		},
		{
			"float",
			types.FLOAT64,
		},
		{
			"text",
			types.STRING,
		},
		{
			"timestamp",
			types.TIME,
		},
		{
			"datetime",
			types.TIME,
		},
	} {
		actual, err := sqlite.TypeTag(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
