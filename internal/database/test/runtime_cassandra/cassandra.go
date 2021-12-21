package runtime_cassandra

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/cassandra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var CassandraDbOpt types.CassandraOpt

func init() {
	cassandraContainer, err := gnomock.Start(cassandra.Preset(
		cassandra.WithVersion("4.0"),
	),
		gnomock.WithUseLocalImagesFirst(),
	)
	if err != nil {
		panic(err)
	}

	CassandraDbOpt = types.CassandraOpt{
		Hosts:    []string{cassandraContainer.DefaultAddress()},
		User:     cassandra.DefaultUser,
		Password: cassandra.DefaultPassword,
		KeySpace: "test",
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c

		_ = gnomock.Stop(cassandraContainer)
	}()
}

func PrepareDB(t *testing.T) (context.Context, *gocql.Session) {
	ctx := context.Background()

	cluster := gocql.NewCluster(CassandraDbOpt.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: CassandraDbOpt.User,
		Password: CassandraDbOpt.Password,
	}
	cluster.Timeout = time.Second * 5

	session, err := cluster.CreateSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Query("DROP KEYSPACE IF EXISTS test").WithContext(ctx).Exec(); err != nil {
		t.Fatal(err)
	}
	return ctx, session
}
