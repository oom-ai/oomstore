package tikv_test

// func prepareStore(t *testing.T) (context.Context, online.Store) {
// 	ctx := context.Background()
// 	store, err := tikv.Open(runtime_tikv.GetOpt())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return ctx, store
// }

// func destroyStore(t *testing.T) func() {
// 	return func() {
// 		store, err := tikv.Open(runtime_tikv.GetOpt())
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer store.Close()
// 		if err = store.DeleteRange(context.Background(), []byte(""), []byte{}); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// }

// func TestOpen(t *testing.T) {
// 	test_impl.TestOpen(t, prepareStore, destroyStore(t))
// }

// func TestPing(t *testing.T) {
// 	test_impl.TestPing(t, prepareStore, destroyStore(t))
// }

// func TestGetExisted(t *testing.T) {
// 	test_impl.TestGetExisted(t, prepareStore, destroyStore(t))
// }

// func TestGetNoRevision(t *testing.T) {
// 	test_impl.TestGetNoRevision(t, prepareStore, destroyStore(t))
// }

// func TestGetNotExistedEntityKey(t *testing.T) {
// 	test_impl.TestGetNotExistedEntityKey(t, prepareStore, destroyStore(t))
// }

// func TestMultiGet(t *testing.T) {
// 	test_impl.TestMultiGet(t, prepareStore, destroyStore(t))
// }

// func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
// 	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, destroyStore(t))
// }

// func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
// 	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, destroyStore(t))
// }

// func TestPush(t *testing.T) {
// 	test_impl.TestPush(t, prepareStore, destroyStore(t))
// }
