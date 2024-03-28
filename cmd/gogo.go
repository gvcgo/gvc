package cmd

import (
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/pkg/dev"
	"github.com/spf13/cobra"
)

func RegisterGopher(cli *Cli) {
	parent := &cobra.Command{
		Use:     "gopher",
		Aliases: []string{"go"},
		Short:   "Some useful comand for gophers.",
		GroupID: cli.groupID,
	}

	build := &cobra.Command{
		Use:                "build",
		Aliases:            []string{"b"},
		Short:              `Compiles go code for multi-platforms [with <-ldflags "-s -w"> builtin].`,
		Long:               `If you are planning to use "-X", then remember to replace any "$" by "#".`,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) > 3 {
				args = os.Args[3:]
			} else {
				args = []string{}
			}
			dev.Build(args...)
		},
	}
	parent.AddCommand(build)

	rename := &cobra.Command{
		Use:     "renameto",
		Aliases: []string{"rt"},
		Short:   "Renames a local package to a new name.",
		Long:    "Example: g go rt <your-new-name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			moduleDir, _ := os.Getwd()
			if ok, _ := gutils.PathIsExist(filepath.Join(moduleDir, "go.mod")); !ok {
				gprint.PrintError("Can not find go.mod in current working dir.")
				return
			}
			dev.RenameLocalModule(moduleDir, args[0])
		},
	}
	parent.AddCommand(rename)

	installBinaries := &cobra.Command{}
	parent.AddCommand(installBinaries)

	cli.rootCmd.AddCommand(parent)
}
