package cmd

import (
	"github.com/gvcgo/gvc/pkg/repo"
	"github.com/spf13/cobra"
)

func RegisterRepo(cli *Cli) {
	parent := &cobra.Command{
		Use:     "repo",
		Aliases: []string{"r"},
		Short:   "Use remote github/gitee repo as OSS.",
		GroupID: cli.groupID,
	}

	picRepo := &cobra.Command{
		Use:     "pic",
		Aliases: []string{"p"},
		Short:   "Uploads pictures to github/gitee.",
		Long:    "Example: g r p <pic_path_1> <pic_path_2> ...",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			repoType, _ := cmd.Flags().GetInt("type")
			repo.UploadPics(repo.RepoType(repoType), args...)
		},
	}
	picRepo.Flags().IntP("type", "t", 0, "repo type, 0: github, 1: gitee")
	parent.AddCommand(picRepo)

	vscode := &cobra.Command{
		Use:     "vscode",
		Aliases: []string{"v"},
		Short:   "Syncs vscode settings/keybindings/extension-list to github/gitee.",
		Run: func(cmd *cobra.Command, args []string) {
			repoType, _ := cmd.Flags().GetInt("type")
			toDownload, _ := cmd.Flags().GetBool("download")
			if !toDownload {
				repo.UploadVSCodeFiles(repo.RepoType(repoType))
				return
			}
			repo.DownloadVSCodeFiles(repo.RepoType(repoType))
		},
	}
	vscode.Flags().IntP("type", "t", 0, "repo type, 0: github, 1: gitee")
	vscode.Flags().BoolP("download", "d", false, "download files from github/gitee")
	parent.AddCommand(vscode)

	dotssh := &cobra.Command{
		Use:     "ssh",
		Aliases: []string{"s"},
		Short:   "Syncs .ssh files to github/gitee.",
		Run: func(cmd *cobra.Command, args []string) {
			repoType, _ := cmd.Flags().GetInt("type")
			toDownload, _ := cmd.Flags().GetBool("download")
			if !toDownload {
				repo.UploadSSHFiles(repo.RepoType(repoType))
				return
			}
			repo.DownloadSSHFiles(repo.RepoType(repoType))
		},
	}
	dotssh.Flags().IntP("type", "t", 0, "repo type, 0: github, 1: gitee")
	dotssh.Flags().BoolP("download", "d", false, "download files from github/gitee")
	parent.AddCommand(dotssh)

	asciinema := &cobra.Command{
		Use:     "asciinema",
		Aliases: []string{"a"},
		Short:   "Syncs asciinema-id file to github/gitee.",
		Run: func(cmd *cobra.Command, args []string) {
			repoType, _ := cmd.Flags().GetInt("type")
			toDownload, _ := cmd.Flags().GetBool("download")
			if !toDownload {
				repo.UploadAsciinemaID(repo.RepoType(repoType))
				return
			}
			repo.DownloadAsciinemaID(repo.RepoType(repoType))
		},
	}
	asciinema.Flags().IntP("type", "t", 0, "repo type, 0: github, 1: gitee")
	asciinema.Flags().BoolP("download", "d", false, "download files from github/gitee")
	parent.AddCommand(asciinema)

	neobox := &cobra.Command{
		Use:     "neobox",
		Aliases: []string{"n"},
		Short:   "Syncs neobox config files to github/gitee.",
		Run: func(cmd *cobra.Command, args []string) {
			repoType, _ := cmd.Flags().GetInt("type")
			toDownload, _ := cmd.Flags().GetBool("download")
			if !toDownload {
				repo.UploadNeoboxConfig(repo.RepoType(repoType))
				return
			}
			repo.DownloadNeoboxConfig(repo.RepoType(repoType))
		},
	}
	neobox.Flags().IntP("type", "t", 0, "repo type, 0: github, 1: gitee")
	neobox.Flags().BoolP("download", "d", false, "download files from github/gitee")
	parent.AddCommand(neobox)

	cli.rootCmd.AddCommand(parent)
}
