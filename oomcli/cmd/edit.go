package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editor = "vi"

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from the default editor.",
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

func openFileByEditor(ctx context.Context, fileName string) error {
	if !checkCommandExist(editor) {
		return fmt.Errorf("%s not found", editor)
	}

	cmd := exec.CommandContext(ctx, editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
