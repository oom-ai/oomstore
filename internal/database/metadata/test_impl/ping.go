package test_impl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	require.NoError(t, store.Ping(ctx))
}
