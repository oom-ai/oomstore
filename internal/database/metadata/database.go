package metadata

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error
	GetEntity(ctx context.Context, name string) (*types.Entity, error)
	ListEntity(ctx context.Context) ([]*types.Entity, error)
	UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error
	// TODO: add all metadata methods
}
