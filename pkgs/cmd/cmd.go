package cmd

import (
	"fmt"

	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

type Cmder struct {
	*cli.App
}

func New() *Cmder {
	c := &Cmder{
		App: &cli.App{
			Usage:       "gvc <Command> <SubCommand>...",
			Description: "A productive tool to manage your development environment.",
			Commands:    []*cli.Command{},
		},
	}
	c.initiate()
	return c
}

func (that *Cmder) vhost() {
	command := &cli.Command{
		Name:        "host",
		Aliases:     []string{"h", "hosts"},
		Usage:       "gvc host",
		Description: "Manage system hosts file.",
		Subcommands: []*cli.Command{},
	}
	fetch := &cli.Command{
		Name:        "fetch",
		Aliases:     []string{"f"},
		Usage:       "gvc host fetch",
		Description: "Fetch github hosts info.",
		Action: func(ctx *cli.Context) error {
			h := vctrl.NewHosts()
			h.Run()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, fetch)

	showpath := &cli.Command{
		Name:        "show",
		Aliases:     []string{"s"},
		Usage:       "gvc host show",
		Description: "Show hosts file path.",
		Action: func(ctx *cli.Context) error {
			h := vctrl.NewHosts()
			h.ShowFilePath()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, showpath)

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

	installexts := &cli.Command{
		Name:        "install-extensions",
		Aliases:     []string{"ie", "iext"},
		Usage:       "gvc vscode install-extensions",
		Description: "Automatically install extensions for vscode.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.InstallExts()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installexts)

	showexts := &cli.Command{
		Name:        "sync-extensions",
		Aliases:     []string{"se", "sext"},
		Usage:       "gvc vscode sync-extensions",
		Description: "Push local installed vscode extensions info to remote webdav.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.SyncInstalledExts()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, showexts)

	getsettings := &cli.Command{
		Name:        "get-settings",
		Aliases:     []string{"gs", "gset"},
		Usage:       "gvc vscode get-settings",
		Description: "Get vscode settings(keybindings include) info from remote webdav.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.GetSettings()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, getsettings)

	pushsettings := &cli.Command{
		Name:        "push-settings",
		Aliases:     []string{"ps", "pset"},
		Usage:       "gvc vscode push-settings",
		Description: "Push vscode settings(keybindings include) info to remote webdav.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.SyncSettings()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, pushsettings)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vconf() {
	command := &cli.Command{
		Name:        "config",
		Aliases:     []string{"conf", "cnf", "c"},
		Usage:       "gvc config <Command>",
		Description: "GVC config file management.",
		Subcommands: []*cli.Command{},
	}
	dav := &cli.Command{
		Name:        "webdav",
		Aliases:     []string{"dav", "w"},
		Usage:       "gvc config webdav",
		Description: "Setup webdav account info to backup settings of vscode and gvc.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.SetupWebdav()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, dav)

	pull := &cli.Command{
		Name:        "pull",
		Aliases:     []string{"pl"},
		Usage:       "gvc config pull",
		Description: "Pull settings from your remote webdav.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.Pull()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, pull)

	push := &cli.Command{
		Name:        "push",
		Aliases:     []string{"ph"},
		Usage:       "gvc config push",
		Description: "Push settings to your remote webdav.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.Push()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, push)

	show := &cli.Command{
		Name:        "show",
		Aliases:     []string{"sh", "s"},
		Usage:       "gvc config pull",
		Description: "Show path to conf files.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			fmt.Println("GVC config file:")
			cnf.ShowPath()
			fmt.Println("WebDAV config file:")
			cnf.ShowDavConfigPath()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, show)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vnvim() {
	command := &cli.Command{
		Name:        "nvim",
		Aliases:     []string{"neovim", "nv", "n"},
		Usage:       "gvc nvim <Command>",
		Description: "GVC neovim management.",
		Subcommands: []*cli.Command{},
	}
	nvims := &cli.Command{
		Name:        "install",
		Aliases:     []string{"ins", "i"},
		Usage:       "gvc nvim install",
		Description: "Install neovim.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewNVim()
			v.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, nvims)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) initiate() {
	that.vhost()
	that.vgo()
	that.vscode()
	that.vconf()
	that.vnvim()
}
