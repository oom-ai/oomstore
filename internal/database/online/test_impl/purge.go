package test_impl

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/stretchr/testify/require"
)

func TestPurgeRemovesSpecifiedRevision(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()
	importSample(t, ctx, store, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision.ID)
	require.NoError(t, err)

	for _, record := range SampleMedium.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			Entity:      SampleMedium.Entity,
			RevisionID:  SampleMedium.Revision.ID,
			EntityKey:   record.EntityKey(),
			FeatureList: SampleMedium.Features,
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(rs))
	}
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()
	importSample(t, ctx, store, &SampleSmall, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision.ID)
	require.NoError(t, err)

	for _, record := range SampleSmall.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			Entity:      SampleSmall.Entity,
			RevisionID:  SampleSmall.Revision.ID,
			EntityKey:   record.EntityKey(),
			FeatureList: SampleSmall.Features,
		})
		require.NoError(t, err)
		for i, f := range SampleSmall.Features {
			compareFeatureValue(t, record.ValueAt(i), rs[f.Name], f.ValueType)
		}
	}
}
