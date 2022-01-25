package informer

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var _ metadata.CacheStore = &Informer{}

type Cache struct {
	Entities *EntityCache
	Features *FeatureCache
	Groups   *GroupCache
}

func NewCache(
	entities types.EntityList,
	features types.FeatureList,
	groups types.GroupList) *Cache {
	return &Cache{
		Entities: &EntityCache{entities},
		Features: &FeatureCache{features},
		Groups:   &GroupCache{groups},
	}
}

func (c *Cache) enrich() {
	c.Groups.Enrich(c.Entities)
	c.Features.Enrich(c.Groups)
}

type Informer struct {
	cache  atomic.Value
	lister func() (*Cache, error)
	ticker *time.Ticker
	quit   chan bool
}

func New(listInterval time.Duration, lister func() (*Cache, error)) (*Informer, error) {
	informer := &Informer{
		cache:  atomic.Value{},
		lister: lister,
		ticker: time.NewTicker(listInterval),
		quit:   make(chan bool),
	}
	if err := informer.Refresh(); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-informer.quit:
				return
			case <-informer.ticker.C:
				if err := informer.Refresh(); err != nil {
					log.Printf("failed refreshing informer: %v\n", err)
				}
				informer.ticker.Reset(listInterval)
			}
		}
	}()
	return informer, nil
}

func (f *Informer) Close() error {
	f.ticker.Stop()
	f.quit <- true
	return nil
}

func (f *Informer) Refresh() error {
	cache, err := f.lister()
	if err != nil {
		return err
	}
	cache.enrich()
	f.cache.Store(cache)
	return nil
}

func (f *Informer) Cache() *Cache {
	return f.cache.Load().(*Cache)
}

func (f *Informer) ListCachedFeature(ctx context.Context, opt metadata.ListCachedFeatureOpt) types.FeatureList {
	return f.Cache().Features.List(opt).Copy()
}
