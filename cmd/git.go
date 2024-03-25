package cmd

import (
	"github.com/gvcgo/gvc/pkg/git"
	"github.com/spf13/cobra"
)

func RegisterGit(cli *Cli) {
	parent := &cobra.Command{
		Use:     "git",
		Aliases: []string{"g"},
		Short:   "Git related CLIs.",
		GroupID: cli.groupID,
	}

	hosts := &cobra.Command{
		Use:     "hosts",
		Aliases: []string{"h"},
		Short:   "Updates hosts file.",
		Run: func(cmd *cobra.Command, args []string) {
			m := git.NewModifier()
			m.Run()
		},
	}
	parent.AddCommand(hosts)

	cli.rootCmd.AddCommand(parent)
}
