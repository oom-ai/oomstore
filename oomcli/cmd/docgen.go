package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var docgenCmd = &cobra.Command{
	Use:   "docgen",
	Short: "Generate oomcli docs",
	Run: func(cmd *cobra.Command, args []string) {
		genDoc()
	},
}

func init() {
	// So that it doesn't docgen 'completion' command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Hide the command for non-dev users
	docgenCmd.Hidden = true

	rootCmd.AddCommand(docgenCmd)
}

func genDoc() {
	var sb strings.Builder
	sb.WriteString("# CLI\n\n")
	genDocTree(&sb, rootCmd)
	fmt.Println(sb.String())
}

// Recursively generate markdown documentation for the given command
func genDocTree(sb *strings.Builder, cmd *cobra.Command) {
	if !cmd.IsAvailableCommand() || cmd.IsAdditionalHelpTopicCommand() {
		return
	}
	out := new(bytes.Buffer)
	err := genMarkdown(cmd, out)
	if err != nil {
		exit(err)
	}
	sb.WriteString(out.String())

	for _, child := range cmd.Commands() {
		genDocTree(sb, child)
	}
}

// The method is mostly a copy-and-paste of GenMarkdownCustom from cobra doc package
// The only difference is that it does not print "SEE ALSO" section
// which is not very useful since our doc is single-page
func genMarkdown(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```bash\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	printOptions(buf, cmd, name)

	_, err := buf.WriteTo(w)
	return err
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("Options\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
}
