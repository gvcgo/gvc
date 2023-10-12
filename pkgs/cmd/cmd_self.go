package cmd

import (
	"fmt"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
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

func (that *Cmder) SetVersionInfo(gitTag, gitHash, gitTime string) {
	that.gitHash = gitHash
	that.gitTag = gitTag
	that.gitTime = gitTime
}

func (that *Cmder) version() {
	command := &cli.Command{
		Name:    "version",
		Aliases: []string{"ver", "vsi"},
		Usage:   "Show gvc version info.",
		Action: func(ctx *cli.Context) error {
			hashTail := that.gitHash
			if len(hashTail) > 8 {
				hashTail = hashTail[len(hashTail)-8:]
			}
			content := fmt.Sprintf(
				"Name: %s\nVersion: %s\nUpdateAt: %s\nHomepage: %s\nEmail: %s",
				"GVC",
				fmt.Sprintf("%s(%s)", that.gitTag, hashTail),
				that.gitTime,
				"https://github.com/moqsien/gvc",
				"moqsien2022@gmail.com",
			)
			gprint.PrintlnByDefault(content)
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) checkUpdate() {
	command := &cli.Command{
		Name:    "check",
		Aliases: []string{"checklatest", "checkupdate"},
		Usage:   "Check and download the latest version of gvc.",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.CheckLatestVersion(that.gitTag)
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

func (that *Cmder) vsshFiles() {
	command := &cli.Command{
		Name:        "ssh-files",
		Aliases:     []string{"sshf", "ssh"},
		Usage:       "Backup your ssh files.",
		Subcommands: []*cli.Command{},
	}
	sshSave := &cli.Command{
		Name:    "save",
		Aliases: []string{"sv", "s"},
		Usage:   "Save your local ssh files to WebDAV.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.GatherSSHFiles()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, sshSave)

	sshDeploy := &cli.Command{
		Name:    "deploy",
		Aliases: []string{"dp", "d"},
		Usage:   "Get ssh files from WebDAV and deploy them to local dir.",
		Action: func(ctx *cli.Context) error {
			dav := vctrl.NewGVCWebdav()
			dav.DeploySSHFiles()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, sshDeploy)
	that.Commands = append(that.Commands, command)
}
