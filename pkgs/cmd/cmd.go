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
	that.vscode()
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
		Usage:       "gvc go <Command>",
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
		Usage:       "gvc go use",
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
		Usage:       "gvc go local",
		Description: "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	rmunused := &cli.Command{
		Name:        "remove-unused",
		Aliases:     []string{"ru"},
		Usage:       "gvc go ru",
		Description: "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmunused)

	rmversion := &cli.Command{
		Name:        "remove-version",
		Aliases:     []string{"rm"},
		Usage:       "gvc go rm",
		Description: "Remove a version.",
		Action: func(ctx *cli.Context) error {
			if version := ctx.Args().First(); version != "" {
				gv := vctrl.NewGoVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmversion)

	genvs := &cli.Command{
		Name:        "add-envs",
		Aliases:     []string{"env", "e", "ae"},
		Usage:       "gvc go env",
		Description: "Add envs for go.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.CheckAndInitEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, genvs)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vscode() {
	command := &cli.Command{
		Name:        "vscode",
		Aliases:     []string{"vsc", "vs", "v"},
		Usage:       "gvc vscode <Command>",
		Description: "VSCode management.",
		Subcommands: []*cli.Command{},
	}
	genvs := &cli.Command{
		Name:        "install",
		Aliases:     []string{"i", "ins"},
		Usage:       "gvc vscode install",
		Description: "Automatically install vscode.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, genvs)
	that.Commands = append(that.Commands, command)
}
