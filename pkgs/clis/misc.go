package clis

import (
	"runtime"

	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

// Homebrew accelerations in China.
func (that *Cli) homebrew() {
	brewCmd := &cobra.Command{
		Use:     "brew",
		Aliases: []string{"B"},
		Short:   "Homebrew accelerations in China.",
		GroupID: that.groupID,
	}

	brewCmd.AddCommand(&cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs homebrew.",
		Run: func(cmd *cobra.Command, args []string) {
			hb := vctrl.NewHomebrew()
			hb.Install()
		},
	})

	brewCmd.AddCommand(&cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Set envs to accelerate homebrew in China.",
		Run: func(cmd *cobra.Command, args []string) {
			v := vctrl.NewHomebrew()
			v.SetEnv()
		},
	})

	that.rootCmd.AddCommand(brewCmd)
}

// Installs sudo command for windows.
func (that *Cli) gsudo() {
	gsudoCmd := &cobra.Command{
		Use:     "gsudo-install",
		Aliases: []string{"gs"},
		Short:   "Installs gsudo for windows.",
		GroupID: that.groupID,
		Run: func(cmd *cobra.Command, args []string) {
			if runtime.GOOS != utils.Windows {
				cmd.Help()
				return
			}
			gs := vctrl.NewGSudo()
			gs.Install(true)
		},
	}
	that.rootCmd.AddCommand(gsudoCmd)
}

