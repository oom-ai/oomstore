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

type listRevisionOption struct {
	GroupName *string
}

var listRevisionOpt listRevisionOption

var listRevisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "list historical revisions given a specific group",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("group") {
			listRevisionOpt.GroupName = nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oneStore := mustOpenOneStore(ctx, oneStoreOpt)
		defer oneStore.Close()

		revisions, err := oneStore.ListRevision(ctx, listRevisionOpt.GroupName)
		if err != nil {
			log.Fatal(err)
		}

		if err := printRevision(revisions); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	listCmd.AddCommand(listRevisionCmd)

	flags := listRevisionCmd.Flags()
	listRevisionOpt.GroupName = flags.StringP("group", "g", "", "feature group")
}

func printRevision(revisions []*types.Revision) error {
	w := csv.NewWriter(os.Stdout)

	if err := w.Write([]string{"Revision", "GroupName", "DataTable", "Description", "CreateTime", "ModifyTime"}); err != nil {
		return err
	}

	for _, r := range revisions {
		if err := w.Write([]string{strconv.Itoa(int(r.Revision)), r.GroupName, r.DataTable, r.Description,
			r.CreateTime.Format(time.RFC3339), r.ModifyTime.Format(time.RFC3339)}); err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}
