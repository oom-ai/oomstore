package sqlutil

import "fmt"

const BatchSize = 10

func OnlineTableName(revisionID int) string {
	return fmt.Sprintf("online_%d", revisionID)
}
