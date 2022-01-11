package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cobra"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types/apply"
)

var editor = "vi"

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from the default editor",
}

func init() {
	rootCmd.AddCommand(editCmd)

	if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}
}

func getTempFile() (*os.File, error) {
	file, err := os.CreateTemp("", "oomstore.*.yml")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func checkCommandExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func editFileByExternalEditor(ctx context.Context, fileName string) error {
	if !checkCommandExist(editor) {
		return errdefs.Errorf("%s not found", editor)
	}

	cmd := exec.CommandContext(ctx, editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func edit(ctx context.Context, oomStore *oomstore.OomStore, fileName string) error {
	if err := editFileByExternalEditor(ctx, fileName); err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	if err := oomStore.Apply(ctx, apply.ApplyOpt{R: file}); err != nil {
		return fmt.Errorf("apply failed: %v", err)
	}
	return nil
}
