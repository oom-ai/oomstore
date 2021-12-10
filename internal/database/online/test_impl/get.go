package test_impl

import (
	"testing"
	"time"

	"github.com/ethhte88/oomstore/internal/database/online"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExisted(t *testing.T, prepareStore PrepareStoreFn) {
	s := &SampleSmall
	ctx, store := prepareStore()
	defer store.Close()
	importSample(t, ctx, store, s)

	for _, target := range s.Data {
		opt := online.GetOpt{
			Entity:      s.Entity,
			RevisionID:  s.Revision.ID,
			FeatureList: s.Features,
			EntityKey:   target.EntityKey(),
		}
		rs, err := store.Get(ctx, opt)
		require.NoError(t, err)

		for i, f := range s.Features {
			if f.ValueType == types.TIME {
				expected, ok := target.ValueAt(i).(time.Time)
				require.Equal(t, true, ok)
				actual, ok := rs[f.Name].(time.Time)
				require.Equal(t, true, ok)
				assert.Equal(t, expected.UnixMilli(), actual.UnixMilli())
			} else {
				assert.Equal(t, target.ValueAt(i), rs[f.Name], "result: %+v", rs)
			}
		}
	}
}

func TestGetNotExistedEntityKey(t *testing.T, prepareStore PrepareStoreFn) {
	s := &SampleSmall
	ctx, store := prepareStore()
	defer store.Close()
	importSample(t, ctx, store, s)

	rs, err := store.Get(ctx, online.GetOpt{
		Entity:      s.Entity,
		RevisionID:  s.Revision.ID,
		FeatureList: s.Features,
		EntityKey:   "not-existed-key",
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(rs), "actual: %+v", rs)
}

func TestMultiGet(t *testing.T, prepareStore PrepareStoreFn) {
	s := &SampleSmall
	ctx, store := prepareStore()
	defer store.Close()
	importSample(t, ctx, store, s)

	keys := []string{"not-existed-key"}
	for _, r := range s.Data {
		keys = append(keys, r.EntityKey())
	}
	rs, err := store.MultiGet(ctx, online.MultiGetOpt{
		Entity:      s.Entity,
		RevisionID:  s.Revision.ID,
		FeatureList: s.Features,
		EntityKeys:  keys,
	})
	require.NoError(t, err)

	for _, record := range s.Data {
		for i, feature := range s.Features {
			if feature.ValueType == types.TIME {
				expected, ok := record.ValueAt(i).(time.Time)
				require.Equal(t, true, ok)
				actual, ok := rs[record.EntityKey()][feature.Name].(time.Time)
				require.Equal(t, true, ok)
				assert.Equal(t, expected.UnixMilli(), actual.UnixMilli())
			} else {
				assert.Equal(t, record.ValueAt(i), rs[record.EntityKey()][feature.Name], "result: %+v", rs)
			}
		}
	}

	require.Equal(t, len(s.Data), len(rs), "actual: %+v", rs)
}
