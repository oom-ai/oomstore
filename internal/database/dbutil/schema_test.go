package dbutil_test

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/stretchr/testify/require"
)

func TestBuildSchema(t *testing.T) {
	cases := []struct {
		description string
		schema      dbutil.Schema
		schemaType  dbutil.SchemaType

		want    string
		wantErr error
	}{
		{
			description: "cassandra schema",
			schemaType:  dbutil.Cassandra,
			schema: dbutil.Schema{
				TableName:  "user",
				EntityName: "user_id",
				Columns: []dbutil.Column{
					{
						Name:   "age",
						DbType: "int",
					},
					{
						Name:   "gender",
						DbType: "varchar",
					},
					{
						Name:   "phone",
						DbType: "182",
					},
				},
			},

			wantErr: nil,
			want: `
CREATE TABLE user (
    user_id  TEXT PRIMARY KEY,
    age  int,
    gender  varchar,
    phone  182
)
`,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			schema, err := dbutil.BuildSchema(c.schema, c.schemaType)
			require.Equal(t, c.want, schema)
			require.Equal(t, c.wantErr, err)
		})
	}
}
