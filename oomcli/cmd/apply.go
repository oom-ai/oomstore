package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

type ApplyOption struct {
	FilePath string
}

var applyOpt ApplyOption

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a change",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		file, err := os.Open(applyOpt.FilePath)
		if err != nil {
			exit(err)
		}
		defer file.Close()

		if err := oomStore.Apply(ctx, apply.ApplyOpt{R: file}); err != nil {
			exitf("apply failed: %+v", err)
		}

		fmt.Fprintln(os.Stderr, "applied")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	flags := applyCmd.Flags()

	flags.StringVarP(&applyOpt.FilePath, "filepath", "f", "", "filepath")
}
