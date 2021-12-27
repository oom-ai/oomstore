package tikv_test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/online/tikv"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_tikv"
)

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx := context.Background()
	store, err := tikv.Open(runtime_tikv.GetOpt())
	if err != nil {
		t.Fatal(err)
	}
	return ctx, store
}

func destroyStore(t *testing.T) func() {
	return func() {
		store, err := tikv.Open(runtime_tikv.GetOpt())
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()

		if err = store.DeleteRange(context.Background(), []byte(""), []byte{}); err != nil {
			t.Fatal(err)
		}
	}

}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, destroyStore(t))
}
