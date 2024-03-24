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

	cut := &cobra.Command{
		Use:     "cut",
		Aliases: []string{"c"},
		Short:   "Removes a certain range of time frames.",
		Long:    "Example: g a c --start=1.00 --end=10.00 <input.cast> <output.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			start, _ := cmd.Flags().GetFloat64("start")
			end, _ := cmd.Flags().GetFloat64("end")
			if len(args) < 2 || end <= 0 || start < 0 {
				cmd.Help()
				return
			}
			ascer.Cut(args[0], args[1], start, end)
		},
	}
	cut.Flags().Float64P("start", "s", 0, "start time in seconds")
	cut.Flags().Float64P("end", "e", 0, "end time in seconds")
	parent.AddCommand(cut)

	speed := &cobra.Command{
		Use:     "speed",
		Aliases: []string{"s"},
		Short:   "Updates the cast speed by a certain factor.",
		Long:    "Example: g a s --factor=1.50 --start=1.00 --end=10.00 <input.cast> <output.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			factor, _ := cmd.Flags().GetFloat64("factor")
			start, _ := cmd.Flags().GetFloat64("start")
			end, _ := cmd.Flags().GetFloat64("end")
			if len(args) < 2 || end <= 0 || start < 0 || factor <= 0 {
				cmd.Help()
				return
			}
			ascer.Speed(args[0], args[1], factor, start, end)
		},
	}
	speed.Flags().Float64P("factor", "f", 0, "speed factor")
	speed.Flags().Float64P("start", "s", 0, "start time in seconds")
	speed.Flags().Float64P("end", "e", 0, "end time in seconds")
	parent.AddCommand(speed)

	quantize := &cobra.Command{
		Use:     "quantize",
		Aliases: []string{"q"},
		Short:   "Updates the cast delays following quantization ranges.",
		Long:    "Example: g a q --range=1.0,2.0 <input.cast> <output.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			ranges, _ := cmd.Flags().GetStringArray("range")
			if len(ranges) == 0 {
				cmd.Help()
				return
			}
			ascer.Quantize(args[0], args[1], ranges)
		},
	}
	quantize.Flags().StringArrayP("range", "r", []string{}, `quantization ranges "min,max" separated by commas`)
	parent.AddCommand(quantize)

	cli.rootCmd.AddCommand(parent)
}
