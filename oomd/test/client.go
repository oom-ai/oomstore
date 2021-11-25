package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

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

	importClient, err := c.ChannelImport(ctx)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("input_file_path")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		if err := importClient.Send(&codegen.ChannelImportRequest{
			GroupName:   "please input your group name",
			Description: "please input you description",
			Row:         fileScanner.Bytes(),
		}); err != nil {
			log.Fatal(err)
		}
	}

	if err := fileScanner.Err(); err != nil {
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
