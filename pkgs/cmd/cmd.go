package cmd

import (
	"fmt"

	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/moqsien/gvc/pkgs/vctrl/vproxy"
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
		Usage:   "Show [gvc] install path.",
		Action: func(ctx *cli.Context) error {
			self := vctrl.NewSelf()
			self.ShowInstallPath()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) startXray() {
	var start bool
	command := &cli.Command{
		Name:    "xray",
		Aliases: []string{"ray", "xry", "x"},
		Usage:   "Start Xray Client.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "start",
				Aliases:     []string{"st", "s"},
				Usage:       "Start Xray Client.",
				Destination: &start,
			},
		},
		Action: func(ctx *cli.Context) error {
			xctrl := vproxy.NewXrayCtrl()
			xctrl.DownloadGeoIP()
			if start {
				xctrl.StartXray()
			} else {
				xctrl.StartShell()
			}
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vhost() {
	command := &cli.Command{
		Name:        "host",
		Aliases:     []string{"h", "hosts"},
		Usage:       "Manage system hosts file.",
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
		Usage:       "Go version control.",
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
		Usage:       "VSCode management.",
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
			gcode := vctrl.NewCode()
			gcode.InstallExts()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installexts)

	showexts := &cli.Command{
		Name:    "sync-extensions",
		Aliases: []string{"se", "sext"},
		Usage:   "Push local installed vscode extensions info to remote webdav.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.SyncInstalledExts()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, showexts)

	getsettings := &cli.Command{
		Name:    "get-settings",
		Aliases: []string{"gs", "gset"},
		Usage:   "Get vscode settings(keybindings include) info from remote webdav.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewCode()
			gcode.GetSettings()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, getsettings)

	pushsettings := &cli.Command{
		Name:    "push-settings",
		Aliases: []string{"ps", "pset"},
		Usage:   "Push vscode settings(keybindings include) info to remote webdav.",
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
		Usage:       "GVC config file management.",
		Subcommands: []*cli.Command{},
	}
	dav := &cli.Command{
		Name:    "webdav",
		Aliases: []string{"dav", "w"},
		Usage:   "Setup webdav account info to backup local settings for gvc, vscode, neovim etc.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.SetupWebdav()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, dav)

	pull := &cli.Command{
		Name:    "pull",
		Aliases: []string{"pl"},
		Usage:   "Pull settings to local backup dir from your remote webdav.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.Pull()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, pull)

	push := &cli.Command{
		Name:    "push",
		Aliases: []string{"ph"},
		Usage:   "Push settings from local backup dir to your remote webdav.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.Push()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, push)

	show := &cli.Command{
		Name:    "show",
		Aliases: []string{"sh", "s"},
		Usage:   "Show path to conf files.",
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

	reset := &cli.Command{
		Name:    "reset",
		Aliases: []string{"rs", "r"},
		Usage:   "Reset config file to default values.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.Reset()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, reset)

	downbackfiles := &cli.Command{
		Name:    "download",
		Aliases: []string{"dl", "d"},
		Usage:   "Download example config files from gitee when backup dir is empty.",
		Action: func(ctx *cli.Context) error {
			cnf := confs.New()
			cnf.UseDefautFiles()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, downbackfiles)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vnvim() {
	command := &cli.Command{
		Name:        "nvim",
		Aliases:     []string{"neovim", "nv", "n"},
		Usage:       "GVC neovim management.",
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
		Usage:       "GVC jdk management.",
		Subcommands: []*cli.Command{},
	}
	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use jdk.",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJavaVersion()
				gv.UseVersion(version)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vuse)

	vshow := &cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "Show available versions.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewJavaVersion()
			gv.ShowVersions()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vshow)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vrust() {
	command := &cli.Command{
		Name:        "rust",
		Aliases:     []string{"rustc", "ru", "r"},
		Usage:       "GVC rust management.",
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
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vnodejs() {
	command := &cli.Command{
		Name:        "nodejs",
		Aliases:     []string{"node", "no"},
		Usage:       "Nodejs version control.",
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

func (that *Cmder) vcygwin() {
	command := &cli.Command{
		Name:        "cygwin",
		Aliases:     []string{"cygw", "cyg", "cy"},
		Usage:       "Cygwin management.",
		Subcommands: []*cli.Command{},
	}
	install := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install Cygwin.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewCygwin()
			v.InstallByDefault("")
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, install)

	ipackage := &cli.Command{
		Name:    "package",
		Aliases: []string{"pack", "p"},
		Usage:   "Install packages for Cygwin.",
		Action: func(ctx *cli.Context) error {
			if packs := ctx.Args().First(); packs != "" {
				v := vctrl.NewCygwin()
				v.InstallByDefault(packs)
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, ipackage)
	that.Commands = append(that.Commands, command)
}

func (that *Cmder) initiate() {
	that.uninstall()
	that.showinfo()
	that.vhost()
	that.vgo()
	that.vscode()
	that.vconf()
	that.vnvim()
	that.vjava()
	that.vrust()
	that.vnodejs()
	that.vpython()
	that.vcygwin()
	that.startXray()
}
