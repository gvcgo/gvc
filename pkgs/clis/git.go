package clis

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

// github download acceleration.
func (that *Cli) github() {
	githubCmd := &cobra.Command{
		Use:     "github",
		Aliases: []string{"gh"},
		Short:   "Github related CLIs.",
		GroupID: that.groupID,
	}

	githubCmd.AddCommand(&cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Short:   "Set a proxy URI for github downloads.",
		Long:    "Example: gh p https://gh.flyinbug.top/gh/",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			vg := vctrl.NewGhDownloader()
			vg.SetReverseProxyForDownload(args[0])
		},
	})

	sourceCodeFlag := "code"
	download := &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Short:   "Download released files or source code from a github repo.",
		Long:    "Example: gh d http://github.com/moqsien/gvc",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			vg := vctrl.NewGhDownloader()
			onlySourceCode, _ := cmd.Flags().GetBool(sourceCodeFlag)
			vg.Download(args[0], onlySourceCode)
		},
	}

	download.Flags().BoolP(sourceCodeFlag, "c", false, "To download source code only.")
	githubCmd.AddCommand(download)

	that.rootCmd.AddCommand(githubCmd)
}

// git installation for windows.
func (that *Cli) git() {

}

// git with proxy written in pure go.
func (that *Cli) gogit() {

}

// lazygit with proxy set/unset.
func (that *Cli) lazygit() {

}
