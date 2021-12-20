package mysql_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/mysql"
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
			"int",
			types.INT64,
		},
		{
			"char(32)",
			types.STRING,
		},
		{
			"year",
			types.TIME,
		},
	} {
		actual, err := mysql.TypeTag(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
