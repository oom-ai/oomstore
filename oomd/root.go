package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/oom-ai/oomstore/oomd/codegen"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/version"
)

var cfgFile string
var defaultCfgFile = filepath.Join(xdg.Home, ".config", "oomd", "config.yaml")

var oomStoreCfg types.OomStoreConfig
var port int

// rootCmd represents the base command when called without and subcommands
var rootCmd = &cobra.Command{
	Use:     "oomd",
	Short:   "a cli tool that lets you start oom feature store grpc backend.",
	Version: version.String(),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore, err := oomstore.Open(ctx, oomStoreCfg)
		if err != nil {
			log.Fatal(err)
		}
		defer oomStore.Close()

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		codegen.RegisterOomDServer(grpcServer, &server{oomstore: oomStore})
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
	if envCfgFile := os.Getenv("OOMD_CONFIG"); envCfgFile != "" {
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
