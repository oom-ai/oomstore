package runtime_redis

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
)

var RedisDbOpt types.RedisOpt

func init() {
	container, err := gnomock.Start(
		redis.Preset(
			redis.WithVersion("6.2.6"),
		),
		gnomock.WithUseLocalImagesFirst(),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c
		_ = gnomock.Stop(container)
	}()

	RedisDbOpt = types.RedisOpt{
		Host:     container.Host,
		Port:     strconv.Itoa(container.DefaultPort()),
		Password: "",
		Database: 0,
	}
}
