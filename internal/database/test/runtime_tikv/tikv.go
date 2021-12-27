package runtime_tikv

import (
	"fmt"
	"os/exec"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	if out, err := exec.Command("oomplay", "init", "tikv").CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func GetOpt() *types.TiKVOpt {
	return &types.TiKVOpt{
		PdAddrs: []string{"127.0.0.1:22379"},
	}
}
