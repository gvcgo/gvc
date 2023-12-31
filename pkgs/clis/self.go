package clis

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/spf13/cobra"
)

func (that *Cli) SetVersionInfo(gitTag, gitHash, gitTime string) {
	that.gitHash = gitHash
	that.gitTag = gitTag
	that.gitTime = gitTime
}

func (that *Cli) showVersion() {
	that.rootCmd.AddCommand(&cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Shows installation info about GVC.",
		GroupID: that.groupID,
		Run: func(cmd *cobra.Command, args []string) {
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

			pattern := "Name: GVC \nVersion: %s\nUpdatedAt: %s\nHomepage: %s\nEmail: %s\n"
			pattern = pattern + "GVCBin: %s\nConfPath: %s\nRemoteConf: %s\nAppsDir: %s\n"
			s := fmt.Sprintf(
				pattern,
				fmt.Sprintf("%s(%s)", that.gitTag, hashTail),
				that.gitTime,
				"https://github.com/moqsien/gvc",
				"moqsien2022@gmail.com",
				// installation info.
				config.GVCDir,
				config.GVConfigPath,
				config.GVCWebdavConfigPath,
				config.GVCInstallDir,
			)
			bp := gprint.NewBlockPrinter(
				s,
				gprint.WithAlign(lipgloss.Left),
				gprint.WithForeground("#FAFAFA"),
				gprint.WithBackground("#874BFD", "#7D56F4"),
				gprint.WithPadding(2, 6),
				gprint.WithHeight(8),
				gprint.WithWidth(78),
				gprint.WithBold(true),
				gprint.WithItalic(true),
			)
			bp.Println()
		},
	})
}

func (that *Cli) checkForUpdate() {
	that.rootCmd.AddCommand(&cobra.Command{
		Use:     "check",
		Aliases: []string{"ch"},
		Short:   "Checks and downloads the latest version of GVC.",
		GroupID: that.groupID,
		Run: func(cmd *cobra.Command, args []string) {
			self := vctrl.NewSelf()
			self.CheckLatestVersion(that.gitTag)
		},
	})
}

func (that *Cli) uninstall() {
	that.rootCmd.AddCommand(&cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"uni"},
		Short:   "Uninstall gvc and Remove all the Apps installed by GVC.",
		GroupID: that.groupID,
		Run: func(cmd *cobra.Command, args []string) {
			self := vctrl.NewSelf()
			self.Uninstall()
		},
	})
}

func (that *Cli) configure() {
	cfgCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Configurations.",
		GroupID: that.groupID,
	}

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "reset",
		Aliases: []string{"R"},
		Short:   "Reset configurations fo GVC to default values.",
		Run: func(cmd *cobra.Command, args []string) {
			dav := vctrl.NewGVCWebdav()
			dav.RestoreDefaultGVConf()
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "Config account info for your remote Repo[WebDAV/Github/Gitee].",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: support github/gitee.
			dav := vctrl.NewGVCWebdav()
			dav.SetWebdavAccount()
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "push",
		Aliases: []string{"p"},
		Short:   "Gather settings and push them to your remote Repo.",
		Run: func(cmd *cobra.Command, args []string) {
			dav := vctrl.NewGVCWebdav()
			dav.GatherAndPushSettings()
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "pull",
		Aliases: []string{"P"},
		Short:   "Pull settings from your remote Repo and apply them.",
		Run: func(cmd *cobra.Command, args []string) {
			dav := vctrl.NewGVCWebdav()
			dav.FetchAndApplySettings()
		},
	})
	that.rootCmd.AddCommand(cfgCmd)
}

// to backup/deploy .ssh files.
func (that *Cli) ssh() {
	sshCmd := &cobra.Command{
		Use:     "ssh",
		Aliases: []string{"s"},
		Short:   "Backups/Deploys your .ssh files.",
		GroupID: that.groupID,
	}

	sshCmd.AddCommand(&cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Backups your .ssh files to remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			dav := vctrl.NewGVCWebdav()
			dav.GatherSSHFiles()
		},
	})

	sshCmd.AddCommand(&cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d"},
		Short:   "Gets .ssh files from remote repo and deploys them.",
		Run: func(cmd *cobra.Command, args []string) {
			dav := vctrl.NewGVCWebdav()
			dav.DeploySSHFiles()
		},
	})
	that.rootCmd.AddCommand(sshCmd)
}
