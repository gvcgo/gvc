package cmd

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

func (that *Cmder) vscode() {
	command := &cli.Command{
		Name:        "vscode",
		Aliases:     []string{"vsc", "vs", "v"},
		Usage:       "VSCode and extensions installation.",
		Subcommands: []*cli.Command{},
	}
	ginstall := &cli.Command{
		Name:    "install",
		Aliases: []string{"i", "ins"},
		Usage:   "Automatically install vscode.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, ginstall)

	installexts := &cli.Command{
		Name:    "install-extensions",
		Aliases: []string{"ie", "iext"},
		Usage:   "Automatically install extensions for vscode.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewGVCWebdav()
			gcode.InstallVSCodeExts("")
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installexts)

	repairgit := &cli.Command{
		Name:    "use-msys2-cygwin-git",
		Aliases: []string{"use-git", "ug"},
		Usage:   "Repair and make use of git.exe from Msys2/Cygwin.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCppManager()
			gcode.RepairGitForVSCode()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, repairgit)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vnvim() {
	command := &cli.Command{
		Name:        "nvim",
		Aliases:     []string{"neovim", "nv", "n"},
		Usage:       "Neovim installation.",
		Subcommands: []*cli.Command{},
	}
	nvims := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install neovim.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewNVim()
			v.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, nvims)
	that.Commands = append(that.Commands, command)
}
