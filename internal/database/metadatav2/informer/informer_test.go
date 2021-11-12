package informer_test

import (
	"context"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/require"
)

func sampleCache() *informer.Cache {
	entities := typesv2.EntityList{
		&typesv2.Entity{
			ID:     1,
			Length: 10,
			Name:   "user",
		},
	}
	return informer.NewCache(entities, nil, nil, nil)
}

func TestInformer(t *testing.T) {
	ctx := context.Background()

	f, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sampleCache(), nil
	})
	require.NoError(t, err)

	entity, err := f.GetEntity(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, int16(1), entity.ID)
	require.Equal(t, 10, entity.Length)
	require.Equal(t, "user", entity.Name)
}

func TestInformerDeepCopy(t *testing.T) {
	ctx := context.Background()

	f, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sampleCache(), nil
	})
	require.NoError(t, err)

	entity, err := f.GetEntity(ctx, 1)
	require.NoError(t, err)

	// changing this entity should not change the internal state of the informer
	entity.Length = 20

	entity, err = f.GetEntity(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, int16(1), entity.ID)
	require.Equal(t, 10, entity.Length)
	require.Equal(t, "user", entity.Name)
}
