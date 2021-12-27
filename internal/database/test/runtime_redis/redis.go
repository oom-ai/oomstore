package runtime_redis

import (
	"fmt"
	"os/exec"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	if out, err := exec.Command("oomplay", "init", "redis").CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func GetOpt() *types.RedisOpt {
	return &types.RedisOpt{
		Host: "localhost",
		Port: "26379",
	}
}
