package informer

import "github.com/oom-ai/oomstore/pkg/oomstore/typesv2"

type EntityCache struct {
	typesv2.EntityList
}

func (c *EntityCache) List() typesv2.EntityList {
	return c.EntityList
}
