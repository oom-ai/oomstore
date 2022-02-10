package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/oom-ai/oomstore/oomagent/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/version"
)

var cfgFile string
var defaultCfgFile = filepath.Join(xdg.Home, ".config", "oomstore", "config.yaml")

var oomStoreCfg types.OomStoreConfig
var port int

// rootCmd represents the base command when called without and subcommands
var rootCmd = &cobra.Command{
	Use:     "oomagent",
	Short:   "oomstore daemon",
	Version: version.String(),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomstore, err := oomstore.Open(ctx, oomStoreCfg)
		if err != nil {
			log.Fatal(err)
		}
		defer oomstore.Close()

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// write the listening address
		tmpdir := filepath.Join(os.TempDir(), "oomagent", strconv.Itoa(os.Getpid()))
		if err := os.MkdirAll(tmpdir, 0755); err != nil {
			log.Fatalf("failed to create temp directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmpdir, "address"), []byte(lis.Addr().String()), 0644); err != nil {
			log.Fatalf("failed to write listen address: %v", err)
		}

		grpcServer := grpc.NewServer()
		codegen.RegisterOomAgentServer(grpcServer, &server{oomstore: oomstore})

		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		go func() {
			<-exit
			grpcServer.GracefulStop()
			oomstore.Close()
			os.Exit(0)
		}()

		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	pFlags := rootCmd.Flags()
	pFlags.StringVar(&cfgFile, "config", defaultCfgFile, "config file")
	pFlags.IntVarP(&port, "port", "p", 50051, "The server port")

	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}

// initConfig reads in config file
func initConfig() {
	if envCfgFile := os.Getenv("OOMAGENT_CONFIG"); envCfgFile != "" {
		cfgFile = envCfgFile
	}
	cfgContent, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading config file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(cfgContent, &oomStoreCfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed loading config : %v\n", err)
		os.Exit(1)
	}
}
