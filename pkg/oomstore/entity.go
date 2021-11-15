package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) GetEntity(ctx context.Context, id int16) (*typesv2.Entity, error) {
	return s.metadatav2.GetEntity(ctx, id)
}

func (s *OomStore) GetEntityByName(ctx context.Context, name string) (*typesv2.Entity, error) {
	return s.metadatav2.GetEntityByName(ctx, name)
}

func (s *OomStore) ListEntity(ctx context.Context) typesv2.EntityList {
	return s.metadatav2.ListEntity(ctx)
}

func (s *OomStore) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	return s.metadatav2.CreateEntity(ctx, opt)
}

func (s *OomStore) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return s.metadatav2.UpdateEntity(ctx, opt)
}
