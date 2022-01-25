package oomstore

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var defaultPushProcessorCfg = types.PushProcessorConfig{
	FlushInterval: 2 * time.Minute,
	BufferSize:    1000,
}

type GroupBuffer struct {
	groupID  int
	records  []types.StreamRecord
	lastPush time.Time
}

type PushProcessor struct {
	bufferSize    int
	flushInterval time.Duration
	notifyQuit    chan struct{}
	waitQuit      chan struct{}

	ticker *time.Ticker
	ch     chan types.StreamRecord
	buffer sync.Map // GroupID -> GroupBuffer
}

func (s *OomStore) InitPushProcessor(ctx context.Context, cfg *types.PushProcessorConfig) {
	if cfg == nil {
		cfg = &defaultPushProcessorCfg
	}

	// tick at least once every 10 seconds
	maxTickInterval := 10 * time.Second
	tickInterval := cfg.FlushInterval
	if cfg.FlushInterval > maxTickInterval {
		tickInterval = maxTickInterval
	}

	processor := &PushProcessor{
		bufferSize:    cfg.BufferSize,
		flushInterval: cfg.FlushInterval,
		notifyQuit:    make(chan struct{}),
		waitQuit:      make(chan struct{}),

		ch:     make(chan types.StreamRecord),
		ticker: time.NewTicker(tickInterval),
	}
	s.pushProcessor = processor

	go func() {
		defer func() {
			processor.waitQuit <- struct{}{}
		}()

		for {
			select {
			case <-processor.notifyQuit:
				processor.buffer.Range(func(key, value interface{}) bool {
					groupID := key.(int)
					b := value.(GroupBuffer)
					if len(b.records) > 0 {
						if err := processor.pushToOffline(ctx, s, groupID); err != nil {
							fmt.Fprintf(os.Stderr, "Error pushing to offline store: %+v", err)
						}
					}
					return true
				})
				return
			case <-processor.ticker.C:
				processor.buffer.Range(func(key, value interface{}) bool {
					groupID := key.(int)
					b := value.(GroupBuffer)
					if len(b.records) > 0 && time.Since(b.lastPush) >= processor.flushInterval {
						if err := processor.pushToOffline(ctx, s, groupID); err != nil {
							fmt.Fprintf(os.Stderr, "Error pushing to offline store: %+v", err)
						}
					}
					return true
				})
			case record := <-processor.ch:
				groupID := record.GroupID
				value, _ := processor.buffer.Load(groupID)
				b := value.(GroupBuffer)

				b.records = append(b.records, record)
				processor.buffer.Store(groupID, b)

				if len(b.records) >= processor.bufferSize {
					if err := processor.pushToOffline(ctx, s, groupID); err != nil {
						fmt.Fprintf(os.Stderr, "Error pushing to offline store: %+v", err)
					}
				}
			}
		}
	}()
}

func (p *PushProcessor) Close() error {
	p.ticker.Stop()
	p.notifyQuit <- struct{}{}

	<-p.waitQuit
	return nil
}

func (p *PushProcessor) Push(record types.StreamRecord) {
	if _, ok := p.buffer.Load(record.GroupID); !ok {
		p.buffer.Store(record.GroupID, GroupBuffer{
			groupID: record.GroupID,
			records: make([]types.StreamRecord, 0, p.bufferSize),
		})
	}
	p.ch <- record
}

func (p *PushProcessor) pushToOffline(ctx context.Context, s *OomStore, groupID int) error {
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		GroupID: &groupID,
	})
	entity := features[0].Entity()

	value, _ := p.buffer.Load(groupID)
	b := value.(GroupBuffer)
	group, err := s.metadata.GetCachedGroup(ctx, groupID)
	if err != nil {
		return err
	}
	if err = s.pushStreamingRecords(ctx, b.records, entity.Name, group, features); err != nil {
		return err
	}

	b.records = make([]types.StreamRecord, 0, p.bufferSize)
	b.lastPush = time.Now()
	p.buffer.Store(groupID, b)
	return err
}

func (s *OomStore) newRevisionForStream(ctx context.Context, groupID int, revision int64) error {
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		GroupID: &groupID,
	})
	entity := features[0].Entity()

	if err := s.offline.CreateTable(ctx, offline.CreateTableOpt{
		TableName: dbutil.OfflineStreamCdcTableName(groupID, revision),
		Entity:    entity,
		Features:  features,
		TableType: types.TableStreamCdc,
	}); err != nil {
		return err
	}

	snapshotTable := ""
	cdcTable := dbutil.OfflineStreamCdcTableName(groupID, revision)
	if _, _, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
		GroupID:       groupID,
		Revision:      revision,
		SnapshotTable: &snapshotTable,
		CdcTable:      &cdcTable,
	}); err != nil {
		return err
	}
	return nil
}

func lastRevisionForStream(snapshotInterval int64, unixMill int64) int64 {
	return (unixMill / snapshotInterval) * snapshotInterval
}
