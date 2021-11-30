package informer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type Cache struct {
	Entities  *EntityCache
	Features  *FeatureCache
	Groups    *GroupCache
	Revisions *RevisionCache
}

func NewCache(
	entities types.EntityList,
	features types.FeatureList,
	groups types.GroupList,
	revisions types.RevisionList) *Cache {
	return &Cache{
		Entities:  &EntityCache{entities},
		Features:  &FeatureCache{features},
		Groups:    &GroupCache{groups},
		Revisions: &RevisionCache{revisions},
	}
}

func (c *Cache) enrich() {
	c.Groups.Enrich(c.Entities)
	c.Features.Enrich(c.Groups)
	// TODO: caching revision data is not necessary, but currently we do it for simplicity
	c.Revisions.Enrich(c.Groups)
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

// Get
func (f *Informer) CacheGetFeature(ctx context.Context, id int) (*types.Feature, error) {
	if feature := f.Cache().Features.Find(func(f *types.Feature) bool {
		return f.ID == id
	}); feature == nil {
		return nil, errdefs.NotFound(fmt.Errorf("feature %d not found", id))
	} else {
		return feature.Copy(), nil
	}
}

func (f *Informer) CacheGetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	if feature := f.Cache().Features.Find(func(f *types.Feature) bool {
		return f.Name == name
	}); feature == nil {
		return nil, errdefs.NotFound(fmt.Errorf("feature '%s' not found", name))
	} else {
		return feature.Copy(), nil
	}
}

// List
func (f *Informer) CacheListFeature(ctx context.Context, opt metadata.ListFeatureOpt) types.FeatureList {
	return f.Cache().Features.List(opt).Copy()
}
