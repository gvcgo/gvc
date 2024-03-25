package cmd

import (
	"github.com/gvcgo/gvc/pkg/cloc"
	"github.com/spf13/cobra"
)

type CCtx struct {
	cmd  *cobra.Command
	args []string
}

func (c *CCtx) String(name string) string {
	r, _ := c.cmd.Flags().GetString(name)
	return r
}

func (c *CCtx) Bool(name string) bool {
	r, _ := c.cmd.Flags().GetBool(name)
	return r
}

func (c *CCtx) Args() []string {
	return c.args
}

func RegisterCloc(cli *Cli) {
	parent := &cobra.Command{
		Use:     "cloc",
		Aliases: []string{"cl"},
		Short:   "Counts lines of code.",
		Long:    "Example: cloc <your_path>",
		GroupID: cli.groupID,
		Run: func(cmd *cobra.Command, args []string) {
			cloc := cloc.NewCloc(&CCtx{cmd: cmd})
			cloc.Run()
		},
	}

	parent.Flags().BoolP(cloc.FlagByFile, "f", false, "Report results for every encountered source file.")
	parent.Flags().BoolP(cloc.FlagDebug, "b", false, "Dump debug log for developer.")
	parent.Flags().BoolP(cloc.FlagSkipDuplicated, "s", false, "Skip duplicated files.")
	parent.Flags().BoolP(cloc.FlagShowLang, "l", false, "Print about all languages and extensions.")
	parent.Flags().StringP(cloc.FlagSortTag, "t", "name", `Sort based on a certain column["name", "files", "blank", "comment", "code"].`)
	parent.Flags().StringP(cloc.FlagOutputType, "o", "default", "Show summary only.")
	parent.Flags().StringP(cloc.FlagExcludeExt, "e", "", "Exclude file name extensions (separated commas).")
	parent.Flags().StringP(cloc.FlagIncludeLang, "L", "", "Include language name (separated commas).")
	parent.Flags().StringP(cloc.FlagMatch, "m", "", "Include file name (regex).")
	parent.Flags().StringP(cloc.FlagNotMatch, "M", "", "Exclude file name (regex).")
	parent.Flags().StringP(cloc.FlagMatchDir, "d", "", "Include dir name (regex).")
	parent.Flags().StringP(cloc.FlagNotMatchDir, "D", "", "Exclude dir name (regex).")

	cli.rootCmd.AddCommand(parent)
}
