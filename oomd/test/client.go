package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/oom-ai/oomstore/oomd/codegen"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// The method temporarily serves the testing purpose while we iterate fast
// Feel free to add what you want to test here
// We will delete this file anyway once we have more formal tests
func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := codegen.NewOomDClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.OnlineGet(ctx, &codegen.OnlineGetRequest{EntityKey: "1001", FeatureNames: []string{"state", "has_2fa_installed"}})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("Got: %v", r)
}
