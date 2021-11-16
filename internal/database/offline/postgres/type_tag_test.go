package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestValueTypeTag(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{
			"bit varying",
			types.BYTES,
		},
		{
			"int",
			types.INT64,
		},
		{
			"varchar(32)",
			types.STRING,
		},
		{
			"timestamp without time zone",
			types.TIME,
		},
	} {
		actual, err := postgres.TypeTag(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}
