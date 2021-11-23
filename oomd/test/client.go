package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/oom-ai/oomstore/oomd/codegen"
	"google.golang.org/grpc"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	function = flag.String("func", "", "the function to be tested")
)

// The method temporarily serves the testing purpose while we iterate fast
// Feel free to add what you want to test here
// We will delete this file anyway once we have more formal tests
func main() {
	flag.Parse()

	switch *function {
	case "online-get":
		OnlineGet()
	case "import":
		Import()
	default:
		log.Fatalf("invalid function: %s", *function)
	}
}

func OnlineGet() {
	c, cancel := prepareOomDClient(*addr)
	defer cancel()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.OnlineGet(ctx, &codegen.OnlineGetRequest{EntityKey: "1", FeatureNames: []string{"model"}})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("Got: %v", r)
}

func Import() {
	c, cancel := prepareOomDClient(*addr)
	defer cancel()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	importClient, err := c.Import(ctx)
	if err != nil {
		log.Fatal(err)
	}

	inbytes, err := ioutil.ReadFile("path-to-data-file")
	if err != nil {
		log.Fatal(err)
	}
	if err := importClient.Send(&codegen.ImportRequest{
		GroupName:   "please input your group name",
		Description: "please input you description",
		Row: []*any.Any{{
			Value: inbytes,
		}},
	}); err != nil {
		log.Fatal(err)
	}

	importRes, err := importClient.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	log.Printf("Import Got: %v", importRes)
}

func prepareOomDClient(addr string) (c codegen.OomDClient, cancel func()) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return codegen.NewOomDClient(conn), func() { conn.Close() }
}
