package cmd

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cobra"
)

var registerEntityOpt types.CreateEntityOpt

var registerEntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "register a new entity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		onestore, err := onestore.Open(ctx, oneStoreOpt)
		if err != nil {
			return err
		}
		_, err = onestore.CreateEntity(ctx, registerEntityOpt)
		return err
	},
}

func init() {
	registerCmd.AddCommand(registerEntityCmd)

	flags := registerEntityCmd.Flags()

	flags.IntVarP(&registerEntityOpt.Length, "length", "l", 0, "entity value length")
	_ = registerEntityCmd.MarkFlagRequired("length")

	flags.StringVarP(&registerEntityOpt.Description, "description", "d", "", "entity description")
}
