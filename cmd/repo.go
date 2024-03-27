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
		Short:   "Upload pictures to github/gitee.",
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

	cli.rootCmd.AddCommand(parent)
}
