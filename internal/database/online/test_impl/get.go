package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/require"
)

func TestGetExisted(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	s := &SampleSmall
	ctx, store := prepareStore(t)
	defer store.Close()
	importSample(t, ctx, store, s)

	for _, target := range s.Data {
		opt := online.GetOpt{
			Entity:     s.Entity,
			RevisionID: s.Revision.ID,
			Features:   s.Features,
			EntityKey:  target.EntityKey(),
		}
		rs, err := store.Get(ctx, opt)
		require.NoError(t, err)

		for i, f := range s.Features {
			compareFeatureValue(t, target.ValueAt(i), rs[f.FullName], f.ValueType)
		}
	}
}

func TestGetNotExistedEntityKey(t *testing.T, prepareStore PrepareStoreFn, destroystore DestroyStoreFn) {
	t.Cleanup(destroystore)

	s := &SampleSmall
	ctx, store := prepareStore(t)
	defer store.Close()
	importSample(t, ctx, store, s)

	rs, err := store.Get(ctx, online.GetOpt{
		Entity:     s.Entity,
		RevisionID: s.Revision.ID,
		Features:   s.Features,
		EntityKey:  "not-existed-key",
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(rs), "actual: %+v", rs)
}

func TestMultiGet(t *testing.T, prepareStore PrepareStoreFn, destroystore DestroyStoreFn) {
	t.Cleanup(destroystore)

	s := &SampleSmall
	ctx, store := prepareStore(t)
	defer store.Close()
	importSample(t, ctx, store, s)

	keys := []string{"not-existed-key"}
	for _, r := range s.Data {
		keys = append(keys, r.EntityKey())
	}
	rs, err := store.MultiGet(ctx, online.MultiGetOpt{
		Entity:     s.Entity,
		RevisionID: s.Revision.ID,
		Features:   s.Features,
		EntityKeys: keys,
	})
	require.NoError(t, err)
	for _, record := range s.Data {
		for i, feature := range s.Features {
			compareFeatureValue(t, record.ValueAt(i), rs[record.EntityKey()][feature.FullName], feature.ValueType)
		}
	}

	require.Equal(t, len(s.Data), len(rs), "actual: %+v", rs)
}
