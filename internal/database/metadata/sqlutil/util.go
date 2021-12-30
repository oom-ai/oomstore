package sqlutil

import "fmt"

func OfflineBatchTableName(groupID int, revisionID int64) string {
	return fmt.Sprintf("offline_batch_%d_%d", groupID, revisionID)
}

func OfflineStreamSnapshotTableName(groupID int, revision int64) string {
	return fmt.Sprintf("offline_stream_snapshot_%d_%d", groupID, revision)
}

func OfflineStreamCdcTableName(groupID int, revision int64) string {
	return fmt.Sprintf("offline_stream_cdc_%d_%d", groupID, revision)
}
