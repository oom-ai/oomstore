package sqlutil

import "fmt"

func OfflineBatchTableName(groupID int, revisionID int64) string {
	return fmt.Sprintf("offline_batch_%d_%d", groupID, revisionID)
}
