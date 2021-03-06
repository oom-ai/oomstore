package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/require"
)

func TestPurgeRemovesSpecifiedRevision(t *testing.T, prepareStore PrepareStoreFn, destroystore DestroyStoreFn) {
	t.Cleanup(destroystore)

	ctx, store := prepareStore(t)
	defer store.Close()
	importSample(t, ctx, store, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision.ID)
	require.NoError(t, err)

	for _, record := range SampleMedium.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			EntityKey:  record.EntityKey(),
			Group:      *SampleMedium.Revision.Group,
			Features:   SampleMedium.Features,
			RevisionID: &SampleMedium.Revision.ID,
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(rs))
	}
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T, prepareStore PrepareStoreFn, destroystore DestroyStoreFn) {
	t.Cleanup(destroystore)

	ctx, store := prepareStore(t)
	defer store.Close()
	importSample(t, ctx, store, &SampleSmall, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision.ID)
	require.NoError(t, err)

	for _, record := range SampleSmall.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			EntityKey:  record.EntityKey(),
			Group:      *SampleSmall.Revision.Group,
			Features:   SampleSmall.Features,
			RevisionID: &SampleSmall.Revision.ID,
		})
		require.NoError(t, err)
		for i, f := range SampleSmall.Features {
			compareFeatureValue(t, record.ValueAt(i), rs[f.FullName()], f.ValueType)
		}
	}
}
