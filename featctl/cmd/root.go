package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/version"
)

var cfgFile string
var defaultCfgFile = filepath.Join(xdg.ConfigHome, "featctl", "config.yaml")

var oomStoreOpt types.OomStoreOpt

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "featctl",
	Short:   "a cli tool lets you control the oom feature store.",
	Version: version.String(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// panic("sjdfiod")
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

	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVar(&cfgFile, "config", defaultCfgFile, "config file")
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}

// initConfig reads in config file
func initConfig() {
	if envCfgFile := os.Getenv("FEATCTL_CONFIG"); envCfgFile != "" {
		cfgFile = envCfgFile
	}
	cfgContent, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(fmt.Sprintf("failed reading config file: %v\n", err))
	}

	if err := yaml.Unmarshal(cfgContent, &oomStoreOpt); err != nil {
		panic(fmt.Sprintf("failed loading config : %v\n", err))
	}
}
