package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	_, store := prepareStore()
	assert.NotNil(t, store)
}
