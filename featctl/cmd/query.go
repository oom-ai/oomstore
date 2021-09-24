package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/onestore-ai/onestore/featctl/pkg/query"
)

var queryOpt query.Option

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query feature values",
	Example: `
1. featctl query --group user_info -n sex,city
2. featctl query --group user_info -n sex,'user name'
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		queryOpt.DBOption = dbOption

		featurenames := make([]string, 0, cap(queryOpt.FeatureNames))
		for _, name := range queryOpt.FeatureNames {
			if name != "entity_key" {
				featurenames = append(featurenames, name)
			}
		}
		queryOpt.FeatureNames = featurenames
	},
	Run: func(cmd *cobra.Command, args []string) {
		query.Run(context.Background(), &queryOpt)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	flags := queryCmd.Flags()

	flags.StringVarP(&queryOpt.Group, "group", "g", "", "feature group")
	_ = queryCmd.MarkFlagRequired("group")

	flags.StringVarP(&queryOpt.Revision, "revision", "r", "", "revision")

	// https://pkg.go.dev/github.com/spf13/pflag#StringSlice
	flags.StringSliceVarP(&queryOpt.Entitykeys, "key", "k", nil, "entity keys")
	flags.StringSliceVarP(&queryOpt.FeatureNames, "name", "n", nil, "feature names")
}
