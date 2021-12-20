package snowflake_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/snowflake"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestValueTypeTag(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{
			"binary",
			types.BYTES,
		},
		{
			"SMALLINT",
			types.INT64,
		},
		{
			"CHARACTER",
			types.STRING,
		},
		{
			"datetime",
			types.TIME,
		},
	} {
		actual, err := snowflake.TypeTag(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
