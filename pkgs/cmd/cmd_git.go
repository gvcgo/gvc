package cmd

import (
	"strconv"
	"strings"

	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/urfave/cli/v2"
)

/*
github accelerations
*/
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

	vreverse := &cli.Command{
		Name:    "add-reverse-proxy",
		Aliases: []string{"arproxy", "arp"},
		Usage:   "Add reverse proxy for github.",
		Action: func(ctx *cli.Context) error {
			proxyStr := ctx.Args().First()
			vg := vctrl.NewGhDownloader()
			vg.SetReverseProxyForDownload(proxyStr)
			return nil
		},
	}
	command.Subcommands = append(command.Subcommands, vreverse)

	that.Commands = append(that.Commands, command)
}

func (that *Cmder) vinstallGitWin() {
	command := &cli.Command{
		Name:    "win-git-install",
		Aliases: []string{"wgit", "wgi"},
		Usage:   "Install git for windows.",
		Action: func(ctx *cli.Context) error {
			vg := vctrl.NewGhDownloader()
			vg.InstallGitForWindows()
			return nil
		},
	}
	that.Commands = append(that.Commands, command)
}

/*
git subcommands using proxies
*/
func (that *Cmder) vgit() {
	vg := vctrl.NewGhDownloader()
	var defaultProxy string = vg.ReadDefaultProxy()
	var mannualProxy string
	var disableProxy bool

	gSetProxy := &cli.Command{
		Name:    "git-set-proxy",
		Aliases: []string{"gsproxy", "gsp"},
		Usage:   "Set default proxy for git [default: http://localhost:2023].",
		Action: func(ctx *cli.Context) error {
			vg.SaveDefaultProxy(ctx.Args().First())
			return nil
		},
	}
	that.Commands = append(that.Commands, gSetProxy)

	gclone := &cli.Command{
		Name:      "git-clone",
		Aliases:   []string{"gclone", "gclo"},
		Usage:     "Git Clone using a proxy.",
		ArgsUsage: "specify a repository url like <git@github.com:moqsien/gvc.git>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			projectUrl := ctx.Args().First()
			if projectUrl == "" {
				return nil
			}
			proxyUrl := defaultProxy
			if mannualProxy != "" {
				proxyUrl = mannualProxy
			}
			vg.Clone(projectUrl, proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gclone)

	gpull := &cli.Command{
		Name:    "git-pull",
		Aliases: []string{"gpull", "gpul"},
		Usage:   "Git Pull using a proxy.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			proxyUrl := defaultProxy
			if mannualProxy != "" {
				proxyUrl = mannualProxy
			}
			vg.Pull(proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gpull)

	gpush := &cli.Command{
		Name:    "git-push",
		Aliases: []string{"gpush", "gpus"},
		Usage:   "Git Push using a proxy.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			proxyUrl := defaultProxy
			if mannualProxy != "" {
				proxyUrl = mannualProxy
			}
			vg.Push(proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gpush)

	gcommitPush := &cli.Command{
		Name:      "git-commit-push",
		Aliases:   []string{"gcpush", "gcp"},
		Usage:     "Git commit and push to remote using a proxy.",
		ArgsUsage: `specify commit messages.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
			&cli.BoolFlag{
				Name:        "disable-proxy",
				Aliases:     []string{"d", "dp", "dpxy"},
				Usage:       "Disable proxy usage.",
				Destination: &disableProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			commitMsgList := ctx.Args().Slice()
			commitMsg := "update"
			if len(commitMsg) > 0 {
				commitMsg = strings.Join(commitMsgList, " ")
			}
			var proxyUrl string
			if !disableProxy {
				proxyUrl = defaultProxy
				if mannualProxy != "" {
					proxyUrl = mannualProxy
				}
			}
			vg.CommitAndPush(commitMsg, proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gcommitPush)

	gAddTagPush := &cli.Command{
		Name:      "git-add-tag-push",
		Aliases:   []string{"gaddtag", "gatag", "gat"},
		Usage:     "Git add a new tag and push to remote using a proxy.",
		ArgsUsage: `specify a tag name.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
			&cli.BoolFlag{
				Name:        "disable-proxy",
				Aliases:     []string{"d", "dp", "dpxy"},
				Usage:       "Disable proxy usage.",
				Destination: &disableProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			tag := ctx.Args().First()
			if tag == "" {
				return nil
			}
			var proxyUrl string
			if !disableProxy {
				proxyUrl = defaultProxy
				if mannualProxy != "" {
					proxyUrl = mannualProxy
				}
			}
			vg.AddTagAndPush(tag, proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gAddTagPush)

	gDelTagPush := &cli.Command{
		Name:      "git-del-tag-push",
		Aliases:   []string{"gdeltag", "gdtag", "gdt"},
		Usage:     "Git delete a tag and push to remote using a proxy.",
		ArgsUsage: `specify a tag name.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "proxy",
				Aliases:     []string{"pxy", "p"},
				Usage:       "Specify your proxy.",
				Destination: &mannualProxy,
			},
			&cli.BoolFlag{
				Name:        "disable-proxy",
				Aliases:     []string{"d", "dp", "dpxy"},
				Usage:       "Disable proxy usage.",
				Destination: &disableProxy,
			},
		},
		Action: func(ctx *cli.Context) error {
			tag := ctx.Args().First()
			if tag == "" {
				return nil
			}
			var proxyUrl string
			if !disableProxy {
				proxyUrl = defaultProxy
				if mannualProxy != "" {
					proxyUrl = mannualProxy
				}
			}
			vg.DelTagAndPush(tag, proxyUrl)
			return nil
		},
	}
	that.Commands = append(that.Commands, gDelTagPush)

	gLatestTag := &cli.Command{
		Name:    "git-show-tag-latest",
		Aliases: []string{"gshowtaglatest", "gstag", "gst"},
		Usage:   "Git show the latest tag of a local repository.",
		Action: func(ctx *cli.Context) error {
			vg.ShowLatestTag()
			return nil
		},
	}
	that.Commands = append(that.Commands, gLatestTag)
}
