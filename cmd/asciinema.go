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
		Short:   "Records your terminal in asciinema cast form.",
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

	cut := &cobra.Command{
		Use:     "cut",
		Aliases: []string{"c"},
		Short:   "Removes a certain range of time frames.",
		Long:    "Example: g a c --start=1.0 --end=5.0 <in.cast> <out.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			start, _ := cmd.Flags().GetFloat64("start")
			end, _ := cmd.Flags().GetFloat64("end")
			if len(args) < 2 || end <= start {
				cmd.Help()
				return
			}
			ascer.Cut(args[0], args[1], start, end)
		},
	}
	cut.Flags().Float64P("start", "s", 0, "start time")
	cut.Flags().Float64P("end", "e", 0, "end time")
	parent.AddCommand(cut)

	speed := &cobra.Command{
		Use:     "speed",
		Aliases: []string{"s"},
		Short:   "Updates the cast speed by a certain factor.",
		Long:    "Example: g a s --factor=0.7 --start=1.0 --end=5.0 <in.cast> <out.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			factor, _ := cmd.Flags().GetFloat64("factor")
			start, _ := cmd.Flags().GetFloat64("start")
			end, _ := cmd.Flags().GetFloat64("end")
			if len(args) < 2 || end <= start || factor <= 0 {
				cmd.Help()
				return
			}
			ascer.Speed(args[0], args[1], factor, start, end)
		},
	}
	speed.Flags().Float64P("factor", "f", 0.7, "speed factor")
	speed.Flags().Float64P("start", "s", 0, "start time")
	speed.Flags().Float64P("end", "e", 0, "end time")
	parent.AddCommand(speed)

	quantize := &cobra.Command{
		Use:     "quantize",
		Aliases: []string{"q"},
		Short:   "Updates the cast delays following quantization ranges.",
		Long:    "Example: g a q --ranges=1.0,5.0 <in.cast> <out.cast>",
		Run: func(cmd *cobra.Command, args []string) {
			ranges, _ := cmd.Flags().GetStringArray("ranges")
			if len(ranges) == 0 || len(args) < 2 {
				cmd.Help()
				return
			}
			ascer.Quantize(args[0], args[1], ranges)
		},
	}
	quantize.Flags().StringArrayP("ranges", "r", []string{}, "quantization ranges")
	parent.AddCommand(quantize)

	cli.rootCmd.AddCommand(parent)
}
