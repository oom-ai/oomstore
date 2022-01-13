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
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	Day       = 24 * time.Hour
	Capacity  = 1000
	Period    = 5 * time.Minute
	MinPeriod = 2 * time.Minute
)

type GroupBuffer struct {
	groupID  int
	records  []types.StreamRecord
	lastPush time.Time
}

type StreamPushProcessor struct {
	capacity   int
	minPeriod  time.Duration
	notifyQuit chan struct{}
	waitQuit   chan struct{}

	ticker *time.Ticker
	ch     chan types.StreamRecord
	buffer sync.Map // GroupID -> GroupBuffer
}

func (s *OomStore) InitStreamPushProcessor(ctx context.Context) {
	processor := &StreamPushProcessor{
		// TODO: make Capacity, Period, MinPeriod configurable
		capacity:   Capacity,
		minPeriod:  MinPeriod,
		notifyQuit: make(chan struct{}),
		waitQuit:   make(chan struct{}),

		ch:     make(chan types.StreamRecord),
		ticker: time.NewTicker(Period),
	}
	s.streamPushProcessor = processor

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
					if len(b.records) > 0 && time.Since(b.lastPush) > processor.minPeriod {
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

				if len(b.records) >= processor.capacity {
					if err := processor.pushToOffline(ctx, s, groupID); err != nil {
						fmt.Fprintf(os.Stderr, "Error pushing to offline store: %+v", err)
					}
				}
			}
		}
	}()
}

func (p *StreamPushProcessor) Close() error {
	p.ticker.Stop()
	p.notifyQuit <- struct{}{}

	<-p.waitQuit
	return nil
}

func (p *StreamPushProcessor) Push(record types.StreamRecord) {
	if _, ok := p.buffer.Load(record.GroupID); !ok {
		p.buffer.Store(record.GroupID, GroupBuffer{
			groupID: record.GroupID,
			records: make([]types.StreamRecord, 0, p.capacity),
		})
	}
	p.ch <- record
}

func (p *StreamPushProcessor) pushToOffline(ctx context.Context, s *OomStore, groupID int) error {
	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &groupID,
	})
	if err != nil {
		return err
	}
	entity := features[0].Entity()

	value, _ := p.buffer.Load(groupID)
	b := value.(GroupBuffer)

	buckets := make(map[int64][]types.StreamRecord)
	for _, record := range b.records {
		revision := latestRevision(record.UnixMilli)
		if _, ok := buckets[revision]; !ok {
			buckets[revision] = make([]types.StreamRecord, 0)
		}
		buckets[revision] = append(buckets[revision], record)
	}
	for revision, records := range buckets {
		err = s.offline.Push(ctx, offline.PushOpt{
			GroupID:      groupID,
			Revision:     revision,
			EntityName:   entity.Name,
			FeatureNames: features.Names(),
			Records:      records,
		})
		if err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}
			if err = p.newRevision(ctx, s, groupID, revision); err != nil {
				return err
			}
		}
	}

	b.records = make([]types.StreamRecord, 0, p.capacity)
	b.lastPush = time.Now()
	p.buffer.Store(groupID, b)
	return err
}

func (p *StreamPushProcessor) newRevision(ctx context.Context, s *OomStore, groupID int, revision int64) error {
	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &groupID,
	})
	if err != nil {
		return err
	}
	entity := features[0].Entity()

	if err = s.offline.CreateTable(ctx, offline.CreateTableOpt{
		TableName: dbutil.OfflineStreamCdcTableName(groupID, revision),
		Entity:    entity,
		Features:  features,
		IsCDC:     true,
	}); err != nil {
		return err
	}

	snapshotTable := ""
	cdcTable := dbutil.OfflineStreamCdcTableName(groupID, revision)
	if _, _, err = s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		GroupID:       groupID,
		Revision:      revision,
		SnapshotTable: &snapshotTable,
		CdcTable:      &cdcTable,
	}); err != nil {
		return err
	}
	return nil
}

func latestRevision(unixMilli int64) int64 {
	return (unixMilli / Day.Milliseconds()) * Day.Milliseconds()
}
