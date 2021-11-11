package informer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
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

// Get
func (f *Informer) GetEntity(ctx context.Context, id int16) (*typesv2.Entity, error) {
	if entity := f.Cache().Entities.Find(func(e *typesv2.Entity) bool {
		return e.ID == id
	}); entity == nil {
		return nil, fmt.Errorf("feature entity %d not found", id)
	} else {
		return entity, nil
	}
}

func (f *Informer) GetFeature(ctx context.Context, id int16) (*typesv2.Feature, error) {
	if feature := f.Cache().Features.Find(func(f *typesv2.Feature) bool {
		return f.ID == id
	}); feature == nil {
		return nil, fmt.Errorf("feature %d not found", id)
	} else {
		return feature, nil
	}
}

func (f *Informer) GetFeatureGroup(ctx context.Context, id int16) (*typesv2.FeatureGroup, error) {
	if featureGroup := f.Cache().Groups.Find(func(g *typesv2.FeatureGroup) bool {
		return g.ID == id
	}); featureGroup == nil {
		return nil, fmt.Errorf("feature group %d not found", id)
	} else {
		return featureGroup, nil
	}
}

func (f *Informer) GetRevision(ctx context.Context, opt metadatav2.GetRevisionOpt) (*typesv2.Revision, error) {
	if opt.RevisionId != nil {
		if opt.GroupID != nil || opt.Revision != nil {
			return nil, fmt.Errorf("invalid GetRevisionOpt: %+v", opt)
		}
		revision := f.Cache().Revisions.Find(func(r *typesv2.Revision) bool {
			return r.ID == *opt.RevisionId
		})
		if revision == nil {
			return nil, fmt.Errorf("revision not found")
		}
		return revision, nil
	} else if opt.GroupID != nil && opt.Revision != nil {
		revision := f.Cache().Revisions.Find(func(r *typesv2.Revision) bool {
			return r.Group.ID == *opt.GroupID && r.Revision == *opt.Revision
		})
		if revision == nil {
			return nil, fmt.Errorf("revision not found")
		}
		return revision, nil
	}
	return nil, fmt.Errorf("invalid GetRevisionOpt: %+v", opt)
}

// List
func (f *Informer) ListEntity(ctx context.Context) typesv2.EntityList {
	return f.Cache().Entities.List()
}

func (f *Informer) ListFeature(ctx context.Context, opt metadatav2.ListFeatureOpt) typesv2.FeatureList {
	return f.Cache().Features.List(opt)
}

func (f *Informer) ListFeatureGroup(ctx context.Context, entityID *int16) typesv2.FeatureGroupList {
	return f.Cache().Groups.List(entityID)
}

func (f *Informer) ListRevision(ctx context.Context, opt metadatav2.ListRevisionOpt) typesv2.RevisionList {
	return f.Cache().Revisions.List(opt)
}

// TODO: not necessary anymore ?
func (f *Informer) GetLatestRevision(ctx context.Context, groupID int16) *typesv2.Revision {
	return f.Cache().Revisions.MaxRevision(groupID)
}

// TODO: refactor this into a private function of OomStore
func (f *Informer) BuildRevisionRanges(ctx context.Context, groupID int16) []*metadatav2.RevisionRange {
	return f.Cache().Revisions.BuildRevisionRanges(groupID)
}
