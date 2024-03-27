package repo

import (
	"fmt"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/conf"
	"github.com/gvcgo/gvc/utils"
)

/*
Use github/gitee repo as image OSS for markdown.
*/
const (
	GithubPicUrlPattern   string = "https://github.com/%s/%s/raw/main/%s"
	JsDelivrPicUrlPattern string = "https://cdn.jsdelivr.net/gh/%s/%s@main/%s"
	GiteePicUrlPattern    string = "https://gitee.com/%s/%s/raw/master/%s"
)

func UploadPics(repoType RepoType, picFiles ...string) {
	cfg := conf.NewGVConfig()
	cfg.Load()
	repoName := cfg.GetPicRepo()
	repo := NewRepo(repoType, false)
	for _, picFile := range picFiles {
		if utils.FileIsImage(picFile) {
			fName := filepath.Base(picFile)
			if err := repo.Upload(repoName, fName, picFile); err == nil {
				switch repoType {
				case RepoGithub:
					fmt.Println(gprint.CyanStr(GiteePicUrlPattern, cfg.GetGitUserName(), repoName, fName))
					fmt.Println(gprint.CyanStr(JsDelivrPicUrlPattern, cfg.GetGitUserName(), repoName, fName))
				case RepoGitee:
					fmt.Println(gprint.CyanStr(GiteePicUrlPattern, cfg.GetGiteeUserName(), repoName, fName))
				default:
				}
			}
		}
	}
}
