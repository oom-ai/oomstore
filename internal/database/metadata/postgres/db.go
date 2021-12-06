package postgres

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/metadata/sqlutil"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func (db *DB) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
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

func (db *DB) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	return sqlutil.ListGroup(ctx, db, entityID, groupIDs)
}

func (db *DB) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	return createFeature(ctx, db, opt)
}

func (db *DB) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return sqlutil.UpdateFeature(ctx, db, opt)
}

func (db *DB) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return sqlutil.GetFeature(ctx, db, id)
}

func (db *DB) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	return sqlutil.ListFeature(ctx, db, opt)
}

func (db *DB) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	return sqlutil.GetFeatureByName(ctx, db, name)
}

func (db *DB) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	var (
		revisionID int
		dataTable  string
		err        error
	)
	err = db.WithTransaction(ctx, func(c context.Context, tx metadata.DBStore) error {
		revisionID, dataTable, err = tx.CreateRevision(ctx, opt)
		return err
	})
	return revisionID, dataTable, err
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
