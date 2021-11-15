package oomstore

import (
	"context"
	"fmt"
	"io"

	"github.com/oom-ai/oomstore/internal/database"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type OomStore struct {
	online     online.Store
	offline    offline.Store
	metadatav2 metadata.Store
}

func NewOomStore(online online.Store, offline offline.Store, metadatav2 metadata.Store) *OomStore {
	return &OomStore{
		online:     online,
		offline:    offline,
		metadatav2: metadatav2,
	}
}

func Open(ctx context.Context, opt types.OomStoreConfig) (*OomStore, error) {
	onlineStore, err := database.OpenOnlineStore(opt.OnlineStore)
	if err != nil {
		return nil, err
	}
	offlineStore, err := database.OpenOfflineStore(opt.OfflineStore)
	if err != nil {
		return nil, err
	}
	metadatav2Store, err := database.OpenMetadatav2Store(opt.MetadataStore)
	if err != nil {
		return nil, err
	}

	return &OomStore{
		online:     onlineStore,
		offline:    offlineStore,
		metadatav2: metadatav2Store,
	}, nil
}

func Create(ctx context.Context, opt types.OomStoreConfig) (*OomStore, error) {
	if err := database.CreateMetadataDatabase(ctx, opt.MetadataStore); err != nil {
		return nil, err
	}
	return Open(ctx, opt)
}

func (s *OomStore) Close() error {
	errs := []error{}

	for _, closer := range []io.Closer{s.online, s.offline, s.metadatav2} {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed closing store: %v", errs)
	}
	return nil
}
