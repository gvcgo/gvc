package cmd

import (
	"strconv"

	"github.com/moqsien/gvc/pkgs/utils/sorts"
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

func (that *Cmder) uninstall() {
	command := &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"unins", "delete", "del"},
		Usage:   "[Caution] Delete gvc and softwares installed by gvc!",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.Uninstall()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) showinfo() {
	command := &cli.Command{
		Name:    "show",
		Aliases: []string{"sho", "sh"},
		Usage:   "Show [gvc] installation path and config file path.",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.ShowInstallPath()

			dav := vctrl.NewGVCWebdav()
			dav.ShowConfigPath()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vhost() {
	command := &cli.Command{
		Name:        "host",
		Aliases:     []string{"h", "hosts"},
		Usage:       "Sytem hosts file management(need admistrator or root).",
		Subcommands: []*cli.Command{},
	}
	fetch := &cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Usage:   "Fetch github hosts info.",
		Action: func(ctx *cli.Context) error {
			h := vctrl.NewHosts()
			h.Run(true)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, fetch)

	fetchall := &cli.Command{
		Name:    "fetchall",
		Aliases: []string{"fa"},
		Usage:   "Get all github hosts info with no ping filters.",
		Action: func(ctx *cli.Context) error {
			h := vctrl.NewHosts()
			h.Run()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, fetchall)

	showpath := &cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "Show hosts file path.",
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
		Usage:       "Go version management.",
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
		Usage: "Show remote versions.",
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
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use version.",
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
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	rmunused := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmunused)

	rmversion := &cli.Command{
		Name:    "remove-version",
		Aliases: []string{"rm"},
		Usage:   "Remove a version.",
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
		Name:    "add-envs",
		Aliases: []string{"env", "e", "ae"},
		Usage:   "Add envs for go.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.CheckAndInitEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, genvs)

	var (
		orderByUpdate bool
		libName       string
	)
	vsearch := &cli.Command{
		Name:    "search-package",
		Aliases: []string{"sp", "search"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "package-name",
				Aliases:     []string{"n", "name"},
				Usage:       "Name of the package.",
				Destination: &libName,
			},
			&cli.BoolFlag{
				Name:        "order-by-time",
				Aliases:     []string{"o", "ou"},
				Usage:       "Order by update time.",
				Destination: &orderByUpdate,
			},
		},
		Usage: "Search for third-party packages.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			var orderBy int = sorts.ByImported
			if orderByUpdate {
				orderBy = sorts.ByUpdate
			}
			if libName == "" {
				libName = ctx.Args().First()
			}
			if libName != "" {
				gv.SearchLibs(libName, orderBy)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vsearch)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vscode() {
	command := &cli.Command{
		Name:        "vscode",
		Aliases:     []string{"vsc", "vs", "v"},
		Usage:       "VSCode and extensions installation.",
		Subcommands: []*cli.Command{},
	}
	genvs := &cli.Command{
		Name:    "install",
		Aliases: []string{"i", "ins"},
		Usage:   "Automatically install vscode.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, genvs)

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
			dav.SetAccount()
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

func (that *Cmder) vjava() {
	command := &cli.Command{
		Name:        "java",
		Aliases:     []string{"jdk", "j"},
		Usage:       "Java jdk version management.",
		Subcommands: []*cli.Command{},
	}

	var useInjdk bool
	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use jdk.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "cn",
				Aliases:     []string{"zh", "z"},
				Usage:       "Use injdk.cn as resource url.",
				Destination: &useInjdk,
			},
		},
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJDKVersion()
				gv.IsOfficial = !useInjdk
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show available versions.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "cn",
				Aliases:     []string{"zh", "z"},
				Usage:       "Use injdk.cn as resource url.",
				Destination: &useInjdk,
			},
		},
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJDKVersion()
			gv.IsOfficial = !useInjdk
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJDKVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	vrm := &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an installed version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJDKVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrm)

	vrmall := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"rmu", "ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJDKVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrmall)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vrust() {
	command := &cli.Command{
		Name:        "rust",
		Aliases:     []string{"rustc", "ru", "r"},
		Usage:       "Rust installation.",
		Subcommands: []*cli.Command{},
	}
	iRust := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install the latest rust compiler tools.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewRustInstaller()
			v.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, iRust)
	setEnv := &cli.Command{
		Name:    "setenv",
		Aliases: []string{"env", "se", "e"},
		Usage:   "Set acceleration env for rust.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewRustInstaller()
			v.SetAccelerationEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, setEnv)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vnodejs() {
	command := &cli.Command{
		Name:        "nodejs",
		Aliases:     []string{"node", "no"},
		Usage:       "NodeJS version management.",
		Subcommands: []*cli.Command{},
	}
	vremote := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show remote versions.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewNodeVersion()
			nv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vremote)

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				nv := vctrl.NewNodeVersion()
				nv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewNodeVersion()
			nv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	rmunused := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewNodeVersion()
			nv.RemoveVersion("all")
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmunused)

	rmversion := &cli.Command{
		Name:    "remove-version",
		Aliases: []string{"rm"},
		Usage:   "Remove a version.",
		Action: func(ctx *cli.Context) error {
			if version := ctx.Args().First(); version != "" {
				nv := vctrl.NewNodeVersion()
				nv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmversion)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vpython() {
	command := &cli.Command{
		Name:        "python",
		Aliases:     []string{"py"},
		Usage:       "Python version management.",
		Subcommands: []*cli.Command{},
	}
	vremote := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show remote versions.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewPyVenv()
			nv.ListRemoteVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vremote)

	var useDefault bool
	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use a version.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "accelerate",
				Aliases:     []string{"acc", "a"},
				Usage:       "Use default version[likely 3.11.2] to accelerte installation.",
				Destination: &useDefault,
			},
		},
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				nv := vctrl.NewPyVenv()
				if useDefault {
					nv.InstallVersion(version, true)
				} else {
					nv.InstallVersion(version, false)
				}
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewPyVenv()
			nv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	rmversion := &cli.Command{
		Name:    "remove-version",
		Aliases: []string{"rm"},
		Usage:   "Remove a version.",
		Action: func(ctx *cli.Context) error {
			if version := ctx.Args().First(); version != "" {
				nv := vctrl.NewPyVenv()
				nv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, rmversion)

	updatePyenv := &cli.Command{
		Name:    "update",
		Aliases: []string{"up"},
		Usage:   "Install or update pyenv.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewPyVenv()
			nv.InstallPyenv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, updatePyenv)

	showPath := &cli.Command{
		Name:    "path",
		Aliases: []string{"pth"},
		Usage:   "Show pyenv versions path.",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewPyVenv()
			nv.ShowVersionPath()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, showPath)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vgithub() {
	command := &cli.Command{
		Name:    "github",
		Aliases: []string{"gh"},
		Usage:   "Open github download acceleration websites.",
		Action: func(ctx *cli.Context) error {
			chosenStr := ctx.Args().First()
			chosen, _ := strconv.Atoi(chosenStr)
			vg := vctrl.NewGhDownloader()
			vg.OpenByBrowser(chosen)
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vhomebrew() {
	command := &cli.Command{
		Name:        "homebrew",
		Aliases:     []string{"brew", "hb"},
		Usage:       "Homebrew installation or update.",
		Subcommands: []*cli.Command{},
	}
	install := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install Homebrew.",
		Action: func(ctx *cli.Context) error {
			hb := vctrl.NewHomebrew()
			hb.Install()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, install)

	setEnv := &cli.Command{
		Name:    "setenv",
		Aliases: []string{"env", "se", "e"},
		Usage:   "Set env to accelerate Homebrew in China.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewHomebrew()
			v.SetEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, setEnv)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vlang() {
	command := &cli.Command{
		Name:        "vlang",
		Aliases:     []string{"vl"},
		Usage:       "Vlang installation.",
		Subcommands: []*cli.Command{},
	}
	var force bool
	install := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install Vlang.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "Force to replace old version.",
				Destination: &force,
			},
		},
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewVlang()
			v.Install(force)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, install)

	setEnv := &cli.Command{
		Name:    "setenv",
		Aliases: []string{"env", "se", "e"},
		Usage:   "Set env for Vlang.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewVlang()
			v.CheckAndInitEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, setEnv)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vgradle() {
	command := &cli.Command{
		Name:        "gradle",
		Aliases:     []string{"gra", "gr"},
		Usage:       "Gradle version management.",
		Subcommands: []*cli.Command{},
	}

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use gradle.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewGradleVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show available versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGradleVersion()
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGradleVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	vset := &cli.Command{
		Name:    "set",
		Aliases: []string{"se"},
		Usage:   "Set aliyun repository.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGradleVersion()
			gv.GenInitFile()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vset)

	vrm := &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an installed version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewGradleVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrm)

	vrmall := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"rmu", "ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGradleVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrmall)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vmaven() {
	command := &cli.Command{
		Name:        "maven",
		Aliases:     []string{"mav", "ma"},
		Usage:       "Maven version management.",
		Subcommands: []*cli.Command{},
	}

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use maven.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewMavenVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show available versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewMavenVersion()
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewMavenVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	vset := &cli.Command{
		Name:    "set",
		Aliases: []string{"se"},
		Usage:   "Set mirrors and local repository path.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewMavenVersion()
			gv.GenSettingsFile()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vset)

	vrm := &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an installed version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewMavenVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrm)

	vrmall := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"rmu", "ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewMavenVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrmall)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vflutter() {
	command := &cli.Command{
		Name:        "flutter",
		Aliases:     []string{"flu", "fl"},
		Usage:       "Flutter version management.",
		Subcommands: []*cli.Command{},
	}

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use flutter.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewFlutterVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show available versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	vrm := &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an installed version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewFlutterVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrm)

	vrmall := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"rmu", "ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrmall)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vjulia() {
	command := &cli.Command{
		Name:        "julia",
		Aliases:     []string{"jul", "ju"},
		Usage:       "Julia version management.",
		Subcommands: []*cli.Command{},
	}

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use julia.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJuliaVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "remote",
		Aliases: []string{"r"},
		Usage:   "Show available versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJuliaVersion()
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vlocal := &cli.Command{
		Name:    "local",
		Aliases: []string{"l"},
		Usage:   "Show installed versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJuliaVersion()
			gv.ShowInstalled()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vlocal)

	vrm := &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove an installed version.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJuliaVersion()
				gv.RemoveVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrm)

	vrmall := &cli.Command{
		Name:    "remove-unused",
		Aliases: []string{"rmu", "ru"},
		Usage:   "Remove unused versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJuliaVersion()
			gv.RemoveUnused()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrmall)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vtypst() {
	command := &cli.Command{
		Name:        "typst",
		Aliases:     []string{"ty"},
		Usage:       "Typst installation.",
		Subcommands: []*cli.Command{},
	}
	var force bool
	install := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install Typst.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "Force to replace old version.",
				Destination: &force,
			},
		},
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewTypstVersion()
			v.Install(force)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, install)

	setEnv := &cli.Command{
		Name:    "setenv",
		Aliases: []string{"env", "se", "e"},
		Usage:   "Set env for Typst.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewTypstVersion()
			v.CheckAndInitEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, setEnv)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vxtray() {
	commands := &cli.Command{
		Name:    "xtray-shell",
		Aliases: []string{"xshell", "xs", "x"},
		Usage:   "Start an xtray shell.",
		Action: func(ctx *cli.Context) error {
			xe := vctrl.NewXtrayExa()
			xe.Runner.CtrlShell()
			return nil
		},
	}
	that.Commands = append(that.Commands, commands)

	commandr := &cli.Command{
		Name:    vctrl.XtrayStarterCmd,
		Aliases: []string{"xrunner", "xr"},
		Usage:   "Start an xtray client.",
		Action: func(ctx *cli.Context) error {
			xe := vctrl.NewXtrayExa()
			xe.Runner.Start()
			return nil
		},
	}
	that.Commands = append(that.Commands, commandr)

	commandk := &cli.Command{
		Name:    vctrl.XtrayKeeperCmd,
		Aliases: []string{"xkeeper", "xk"},
		Usage:   "Start an xtray keeper.",
		Action: func(ctx *cli.Context) error {
			xe := vctrl.NewXtrayExa()
			xe.Keeper.Run()
			return nil
		},
	}
	that.Commands = append(that.Commands, commandk)
}

