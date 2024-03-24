package cmd

import (
	"github.com/gvcgo/gvc/pkg/asciinema"
	"github.com/spf13/cobra"
)

func RegisterAsciinema(cli *Cli) {
	parent := &cobra.Command{
		Use:     "asciinema",
		Aliases: []string{"asc", "a"},
		GroupID: cli.groupID,
		Short:   "Record your terminal in asciinema cast form.",
	}
	ascer := asciinema.NewAsciinema()

	auth := &cobra.Command{
		Use:     "auth",
		Aliases: []string{"a"},
		Short:   "Authrization to asciinema.org.",
		Run: func(cmd *cobra.Command, args []string) {
			ascer.Auth()
		},
	}
	parent.AddCommand(auth)

	record := &cobra.Command{
		Use:     "record",
		Aliases: []string{"r"},
		Short:   "Create a record.",
		Long:    "g a record <xxx.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			ascer.Record(args[0])
		},
	}
	parent.AddCommand(record)

	cli.rootCmd.AddCommand(parent)
}
