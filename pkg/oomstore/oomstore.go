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

// Combine online store, offline store, and metadata store instances in one place.
type OomStore struct {
	online   online.Store
	offline  offline.Store
	metadata metadata.Store

	pushProcessor *PushProcessor
}

// Return an OomStore instance given the configuration.
// Under the hood, it setups connections to the underlying databases.
// You should always use this method to create a new OomStore instance in code.
func Open(ctx context.Context, opt types.OomStoreConfig) (*OomStore, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	onlineStore, err := database.OpenOnlineStore(opt.OnlineStore)
	if err != nil {
		return nil, err
	}
	offlineStore, err := database.OpenOfflineStore(ctx, opt.OfflineStore)
	if err != nil {
		return nil, err
	}
	metadataStore, err := database.OpenMetadataStore(opt.MetadataStore)
	if err != nil {
		return nil, err
	}

	store := &OomStore{
		online:   onlineStore,
		offline:  offlineStore,
		metadata: metadataStore,
	}
	store.InitPushProcessor(ctx, opt.PushProcessor)

	return store, nil
}

// Create a new OomStore instance.
func Create(ctx context.Context, opt types.OomStoreConfig) (*OomStore, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	if err := database.CreateMetadataDatabase(ctx, opt.MetadataStore); err != nil {
		return nil, err
	}
	return Open(ctx, opt)
}

// Ping verifies the connections to the backend stores are still alive
func (s *OomStore) Ping(ctx context.Context) error {
	if err := s.online.Ping(ctx); err != nil {
		return err
	}
	if err := s.offline.Ping(ctx); err != nil {
		return err
	}
	if err := s.metadata.Ping(ctx); err != nil {
		return err
	}
	return nil
}

// Close the connections to underlying databases.
func (s *OomStore) Close() error {
	errs := make([]error, 0)

	for _, closer := range []io.Closer{s.pushProcessor, s.online, s.offline, s.metadata} {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed closing store: %v", errs)
	}
	return nil
}

// Return an OomStore instance for internal test purpose ONLY.
// You should NOT use this method directly in any of your code.
func TEST__New(online online.Store, offline offline.Store, metadata metadata.Store) *OomStore {
	return &OomStore{
		online:   online,
		offline:  offline,
		metadata: metadata,
	}
}
