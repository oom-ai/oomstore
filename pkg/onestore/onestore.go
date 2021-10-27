package onestore

import (
	"context"
	"fmt"
	"io"

	"github.com/oom-ai/oomstore/internal/database"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

type OneStore struct {
	online   online.Store
	offline  offline.Store
	metadata metadata.Store
}

func Open(ctx context.Context, opt types.OneStoreOpt) (*OneStore, error) {
	optV2 := opt.ToOneStoreOptV2()

	onlineStore, err := database.OpenOnlineStore(optV2.OnlineStoreOpt)
	if err != nil {
		return nil, err
	}
	offlineStore, err := database.OpenOfflineStore(optV2.OfflineStoreOpt)
	if err != nil {
		return nil, err
	}
	metadataStore, err := database.OpenMetadataStore(optV2.MetaStoreOpt)
	if err != nil {
		return nil, err
	}

	return &OneStore{
		online:   onlineStore,
		offline:  offlineStore,
		metadata: metadataStore,
	}, nil
}

func Create(ctx context.Context, opt types.OneStoreOpt) (*OneStore, error) {
	optV2 := opt.ToOneStoreOptV2()
	if err := database.CreateMetadataDatabase(ctx, optV2.MetaStoreOpt); err != nil {
		return nil, err
	}

	return Open(ctx, opt)
}

func (s *OneStore) Close() error {
	errs := []error{}

	for _, closer := range []io.Closer{s.online, s.offline, s.metadata} {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed closing store: %v", errs)
	}
	return nil
}
