package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func Open(ctx context.Context, opt types.OneStoreOpt) (*OneStore, error) {
	dbOpt := database.Option{
		Host:   opt.Host,
		Port:   opt.Port,
		User:   opt.User,
		Pass:   opt.Pass,
		DbName: opt.Workspace,
	}
	db, err := database.Open(dbOpt)
	if err != nil {
		return nil, err
	}

	return &OneStore{db}, nil
}
