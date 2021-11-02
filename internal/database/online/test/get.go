package test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExisted(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	s := &SampleSmall
	ctx, store := prepareStore()
	importSample(t, ctx, store, s)

	for _, target := range s.Data {
		opt := online.GetOpt{
			EntityName:  s.Entity.Name,
			RevisionId:  s.Revision.ID,
			FeatureList: s.Features,
			EntityKey:   target.EntityKey(),
		}
		rs, err := store.Get(ctx, opt)
		require.NoError(t, err)

		for i, f := range s.Features {
			assert.Equal(t, rs[f.Name], target.ValueAt(i), "result: %+v", rs)
		}
	}
}

func TestGetNotExistedEntityKey(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	s := &SampleSmall
	ctx, store := prepareStore()
	importSample(t, ctx, store, s)

	rs, err := store.Get(ctx, online.GetOpt{
		EntityName:  s.Entity.Name,
		RevisionId:  s.Revision.ID,
		FeatureList: s.Features,
		EntityKey:   "not-existed-key",
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(rs), "actual: %+v", rs)
}

func TestMultiGet(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	s := &SampleSmall
	ctx, store := prepareStore()
	importSample(t, ctx, store, s)

	keys := []string{"not-existed-key"}
	for _, r := range s.Data {
		keys = append(keys, r.EntityKey())
	}
	rs, err := store.MultiGet(ctx, online.MultiGetOpt{
		EntityName:  s.Entity.Name,
		RevisionId:  s.Revision.ID,
		FeatureList: s.Features,
		EntityKeys:  keys,
	})
	require.NoError(t, err)

	for _, record := range s.Data {
		for i, feature := range s.Features {
			assert.Equal(t, record.ValueAt(i), rs[record.EntityKey()][feature.Name])
		}
	}

	require.Equal(t, len(s.Data), len(rs), "actual: %+v", rs)
}
