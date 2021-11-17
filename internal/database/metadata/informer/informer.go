package informer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
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
	groups types.FeatureGroupList,
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
func (f *Informer) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	if entity := f.Cache().Entities.Find(func(e *types.Entity) bool {
		return e.ID == id
	}); entity == nil {
		return nil, fmt.Errorf("feature entity %d not found", id)
	} else {
		return entity.Copy(), nil
	}
}

func (f *Informer) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	if entity := f.Cache().Entities.Find(func(e *types.Entity) bool {
		return e.Name == name
	}); entity == nil {
		return nil, fmt.Errorf("feature entity '%s' not found", name)
	} else {
		return entity.Copy(), nil
	}
}

func (f *Informer) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	if feature := f.Cache().Features.Find(func(f *types.Feature) bool {
		return f.ID == id
	}); feature == nil {
		return nil, fmt.Errorf("feature %d not found", id)
	} else {
		return feature.Copy(), nil
	}
}

func (f *Informer) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	if feature := f.Cache().Features.Find(func(f *types.Feature) bool {
		return f.Name == name
	}); feature == nil {
		return nil, fmt.Errorf("feature '%s' not found", name)
	} else {
		return feature.Copy(), nil
	}
}

func (f *Informer) GetFeatureGroup(ctx context.Context, id int) (*types.FeatureGroup, error) {
	if featureGroup := f.Cache().Groups.Find(func(g *types.FeatureGroup) bool {
		return g.ID == id
	}); featureGroup == nil {
		return nil, fmt.Errorf("feature group %d not found", id)
	} else {
		return featureGroup.Copy(), nil
	}
}

func (f *Informer) GetFeatureGroupByName(ctx context.Context, name string) (*types.FeatureGroup, error) {
	if featureGroup := f.Cache().Groups.Find(func(g *types.FeatureGroup) bool {
		return g.Name == name
	}); featureGroup == nil {
		return nil, fmt.Errorf("feature group '%s' not found", name)
	} else {
		return featureGroup.Copy(), nil
	}
}

func (f *Informer) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	if revision := f.Cache().Revisions.Find(func(r *types.Revision) bool {
		return r.ID == id
	}); revision == nil {
		return nil, fmt.Errorf("revision not found")
	} else {
		return revision.Copy(), nil
	}
}

func (f *Informer) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	if revision := f.Cache().Revisions.Find(func(r *types.Revision) bool {
		return r.GroupID == groupID && r.Revision == revision
	}); revision == nil {
		return nil, fmt.Errorf("revision not found")
	} else {
		return revision.Copy(), nil
	}
}

// List
func (f *Informer) ListEntity(ctx context.Context) types.EntityList {
	return f.Cache().Entities.List().Copy()
}

func (f *Informer) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) types.FeatureList {
	return f.Cache().Features.List(opt).Copy()
}

func (f *Informer) ListFeatureGroup(ctx context.Context, entityID *int) types.FeatureGroupList {
	return f.Cache().Groups.List(entityID).Copy()
}

func (f *Informer) ListRevision(ctx context.Context, groupID *int) types.RevisionList {
	return f.Cache().Revisions.List(groupID).Copy()
}
