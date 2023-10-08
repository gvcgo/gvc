package cmd

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

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

	gbuild := &cli.Command{
		Name:    "build",
		Aliases: []string{"bui", "b"},
		Usage:   `Compiles go code for multi-platforms [with <-ldflags "-s -w"> builtin].`,
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			args := []string{}
			if len(os.Args) > 3 {
				args = os.Args[3:]
			}
			gv.Build(RecoverArgs(args...)...)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, gbuild)

	grename := &cli.Command{
		Name:    "renameTo",
		Aliases: []string{"rnt", "rto"},
		Usage:   `Rename a local go module[gvc go rto NEW_MODULE_NAME].`,
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			newName := ctx.Args().First()
			moduleDir, _ := os.Getwd()
			if ok, _ := utils.PathIsExist(filepath.Join(moduleDir, "go.mod")); !ok {
				gprint.PrintError("Can not find go.mod in current working dir.")
				return nil
			}
			gv.RenameLocalModule(moduleDir, newName)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, grename)

	gdist := &cli.Command{
		Name:    "list-distributions",
		Aliases: []string{"list-dist", "dist", "ld"},
		Usage:   "List the platforms supported by go compilers.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewGoVersion()
			gv.ShowGoDistlist()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, gdist)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vjava() {
	command := &cli.Command{
		Name:        "java",
		Aliases:     []string{"jdk", "j"},
		Usage:       "Java jdk version management.",
		Subcommands: []*cli.Command{},
	}

	vuse := &cli.Command{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "Download and use jdk.}",
		Action: func(ctx *cli.Context) error {
			version := ctx.Args().First()
			if version != "" {
				gv := vctrl.NewJDKVersion()
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
			gv := vctrl.NewJDKVersion()
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

	fixForWin := &cli.Command{
		Name:    "rmfix",
		Aliases: []string{"rfix"},
		Usage:   "Automatically remove python.exe generated by Windows system in ~/AppData/Local/Microsoft/WindowsApps .",
		Action: func(ctx *cli.Context) error {
			nv := vctrl.NewPyVenv()
			nv.FixSystemGenerationsForWin()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, fixForWin)

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

	iVcpkg := &cli.Command{
		Name:    "install-vcpkg",
		Aliases: []string{"insv", "iv"},
		Usage:   "Install vcpkg.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewCppManager()
			v.InstallVCPkg()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, iVcpkg)

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

	genv := &cli.Command{
		Name:    "SetPathEnv",
		Aliases: []string{"env", "path"},
		Usage:   "Automatically set path env for flutter.",
		Action: func(ctx *cli.Context) error {
			gcode := vctrl.NewFlutterVersion()
			gcode.CheckAndInitEnv()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, genv)

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

	vinstallAndroidTools := &cli.Command{
		Name:    "install-android-sdkmanager",
		Aliases: []string{"install-sdkm", "isdkm", "ism"},
		Usage:   "Install android cmdline tools(sdkmanager).",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.InstallAndroidTool()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vinstallAndroidTools)

	vavdCreate := &cli.Command{
		Name:      "install-build-tools-create-avd",
		Aliases:   []string{"ibt", "cavd"},
		Usage:     "Install build-tools, platform-tools, etc. And create avd for android.",
		ArgsUsage: "Specify a avd name.",
		Action: func(ctx *cli.Context) error {
			avdName := ctx.Args().First()
			gv := vctrl.NewFlutterVersion()
			gv.SetupAVD(avdName)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vavdCreate)

	vavdStart := &cli.Command{
		Name:    "start-avd",
		Aliases: []string{"savd"},
		Usage:   "Start an avd.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.StartAVD()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vavdStart)

	vreplace := &cli.Command{
		Name:    "gradle-repo-aliyun",
		Aliases: []string{"repo", "aliyun"},
		Usage:   "use aliyun repo for android gradle.",
		Action: func(ctx *cli.Context) error {
			gv := vctrl.NewFlutterVersion()
			gv.ReplaceMavenRepo()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vreplace)

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

	installAnalyzer := &cli.Command{
		Name:    "install-analyzer",
		Aliases: []string{"insa", "ia"},
		Usage:   "Install v-analyzer and related extension for vscode.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewVlang()
			v.InstallVAnalyzerForVscode()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installAnalyzer)

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

func (that *Cmder) vprotobuf() {
	command := &cli.Command{
		Name:        "proto",
		Aliases:     []string{"protobuf", "protoc", "pt"},
		Usage:       "Protoc installation.",
		Subcommands: []*cli.Command{},
	}
	var force bool
	install := &cli.Command{
		Name:    "install",
		Aliases: []string{"ins", "i"},
		Usage:   "Install protoc.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "Force to replace old version.",
				Destination: &force,
			},
		},
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewProtobuffer()
			v.Install(force)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, install)

	installGoPlugin := &cli.Command{
		Name:    "install-go-plugin",
		Aliases: []string{"igo", "ig"},
		Usage:   "Install protoc-gen-go.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewProtobuffer()
			v.InstallGoProtobufPlugin()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installGoPlugin)

	installGoGrpcPlugin := &cli.Command{
		Name:    "install-grpc-plugin",
		Aliases: []string{"igrpc", "igr"},
		Usage:   "Install protoc-gen-go-grpc.",
		Action: func(ctx *cli.Context) error {
			v := vctrl.NewProtobuffer()
			v.InstallGoProtoGRPCPlugin()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, installGoGrpcPlugin)
	that.Commands = append(that.Commands, command)
}
