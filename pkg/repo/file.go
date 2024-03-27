package repo

import (
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
	"github.com/gvcgo/gvc/utils"
)

/*
Upload file/dir to Repo.
*/
func UploadToRepo(repoType RepoType, encryptEnabled bool, remoteFileName, localFilePath string) {
	cfg := conf.NewGVConfig()
	cfg.Load()
	repoName := cfg.GetBackupRepo()
	repo := NewRepo(repoType, encryptEnabled)
	if err := repo.Upload(repoName, remoteFileName, localFilePath); err != nil {
		gprint.PrintError("upload file failed: %+v", err)
	}
}

/*
Download file/dir from Repo.
*/
func DownloadFromRepo(repoType RepoType, encryptEnabled bool, remoteFileName, localFilePath string) {
	cfg := conf.NewGVConfig()
	cfg.Load()
	repoName := cfg.GetBackupRepo()
	repo := NewRepo(repoType, encryptEnabled)

	backupFileName := fmt.Sprintf("%s.old", localFilePath)
	if ok, _ := gutils.PathIsExist(localFilePath); ok {
		fmt.Println(gprint.CyanStr("File or directory already exists: %s", localFilePath))
		fmt.Println(gprint.YellowStr("Backup the old files or not?[y/N]"))
		var okStr string
		fmt.Scanln(&okStr)
		if strings.ToLower(okStr) == "y" {
			os.RemoveAll(backupFileName)
			os.Rename(localFilePath, backupFileName)
		} else {
			if utils.PathIsDir(localFilePath) {
				os.RemoveAll(localFilePath)
			}
		}
	}

	err := repo.Download(repoName, remoteFileName, localFilePath)
	if err != nil {
		gprint.PrintError("download file failed: %+v", err)
		// recover from backuped files.
		if ok, _ := gutils.PathIsExist(backupFileName); ok {
			os.Rename(backupFileName, localFilePath)
		}
	}
}
