package tidb_test

// var DATABASE string

// func init() {
// 	DATABASE = strings.ToLower(dbutil.RandString(20))
// }

// func prepareStore(t *testing.T) (context.Context, online.Store) {
// 	runtime_tidb.CreateDatabase(DATABASE)

// 	store, err := mysql.Open(runtime_tidb.GetOpt(DATABASE))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return context.Background(), store
// }

// func TestOpen(t *testing.T) {
// 	test_impl.TestOpen(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestGetExisted(t *testing.T) {
// 	test_impl.TestGetExisted(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestGetNoRevision(t *testing.T) {
// 	test_impl.TestGetNoRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestGetNotExistedEntityKey(t *testing.T) {
// 	test_impl.TestGetNotExistedEntityKey(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestMultiGet(t *testing.T) {
// 	test_impl.TestMultiGet(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
// 	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
// 	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestPush(t *testing.T) {
// 	test_impl.TestPush(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestPing(t *testing.T) {
// 	test_impl.TestPing(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }
