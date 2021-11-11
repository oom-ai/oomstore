package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) ListRevision(ctx context.Context, id *int16) typesv2.RevisionList {
	return s.metadatav2.ListRevision(ctx, metadatav2.ListRevisionOpt{GroupID: id})
}

func (s *OomStore) GetRevision(ctx context.Context, opt metadatav2.GetRevisionOpt) (*typesv2.Revision, error) {
	return s.metadatav2.GetRevision(ctx, metadatav2.GetRevisionOpt{
		GroupID:    opt.GroupID,
		Revision:   opt.Revision,
		RevisionId: opt.RevisionId,
	})
}
