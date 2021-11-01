package redis

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/require"
)

func TestGetExisted(t *testing.T) {
	ctx, store := prepare()
	importOpt := importSample(t)

	rs, err := store.Get(ctx, online.GetOpt{
		EntityName:  "age",
		RevisionId:  3,
		EntityKey:   "3215",
		FeatureList: importOpt.Features,
	})
	require.NoError(t, err)

	require.Equal(t, int16(18), rs["age"])
	require.Equal(t, "F", rs["gender"])
}

func TestGetNotExistedEntityKey(t *testing.T) {
	ctx, store := prepare()
	importOpt := importSample(t)

	rs, err := store.Get(ctx, online.GetOpt{
		EntityName:  "age",
		RevisionId:  3,
		EntityKey:   "not-existed-key",
		FeatureList: importOpt.Features,
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(rs), "actual: %+v", rs)
}

func TestMultiGet(t *testing.T) {
	ctx, store := prepare()
	importOpt := importSample(t)

	rs, err := store.MultiGet(ctx, online.MultiGetOpt{
		EntityName:  "age",
		RevisionId:  3,
		EntityKeys:  []string{"3215", "3216", "3217", "not-existed-key"},
		FeatureList: importOpt.Features,
	})
	require.NoError(t, err)

	key := "3215"
	require.Equal(t, int16(18), rs[key]["age"])
	require.Equal(t, "F", rs[key]["gender"])

	key = "3216"
	require.Equal(t, int16(29), rs[key]["age"])
	require.Equal(t, nil, rs[key]["gender"])

	key = "3217"
	require.Equal(t, int16(44), rs[key]["age"])
	require.Equal(t, "M", rs[key]["gender"])

	require.Equal(t, 3, len(rs), "actual: %+v", rs)
}
