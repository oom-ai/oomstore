package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

var joinOpt JoinOpt
var joinOutput *string

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join training label data set with historical feature values",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("output") {
			joinOutput = stringPtr(ASCIITable)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		if err := join(ctx, oomStore, joinOpt, *joinOutput); err != nil {
			log.Fatalf("failed joining historical features: %+v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)

	flags := joinCmd.Flags()

	flags.StringVar(&joinOpt.InputFilePath, "input-file", "", "file path of training label data set")
	_ = joinCmd.MarkFlagRequired("input-file")

	flags.StringSliceVar(&joinOpt.FeatureNames, "feature", nil, "feature names")
	_ = joinCmd.MarkFlagRequired("feature")

	joinOutput = flags.StringP("output", "o", "", "output format")
}
