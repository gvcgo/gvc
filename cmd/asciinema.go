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
		Short:   "Creates a record.",
		Long:    "Example: g a record <xxx.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			ascer.Record(args[0])
		},
	}
	parent.AddCommand(record)

	play := &cobra.Command{
		Use:     "play",
		Aliases: []string{"p"},
		Short:   "Plays a record.",
		Long:    "Example: g a p <xxx.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			ascer.Play(args[0])
		},
	}
	parent.AddCommand(play)

	upload := &cobra.Command{
		Use:     "upload",
		Aliases: []string{"u"},
		Short:   "Uploads a record file to asciinema.org.",
		Long:    "Example: g a u <xxx.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			ascer.Upload(args[0])
		},
	}
	parent.AddCommand(upload)

	convert := &cobra.Command{
		Use:     "convert-to-gif",
		Aliases: []string{"cg"},
		Short:   "Converts an asciinema cast to gif.",
		Long:    "Example: g a cg <input.cast> <output.gif>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			ascer.ConvertToGif(args[0], args[1])
		},
	}
	parent.AddCommand(convert)

	cli.rootCmd.AddCommand(parent)
}
