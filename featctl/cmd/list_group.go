package cmd

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

type listFeatureGroupOption struct {
	EntityName *string
}

var listFeatureGroupOpt listFeatureGroupOption

var listFeatureGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "list feature groups",
	Example: `1. featctl list group
2. featctl list group --entity device
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("entity") {
			listFeatureGroupOpt.EntityName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		groups, err := oneStore.ListFeatureGroup(ctx, listFeatureGroupOpt.EntityName)
		if err != nil {
			log.Fatal(err)
		}
		if err := printFeatureGroups(groups); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	listCmd.AddCommand(listFeatureGroupCmd)

	flags := listFeatureGroupCmd.Flags()

	listFeatureGroupOpt.EntityName = flags.StringP("entity", "", "", "use to filter groups")
}

func printFeatureGroups(featureGroups []*types.FeatureGroup) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write([]string{"Name", "Entity", "Description", "Revision", "DataTable", "CreateTime", "ModifyTime"}); err != nil {
		return err
	}

	for _, g := range featureGroups {
		revision := ""
		if g.Revision != nil {
			revision = strconv.FormatInt(*g.Revision, 10)
		}

		dataTable := ""
		if g.DataTable != nil {
			dataTable = *g.DataTable
		}

		if err := w.Write([]string{g.Name, g.EntityName, g.Description, revision, dataTable,
			g.CreateTime.Format(time.RFC3339), g.ModifyTime.Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
