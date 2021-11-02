package test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/stretchr/testify/require"
)

func TestPurgeRemovesSpecifiedRevision(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore()
	importSample(t, ctx, store, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision)
	require.NoError(t, err)

	for _, record := range SampleMedium.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			EntityName:  SampleMedium.Entity.Name,
			RevisionId:  SampleMedium.Revision.ID,
			EntityKey:   record.EntityKey(),
			FeatureList: SampleMedium.Features,
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(rs))
	}
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore()
	importSample(t, ctx, store, &SampleSmall, &SampleMedium)

	err := store.Purge(ctx, SampleMedium.Revision)
	require.NoError(t, err)

	for _, record := range SampleSmall.Data {
		rs, err := store.Get(ctx, online.GetOpt{
			EntityName:  SampleSmall.Entity.Name,
			RevisionId:  SampleSmall.Revision.ID,
			EntityKey:   record.EntityKey(),
			FeatureList: SampleSmall.Features,
		})
		require.NoError(t, err)
		for i, f := range SampleSmall.Features {
			require.Equal(t, record.ValueAt(i), rs[f.Name])
		}
	}
}
