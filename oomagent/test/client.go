package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/oom-ai/oomstore/oomagent/codegen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	c, cancel := prepareOomAgentClient(*addr)
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
	c, cancel := prepareOomAgentClient(*addr)
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

	groupName := "please input your group name"
	description := "please input you description"
	for fileScanner.Scan() {
		if err := importClient.Send(&codegen.ChannelImportRequest{
			GroupName:   &groupName,
			Description: &description,
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

func prepareOomAgentClient(addr string) (c codegen.OomAgentClient, cancel func()) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return codegen.NewOomAgentClient(conn), func() { conn.Close() }
}
