package informer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

type Cache struct {
	Entities  *EntityCache
	Features  *FeatureCache
	Groups    *GroupCache
	Revisions *RevisionCache
}

func NewCache(
	entities typesv2.EntityList,
	features typesv2.FeatureList,
	groups typesv2.FeatureGroupList,
	revisions typesv2.RevisionList) *Cache {
	return &Cache{
		Entities:  &EntityCache{entities},
		Features:  &FeatureCache{features},
		Groups:    &GroupCache{groups},
		Revisions: &RevisionCache{revisions},
	}
}

func (c *Cache) enrich() {
	c.Groups.Enrich(c.Entities, c.Revisions)
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

func (f *Informer) GetEntity(ctx context.Context, name string) *typesv2.Entity {
	return f.Cache().Entities.Find(func(e *typesv2.Entity) bool {
		return e.Name == name
	})
}

func (f *Informer) ListEntity(ctx context.Context) typesv2.EntityList {
	return f.Cache().Entities.List()
}

func (f *Informer) GetFeature(ctx context.Context, name string) *typesv2.Feature {
	return f.Cache().Features.Find(func(f *typesv2.Feature) bool {
		return f.Name == name
	})
}

func (f *Informer) ListFeature(ctx context.Context, opt types.ListFeatureOpt) typesv2.FeatureList {
	return f.Cache().Features.List(opt)
}

func (f *Informer) GetFeatureGroup(ctx context.Context, name string) (*typesv2.FeatureGroup, error) {
	if featureGroup := f.Cache().Groups.Find(func(g *typesv2.FeatureGroup) bool {
		return g.Name == name
	}); featureGroup == nil {
		return nil, fmt.Errorf("feature group name=%s not found", name)
	} else {
		return featureGroup, nil
	}
}

func (f *Informer) ListFeatureGroup(ctx context.Context, entityName *string) typesv2.FeatureGroupList {
	return f.Cache().Groups.List(entityName)
}

// TODO: refactor this
func (f *Informer) GetRevision(ctx context.Context, opt metadatav2.GetRevisionOpt) (*typesv2.Revision, error) {
	if opt.RevisionId != nil {
		return f.Cache().Revisions.Find(func(r *typesv2.Revision) bool {
			return r.ID == *opt.RevisionId
		}), nil
	} else if opt.GroupName != nil && opt.Revision != nil {
		return f.Cache().Revisions.Find(func(r *typesv2.Revision) bool {
			return r.Group.Name == *opt.GroupName && r.Revision == *opt.Revision
		}), nil
	}
	return nil, fmt.Errorf("invalid GetRevisionOpt: %+v", opt)
}

func (f *Informer) ListRevision(ctx context.Context, opt metadatav2.ListRevisionOpt) typesv2.RevisionList {
	return f.Cache().Revisions.List(opt)
}

// TODO: not necessary anymore ?
func (f *Informer) GetLatestRevision(ctx context.Context, groupName string) *typesv2.Revision {
	return f.Cache().Revisions.MaxRevision(groupName)
}

// TODO: refactor this into a private function of OomStore
func (f *Informer) BuildRevisionRanges(ctx context.Context, groupName string) []*types.RevisionRange {
	return f.Cache().Revisions.BuildRevisionRanges(groupName)
}