func (that *Cli) browser() {
	browserCmd := &cobra.Command{
		Use:     "browser",
		Aliases: []string{"b"},
		Short:   "Manages data for different browsers.",
		GroupID: that.groupID,
	}

	browserCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Lists the supported browsers and dir to save your data.",
		Run: func(cmd *cobra.Command, args []string) {
			b := vctrl.NewBrowser()
			b.ShowSupportedBrowser()
			b.ShowBackupPath()
		},
	})

	browserCmd.AddCommand(&cobra.Command{
		Use:     "upload",
		Aliases: []string{"u"},
		Short:   "Uploads your browser data to remote repo.",
		Long:    `Example: b u <browser_name: To see browser name list, please use command "b list">.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			b := vctrl.NewBrowser()
			b.HandleBrowserFiles(args[0], false)
		},
	})

	browserCmd.AddCommand(&cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Short:   "Downloads browser data from remote repo.",
		Long:    "Example: b d.",
		Run: func(cmd *cobra.Command, args []string) {
			b := vctrl.NewBrowser()
			b.HandleBrowserFiles("", true)
		},
	})

	that.rootCmd.AddCommand(browserCmd)
}

func (that *Cli) docker() {
	dockerCmd := &cobra.Command{
		Use:     "docker",
		Aliases: []string{"d"},
		Short:   "Docker related CLIs.",
		GroupID: that.groupID,
	}

	dockerCmd.AddCommand(&cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs docker.",
		Run: func(cmd *cobra.Command, args []string) {
			dv := vctrl.NewVDocker()
			dv.Install()
		},
	})

	dockerCmd.AddCommand(&cobra.Command{
		Use:     "mirrors",
		Aliases: []string{"m"},
		Short:   "Shows mirrors available in China(for accelerations).",
		Run: func(cmd *cobra.Command, args []string) {
			dv := vctrl.NewVDocker()
			dv.ShowRegistryMirrorInChina()
		},
	})

	that.rootCmd.AddCommand(dockerCmd)
}

func (that *Cli) asciinema() {
	asnemaCmd := &cobra.Command{
		Use:     "asciinema",
		Aliases: []string{"asc", "a"},
		Short:   "Asciinema related CLIs(terminal recorder).",
		GroupID: that.groupID,
	}

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "record",
		Aliases: []string{"r"},
		Short:   "Records your terminal operations.",
		Long:    "Example: a r <your_file_name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			a := vctrl.NewAsciiCast()
			a.Rec(args[0])
		},
	})

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "paly",
		Aliases: []string{"p"},
		Short:   "Plays an asciinema file.",
		Long:    "Example: a p <your_file_path>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			a := vctrl.NewAsciiCast()
			a.Play(args[0])
		},
	})

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "auth",
		Aliases: []string{"a"},
		Short:   "Binds your installation id to asciinema.org account.",
		Run: func(cmd *cobra.Command, args []string) {
			a := vctrl.NewAsciiCast()
			a.Auth()
		},
	})

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "upload",
		Aliases: []string{"u"},
		Short:   "Uploads an asciinema file to asciinema.org.",
		Long:    "Example: a u <your_file_path>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			a := vctrl.NewAsciiCast()
			a.Upload(args[0])
		},
	})

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "upload-conf",
		Aliases: []string{"U"},
		Short:   "Uploads asciinema config file to remote repo.",
		Long:    "Example: a U",
		Run: func(cmd *cobra.Command, args []string) {
			a := vctrl.NewAsciiCast()
			a.HandleAsciinemaConf(false)
		},
	})

	asnemaCmd.AddCommand(&cobra.Command{
		Use:     "download-conf",
		Aliases: []string{"d"},
		Short:   "Downloads asciinema config file from remote repo.",
		Long:    "Example: a d",
		Run: func(cmd *cobra.Command, args []string) {
			a := vctrl.NewAsciiCast()
			a.HandleAsciinemaConf(true)
		},
	})

	that.rootCmd.AddCommand(asnemaCmd)
}

func (that *Cli) gpt() {
	gptCmd := &cobra.Command{
		Use:     "gpt",
		Aliases: []string{"G"},
		Short:   "Starts the ChatGPT/Spark bot.",
		GroupID: that.groupID,
	}

	gptCmd.AddCommand(&cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Starts the ChatGPT/Spark bot.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewVGPT()
			gv.Run()
		},
	})

	gptCmd.AddCommand(&cobra.Command{
		Use:     "upload-conf",
		Aliases: []string{"u"},
		Short:   "Uploads config files to remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewVGPT()
			gv.HandleGPTConf(false)
		},
	})

	gptCmd.AddCommand(&cobra.Command{
		Use:     "download-conf",
		Aliases: []string{"d"},
		Short:   "Downloads config files from remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			gv := vctrl.NewVGPT()
			gv.HandleGPTConf(true)
		},
	})

	that.rootCmd.AddCommand(gptCmd)
}

type CCtx struct {
	cmd  *cobra.Command
	args []string
}

func (c *CCtx) String(name string) string {
	r, _ := c.cmd.Flags().GetString(name)
	return r
}

func (c *CCtx) Bool(name string) bool {
	r, _ := c.cmd.Flags().GetBool(name)
	return r
}

func (c *CCtx) Args() []string {
	return c.args
}

func (that *Cli) cloc() {
	clCmd := &cobra.Command{
		Use:     "cloc",
		Aliases: []string{"cl"},
		Short:   "Counts lines of code.",
		Long:    "Example: cloc <your_path>",
		GroupID: that.groupID,
		Run: func(cmd *cobra.Command, args []string) {
			cloc := vctrl.NewCloc(&CCtx{cmd: cmd})
			cloc.Run()
		},
	}

	clCmd.Flags().BoolP(vctrl.FlagByFile, "f", false, "Report results for every encountered source file.")
	clCmd.Flags().BoolP(vctrl.FlagDebug, "b", false, "Dump debug log for developer.")
	clCmd.Flags().BoolP(vctrl.FlagSkipDuplicated, "s", false, "Skip duplicated files.")
	clCmd.Flags().BoolP(vctrl.FlagShowLang, "l", false, "Print about all languages and extensions.")
	clCmd.Flags().StringP(vctrl.FlagSortTag, "t", "name", `Sort based on a certain column["name", "files", "blank", "comment", "code"].`)
	clCmd.Flags().StringP(vctrl.FlagOutputType, "o", "default", "Show summary only.")
	clCmd.Flags().StringP(vctrl.FlagExcludeExt, "e", "", "Exclude file name extensions (separated commas).")
	clCmd.Flags().StringP(vctrl.FlagIncludeLang, "L", "", "Include language name (separated commas).")
	clCmd.Flags().StringP(vctrl.FlagMatch, "m", "", "Include file name (regex).")
	clCmd.Flags().StringP(vctrl.FlagNotMatch, "M", "", "Exclude file name (regex).")
	clCmd.Flags().StringP(vctrl.FlagMatchDir, "d", "", "Include dir name (regex).")
	clCmd.Flags().StringP(vctrl.FlagNotMatchDir, "D", "", "Exclude dir name (regex).")

	that.rootCmd.AddCommand(clCmd)
}

func (that *Cli) picRepo() {
	prCmd := &cobra.Command{
		Use:     "pic-repo",
		Aliases: []string{"P"},
		Short:   "Uses github/gitee as picture repositary, especially for markdown.",
		GroupID: that.groupID,
	}

	prCmd.AddCommand(&cobra.Command{
		Use:     "upload",
		Aliases: []string{"u"},
		Short:   "Uploads local picture to remote repo.",
		Long:    "Example: u <path_to_picture...>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			picRepo := vctrl.NewPicRepo()
			for _, arg := range args {
				picRepo.UploadPic(arg)
			}
		},
	})

	prCmd.AddCommand(&cobra.Command{
		Use:     "set-repo",
		Aliases: []string{"s"},
		Short:   "Sets repo name.",
		Run: func(cmd *cobra.Command, args []string) {
			picRepo := vctrl.NewPicRepo()
			picRepo.SetPicRepoName()
		},
	})

	that.rootCmd.AddCommand(prCmd)
}
