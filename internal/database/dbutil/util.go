package dbutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type RowMap = map[string]interface{}

func TempTable(prefix string) string {
	return fmt.Sprintf("tmp_%s_%d", prefix, time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randInt(size int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(size)))
	return int(n.Int64())
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randInt(len(letterRunes))]
	}
	return string(b)
}

func Fill(size int, elem, sep string) string {
	r := make([]string, 0, size)

	for i := 0; i < size; i++ {
		r = append(r, elem)
	}

	return strings.Join(r, sep)
}

func OfflineBatchSnapshotTableName(groupID int, revisionID int64) string {
	return fmt.Sprintf("offline_batch_snapshot_%d_%d", groupID, revisionID)
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
