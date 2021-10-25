package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/metadata"
	"github.com/onestore-ai/onestore/internal/database/offline"
	"github.com/onestore-ai/onestore/internal/database/online"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type OneStore struct {
	db *database.DB

	online   online.Store
	offline  offline.Store
	metadata metadata.Store
}

func Open(ctx context.Context, opt types.OneStoreOpt) (*OneStore, error) {
	db, err := database.Open(toDatabaseOption(&opt))
	if err != nil {
		return nil, err
	}

	onlineStore, err := online.OpenPostgresDB(online.OnlineStoreOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Pass:     opt.Pass,
		Database: opt.Workspace,
	})
	if err != nil {
		return nil, err
	}

	return &OneStore{
		db:     db,
		online: onlineStore,
	}, nil
}

func Create(ctx context.Context, opt types.OneStoreOpt) (*OneStore, error) {
	if err := database.CreateDatabase(ctx, toDatabaseOption(&opt)); err != nil {
		return nil, err
	}

	return Open(ctx, opt)
}

func (s *OneStore) Close() error {
	return s.db.Close()
}

func toDatabaseOption(opt *types.OneStoreOpt) database.Option {
	return database.Option{
		Host:   opt.Host,
		Port:   opt.Port,
		User:   opt.User,
		Pass:   opt.Pass,
		DbName: opt.Workspace,
	}
}
