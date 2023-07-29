package cmd

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

func (that *Cmder) showinfo() {
	command := &cli.Command{
		Name:    "show",
		Aliases: []string{"sho", "sh"},
		Usage:   "Show [gvc] installation path and config file path.",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.ShowPath()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) version() {
	command := &cli.Command{
		Name:    "version",
		Aliases: []string{"ver", "vsi"},
		Usage:   "Show gvc version info.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewSelf()
			v.ShowVersion()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) uninstall() {
	command := &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"unins", "delete", "del"},
		Usage:   "[Caution] Remove gvc and softwares installed by gvc!",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.Uninstall()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vconf() {
	command := &cli.Command{
		Name:        "config",
		Aliases:     []string{"conf", "cnf", "c"},
		Usage:       "Config file management for gvc.",
		Subcommands: []*cli.Command{},
	}
	dav := &cli.Command{
		Name:    "webdav",
		Aliases: []string{"dav", "w"},
		Usage:   "Setup webdav account info.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.SetWebdavAccount()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, dav)

	pull := &cli.Command{
		Name:    "pull",
		Aliases: []string{"pl"},
		Usage:   "Pull settings from remote webdav and apply them to applications.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.FetchAndApplySettings()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, pull)

	push := &cli.Command{
		Name:    "push",
		Aliases: []string{"ph"},
		Usage:   "Gather settings from applications and sync them to remote webdav.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.GatherAndPushSettings()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, push)

	reset := &cli.Command{
		Name:    "reset",
		Aliases: []string{"rs", "r"},
		Usage:   "Reset the gvc config file to default values.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.RestoreDefaultGVConf()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, reset)

	that.Commands = append(that.Commands, command)
}
