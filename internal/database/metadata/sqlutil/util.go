package sqlutil

import "fmt"

func OfflineBatchTableName(groupID int, revisionID int64) string {
	return fmt.Sprintf("offline_batch_%d_%d", groupID, revisionID)
}

func OfflineStreamSnapshotTableName(groupID int, revisionID int) string {
	return fmt.Sprintf("offline_stream_snapshot_%d_%d", groupID, revisionID)
}

func OfflineStreamCdcTableName(groupID int, revisionID int) string {
	return fmt.Sprintf("offline_stream_cdc_%d_%d", groupID, revisionID)
}
