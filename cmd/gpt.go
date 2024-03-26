package cmd

import (
	"github.com/gvcgo/gvc/pkg/gpt"
	"github.com/spf13/cobra"
)

func RegisterGPT(cli *Cli) {
	parent := &cobra.Command{
		Use:     "gpt",
		Aliases: []string{"G"},
		GroupID: cli.groupID,
		Short:   "ChatGPT or FlyTek spark bot.",
		Run: func(cmd *cobra.Command, args []string) {
			g := gpt.NewGPT()
			g.Run()
		},
	}
	cli.rootCmd.AddCommand(parent)
}
