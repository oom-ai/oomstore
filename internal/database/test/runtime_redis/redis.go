package runtime_redis

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func init() {
	opt := GetOpt()
	if out, err := exec.Command(
		"oomplay", "init", "redis",
		"--port", opt.Port,
		"--password", opt.Password,
		"--database", strconv.Itoa(opt.Database),
	).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %s", err, out))
	}
}

func GetOpt() *types.RedisOpt {
	return &types.RedisOpt{
		Host:     "localhost",
		Port:     "6379",
		Password: "test",
		Database: 0,
	}
}
