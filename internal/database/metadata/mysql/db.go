package mysql

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	ER_DUP_ENTRY = 1062
)

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := &Tx{Tx: tx}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fn(ctx, txStore)
}

func (db *DB) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	return createEntity(ctx, db, opt)
}

func (db *DB) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return sqlutil.UpdateEntity(ctx, db, opt)
}

func (db *DB) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return sqlutil.GetEntity(ctx, db, id)
}

func (db *DB) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return sqlutil.GetEntityByName(ctx, db, name)
}

func (db *DB) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	return sqlutil.ListEntity(ctx, db, entityIDs)
}

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	var featureID int
	var err error
	err = db.WithTransaction(ctx, func(c context.Context, tx metadata.DBStore) error {
		featureID, err = tx.CreateFeature(c, opt)
		return err
	})

	return featureID, err
}

func (db *DB) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return sqlutil.UpdateFeature(ctx, db, opt)
}

func (db *DB) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return sqlutil.GetFeature(ctx, db, id)
}

func (db *DB) GetFeatureByName(ctx context.Context, groupName, featureName string) (*types.Feature, error) {
	return sqlutil.GetFeatureByName(ctx, db, groupName, featureName)
}

func (db *DB) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	return sqlutil.ListFeature(ctx, db, opt)
}

func (db *DB) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	return createGroup(ctx, db, opt)
}

func (db *DB) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	return sqlutil.UpdateGroup(ctx, db, opt)
}

func (db *DB) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return sqlutil.GetGroup(ctx, db, id)
}

func (db *DB) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return sqlutil.GetGroupByName(ctx, db, name)
}

func (db *DB) ListGroup(ctx context.Context, opt metadata.ListGroupOpt) (types.GroupList, error) {
	return sqlutil.ListGroup(ctx, db, opt)
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, error) {
	var (
		revisionID int
		err        error
	)
	err = db.WithTransaction(ctx, func(c context.Context, tx metadata.DBStore) error {
		revisionID, err = tx.CreateRevision(c, opt)
		return err
	})
	return revisionID, err
}

func (db *DB) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	return sqlutil.UpdateRevision(ctx, db, opt)
}

func (db *DB) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return sqlutil.GetRevision(ctx, db, id)
}

func (db *DB) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	return sqlutil.GetRevisionBy(ctx, db, groupID, revision)
}

func (db *DB) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	return sqlutil.ListRevision(ctx, db, groupID)
}
