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
				gprint.WithWidth(70),
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
		Short:   "Resets configs for GVC to default values.",
		Run: func(cmd *cobra.Command, args []string) {
			s := vctrl.NewSelf()
			s.ResetConf()
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "repo",
		Aliases: []string{"r"},
		Short:   "Sets account info for your remote Repo[Github/Gitee].",
		Run: func(cmd *cobra.Command, args []string) {
			s := vctrl.NewSynchronizer()
			s.Setup()
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "upload-conf",
		Aliases: []string{"u"},
		Short:   "Uploads gvc config file to your remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			s := vctrl.NewSelf()
			s.HandleGvcConfigFile(false)
		},
	})

	cfgCmd.AddCommand(&cobra.Command{
		Use:     "download-conf",
		Aliases: []string{"d"},
		Short:   "Downloads gvc config file from your remote repo.",
		Run: func(cmd *cobra.Command, args []string) {
			s := vctrl.NewSelf()
			s.HandleGvcConfigFile(true)
		},
	})
	that.rootCmd.AddCommand(cfgCmd)
}

// to backup/deploy .ssh files.
// func (that *Cli) ssh() {
// 	sshCmd := &cobra.Command{
// 		Use:     "ssh",
// 		Aliases: []string{"s"},
// 		Short:   "Backups/Deploys your .ssh files.",
// 		GroupID: that.groupID,
// 	}

// 	sshCmd.AddCommand(&cobra.Command{
// 		Use:     "backup",
// 		Aliases: []string{"b"},
// 		Short:   "Backups your .ssh files to remote repo.",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			dav := vctrl.NewGVCWebdav()
// 			dav.GatherSSHFiles()
// 		},
// 	})

// 	sshCmd.AddCommand(&cobra.Command{
// 		Use:     "deploy",
// 		Aliases: []string{"d"},
// 		Short:   "Gets .ssh files from remote repo and deploys them.",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			dav := vctrl.NewGVCWebdav()
// 			dav.DeploySSHFiles()
// 		},
// 	})
// 	that.rootCmd.AddCommand(sshCmd)
// }
