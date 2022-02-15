package runtime_tikv

// func use_host_tikv() bool {
// 	return os.Getenv("USE_HOST_TIKV") == "1" || runtime.GOOS != "linux"
// }

// func init() {
// 	if use_host_tikv() {
// 		if _, err := exec.Command("sh", "-c", "tiup status | grep RUNNING").CombinedOutput(); err != nil {
// 			panic("cannot init: start tiup playground by `tiup playground ^5 --without-monitor --mode=tikv-slim` and retry")
// 		}
// 	} else {
// 		if out, err := exec.Command("oomplay", "init", "tikv").CombinedOutput(); err != nil {
// 			panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
// 		}
// 	}
// }

// func GetOpt() *types.TiKVOpt {
// 	if use_host_tikv() {
// 		return &types.TiKVOpt{
// 			PdAddrs: []string{"127.0.0.1:2379"},
// 		}
// 	} else {
// 		return &types.TiKVOpt{
// 			PdAddrs: []string{"127.0.0.1:22379"},
// 		}
// 	}
// }
