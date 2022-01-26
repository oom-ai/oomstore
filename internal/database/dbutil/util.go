package dbutil

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type RowMap = map[string]interface{}

func TempTable(prefix string) string {
	return fmt.Sprintf("tmp_%s_%d", prefix, time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(b)
}

func Fill(size int, elem string, sep string) string {
	r := make([]string, 0, size)

	for i := 0; i < size; i++ {
		r = append(r, elem)
	}

	return strings.Join(r, sep)
}

func OfflineBatchTableName(groupID int, revisionID int64) string {
	return fmt.Sprintf("offline_batch_%d_%d", groupID, revisionID)
}

func OfflineStreamSnapshotTableName(groupID int, revision int64) string {
	return fmt.Sprintf("offline_stream_snapshot_%d_%d", groupID, revision)
}

func OfflineStreamCdcTableName(groupID int, revision int64) string {
	return fmt.Sprintf("offline_stream_cdc_%d_%d", groupID, revision)
}

func OnlineBatchTableName(revisionID int) string {
	return fmt.Sprintf("online_batch_%d", revisionID)
}

func OnlineStreamTableName(groupID int) string {
	return fmt.Sprintf("online_stream_%d", groupID)
}
