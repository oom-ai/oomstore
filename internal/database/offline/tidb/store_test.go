package tidb_test

// var DATABASE string

// func init() {
// 	DATABASE = strings.ToLower(dbutil.RandString(20))
// }

// func prepareStore(t *testing.T) (context.Context, offline.Store) {
// 	runtime_tidb.CreateDatabase(DATABASE)

// 	store, err := mysql.Open(runtime_tidb.GetOpt(DATABASE))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return context.Background(), store
// }

// func TestPing(t *testing.T) {
// 	test_impl.TestPing(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestExport(t *testing.T) {
// 	test_impl.TestExport(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestImport(t *testing.T) {
// 	test_impl.TestImport(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestJoin(t *testing.T) {
// 	test_impl.TestJoin(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
// }

// func TestTableSchema(t *testing.T) {
// 	test_impl.TestTableSchema(t, prepareStore, runtime_tidb.DestroyStore(DATABASE), func(ctx context.Context) {
// 		db, err := mysql.Open(runtime_tidb.GetOpt(DATABASE))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer db.Close()
// 		if _, err = db.ExecContext(ctx, "create table `offline_batch_1_1`(`user` varchar(16), `age` smallint, `unix_milli` int)"); err != nil {
// 			t.Fatal(err)
// 		}
// 		if _, err = db.ExecContext(ctx, "insert into `offline_batch_1_1` VALUES ('1', 1, 1), ('2', 2, 100)"); err != nil {
// 			t.Fatal(err)
// 		}
// 	})
// }
