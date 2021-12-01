package informer

import (
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type EntityCache struct {
	types.EntityList
}

func (c *EntityCache) List() types.EntityList {
	return c.EntityList
}
