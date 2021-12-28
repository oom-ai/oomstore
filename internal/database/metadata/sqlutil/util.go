package sqlutil

import "fmt"

func OfflineBatchTableName(groupID, revisionID int) string {
	return fmt.Sprintf("offline_batch_%d_%d", groupID, revisionID)
}
