package runtime_cassandra

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/gocql/gocql"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func PrepareDB() (context.Context, *gocql.Session) {
	ctx := context.Background()
	opt := GetOpt("")

	cluster := gocql.NewCluster(opt.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: opt.User,
		Password: opt.Password,
	}
	cluster.Timeout = time.Second * 5

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return ctx, session
}

func init() {
	// "dummy" db will not actually be used during testing
	opt := GetOpt("dummy")
	if out, err := exec.Command(
		"oomplay", "init", "cassandra",
		"--port",
		"9042",
		"--keyspace",
		opt.KeySpace,
	).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("oomplay failed with error: %v, output: %v", err, out))
	}
}

func GetOpt(keySpace string) *types.CassandraOpt {
	return &types.CassandraOpt{
		Hosts:    []string{"127.0.0.1:9042"},
		User:     "test",
		Password: "test",
		KeySpace: keySpace,
	}
}
