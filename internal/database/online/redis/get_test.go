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
