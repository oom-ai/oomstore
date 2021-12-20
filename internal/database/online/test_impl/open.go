package test_impl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T, prepareStore PrepareStoreFn) {
	_, store := prepareStore(t)
	assert.NotNil(t, store)
	store.Close()
}
