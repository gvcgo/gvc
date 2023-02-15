package cmd

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

type Cmder struct {
	*cli.App
}

func New() *Cmder {
	c := &Cmder{
		App: &cli.App{
			Commands: []*cli.Command{},
		},
	}
	c.initiate()
	return c
}

func (that *Cmder) initiate() {
	that.vhost()
	that.vgo()
}

func (that *Cmder) vhost() {
	command := &cli.Command{
		Name:        "host",
		Aliases:     []string{"h"},
		Usage:       "gvc host",
		Description: "Fetch hosts for github.",
		Action: func(ctx *cli.Context) error {
			h := vctrl.New()
			h.Run()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vgo() {
	command := &cli.Command{
		Name:        "go",
		Aliases:     []string{"g"},
		Usage:       "gvc go",
		Description: "Go version control.",
		Subcommands: []*cli.Command{},
	}
	var showall bool
	vremote := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "show-all",
				Aliases:     []string{"a", "all"},
				Usage:       "Show all remote versions.",
				Destination: &showall,
			},
		},
		Usage:       "gvc go r",
		Description: "Show remote versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			arg := vctrl.ShowStable
			if showall {
				arg = vctrl.ShowAll
			}
			gv.ShowRemoteVersions(arg)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vremote)

	vuse := &cli.Command{
		Name:        "use",
		Aliases:     []string{"u"},
		Usage:       "gvc use",
		Description: "Download and use version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewGoVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vlocal := &cli.Command{
		Name:        "local",
		Aliases:     []string{"l"},
		Usage:       "gvc local",
		Description: "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	that.Commands = append(that.Commands, command)
}
