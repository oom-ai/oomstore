package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var describeGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "show details of a specific group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		oomStore := mustOpenOomStore(ctx, oomStoreCfg)
		defer oomStore.Close()

		groupName := args[0]
		group, err := oomStore.GetGroupByName(ctx, groupName)
		if err != nil {
			log.Fatalf("failed getting group %s, err %v\n", groupName, err)
		}
		fmt.Println(group.String())
	},
}

func init() {
	describeCmd.AddCommand(describeGroupCmd)
}