func (that *Cmder) vbrowser() {
	command := &cli.Command{
		Name:        "browser",
		Aliases:     []string{"br"},
		Usage:       "Browser data management.",
		Subcommands: []*cli.Command{},
	}

	vshow := &cli.Command{
		Name:    "show-info",
		Aliases: []string{"show", "sh"},
		Usage:   "Show supported browsers and data restore dir.",
		Action: func(ctx *cli.Context) error {
			b := vctrl.NewBrowser()
			b.ShowSupportedBrowser()
			b.ShowBackupPath()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)

	vpush := &cli.Command{
		Name:      "push",
		Aliases:   []string{"psh", "pu"},
		Usage:     "Push browser Bookmarks/Password/ExtensionInfo to webdav.",
		ArgsUsage: "gvc browser push xxx",
		Action: func(ctx *cli.Context) error {
			browserName := ctx.Args().First()
			if browserName != "" {
				b := vctrl.NewBrowser()
				b.Save(browserName, true)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vpush)

	vsave := &cli.Command{
		Name:      "save",
		Aliases:   []string{"sa", "s"},
		Usage:     "Save browser Bookmarks/Password/ExtensionInfo to local dir.",
		ArgsUsage: "gvc browser save xxx",
		Action: func(ctx *cli.Context) error {
			browserName := ctx.Args().First()
			if browserName != "" {
				b := vctrl.NewBrowser()
				b.Save(browserName, false)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vsave)

	vpull := &cli.Command{
		Name:    "pull",
		Aliases: []string{"pul", "pl"},
		Usage:   "Pull browser data from webdav to local dir.",
		Action: func(ctx *cli.Context) error {
			b := vctrl.NewBrowser()
			b.PullData()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vpull)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vcpp() {
	command := &cli.Command{
		Name:        "cpp",
		Usage:       "C/C++ management.",
		Subcommands: []*cli.Command{},
	}
	iMsys2 := &cli.Command{
		Name:    "install-msys2",
		Aliases: []string{"insm", "im"},
		Usage:   "Install the latest msys2.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewCppManager()
			v.InstallMsys2()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, iMsys2)

	uMsys2 := &cli.Command{
		Name:    "uninstall-msys2",
		Aliases: []string{"unim", "um", "remove", "rm"},
		Usage:   "Uninstall msys2.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewCppManager()
			v.UninstallMsys2()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, uMsys2)

	iCygwin := &cli.Command{
		Name:    "install-cygwin",
		Aliases: []string{"insc", "ic"},
		Usage:   "Install Cygwin.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewCppManager()
			v.InstallCygwin("")
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, iCygwin)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) initiate() {
	that.vgo()
	that.vpython()
	that.vjava()
	that.vmaven()
	that.vgradle()
	that.vnodejs()
	that.vflutter()
	that.vjulia()
	that.vrust()
	that.vcpp()
	that.vtypst()
	that.vlang()

	that.vscode()
	that.vnvim()
	that.vxtray()
	that.vbrowser()
	that.vhomebrew()
	that.vhost()
	that.vgithub()

	that.vconf()
	that.showinfo()
	that.uninstall()
}
