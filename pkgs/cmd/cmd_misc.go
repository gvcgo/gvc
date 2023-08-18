package cmd

import (
	"runtime"
	"strconv"

	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

func (that *Cmder) vhost() {
	command := &cli.Command{
		Name:        vctrl.HostsCmd,
		Aliases:     []string{"h", "host"},
		Usage:       "Sytem hosts file management(need admistrator or root).",
		Subcommands: []*cli.Command{},
	}

	var previlege bool
	fetch := &cli.Command{
		Name:    vctrl.HostsFileFetchCmd,
		Aliases: []string{"f"},
		Usage:   "Fetch github hosts info.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        vctrl.HostsFlagName,
				Aliases:     []string{"pre", "p"},
				Usage:       "Use admin previlege for windows.",
				Destination: &previlege,
			},
		},
		Action: func(ctx *cli.Context) error {
			h := vctrl.NewHosts()
			if runtime.GOOS != utils.Windows || previlege {
				h.Run()
			} else {
				h.WinRunAsAdmin()
			}
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, fetch)

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

func (that *Cmder) vgithub() {
	command := &cli.Command{
		Name:        "github",
		Aliases:     []string{"gh"},
		Usage:       "Github download speedup.",
		Subcommands: []*cli.Command{},
	}

	var isSourceCode bool
	vdownload := &cli.Command{
		Name:    "download",
		Aliases: []string{"dl", "d"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "code",
				Aliases:     []string{"co", "c"},
				Usage:       "Download only source code.",
				Destination: &isSourceCode,
			},
		},
		Usage: "Download files from github project.",
		Action: func(ctx *cli.Context) error {
			githubProjectUrl := ctx.Args().First()
			vg := vctrl.NewGhDownloader()
			vg.Download(githubProjectUrl, isSourceCode)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vdownload)

	vopen := &cli.Command{
		Name:    "openbrowser",
		Aliases: []string{"open", "ob"},
		Usage:   "Open acceleration website in browser.",
		Action: func(ctx *cli.Context) error {
			chosenStr := ctx.Args().First()
			chosen, _ := strconv.Atoi(chosenStr)
			vg := vctrl.NewGhDownloader()
			vg.OpenByBrowser(chosen)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vopen)

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

func (that *Cmder) vcloc() {
	command := &cli.Command{
		Name:    "cloc",
		Aliases: []string{"cl"},
		Usage:   "Count lines of code.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    vctrl.FlagByFile,
				Aliases: []string{"bf"},
				Value:   false,
				Usage:   "Report results for every encountered source file.",
			},
			&cli.BoolFlag{
				Name:    vctrl.FlagDebug,
				Aliases: []string{"de", "d"},
				Value:   false,
				Usage:   "Dump debug log for developer.",
			},
			&cli.BoolFlag{
				Name:    vctrl.FlagSkipDuplicated,
				Aliases: []string{"skipdup", "sd"},
				Value:   false,
				Usage:   "Skip duplicated files.",
			},
			&cli.BoolFlag{
				Name:    vctrl.FlagShowLang,
				Aliases: []string{"shlang", "sl"},
				Value:   false,
				Usage:   "Print about all languages and extensions.",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagSortTag,
				Aliases: []string{"sort", "st"},
				Value:   "name",
				Usage:   `Sort based on a certain column["name", "files", "blank", "comment", "code"].`,
			},
			&cli.StringFlag{
				Name:    vctrl.FlagOutputType,
				Aliases: []string{"output", "ot"},
				Value:   "default",
				Usage:   "Output type [values: default,cloc-xml,sloccount,json].",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagExcludeExt,
				Aliases: []string{"excl", "ee"},
				Usage:   "Exclude file name extensions (separated commas).",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagIncludeLang,
				Aliases: []string{"langs", "il"},
				Usage:   "Include language name (separated commas).",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagMatch,
				Aliases: []string{"mat", "m"},
				Usage:   "Include file name (regex).",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagNotMatch,
				Aliases: []string{"nmat", "nm"},
				Usage:   "Exclude file name (regex).",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagMatchDir,
				Aliases: []string{"matd", "md"},
				Usage:   "Include dir name (regex).",
			},
			&cli.StringFlag{
				Name:    vctrl.FlagNotMatchDir,
				Aliases: []string{"nmatd", "nmd"},
				Usage:   "Exclude dir name (regex).",
			},
		},
		Action: func(ctx *cli.Context) error {
			cloc := vctrl.NewCloc(ctx)
			cloc.Run()
			return nil
		},
	}

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) asciinema() {
	command := &cli.Command{
		Name:        "asciinema",
		Aliases:     []string{"ascii", "asc"},
		Usage:       "Asciinema terminal recorder.",
		Subcommands: []*cli.Command{},
	}
	vrec := &cli.Command{
		Name:      "record",
		Aliases:   []string{"rec", "r"},
		Usage:     "Record terminal operations.",
		ArgsUsage: "gvc asciinema record xxx",
		Action: func(ctx *cli.Context) error {
			a := vctrl.NewAsciiCast()
			a.Rec(ctx.Args().First())
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vrec)

	vplay := &cli.Command{
		Name:      "play",
		Aliases:   []string{"pl", "p"},
		Usage:     "Play local asciinema file.",
		ArgsUsage: "gvc asciinema play xxx",
		Action: func(ctx *cli.Context) error {
			a := vctrl.NewAsciiCast()
			a.Play(ctx.Args().First())
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vplay)

	vauth := &cli.Command{
		Name:    "auth",
		Aliases: []string{"au", "a"},
		Usage:   "Bind local install-id to your asciinem.org account.",
		Action: func(ctx *cli.Context) error {
			a := vctrl.NewAsciiCast()
			a.Auth()
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vauth)

	vup := &cli.Command{
		Name:      "upload",
		Aliases:   []string{"up", "u"},
		Usage:     "Upload local asciinema file to asciinema.org.",
		ArgsUsage: "gvc asciinema upload xxx",
		Action: func(ctx *cli.Context) error {
			a := vctrl.NewAsciiCast()
			a.Upload(ctx.Args().First())
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vup)

	that.Commands = append(that.Commands, command)
}
