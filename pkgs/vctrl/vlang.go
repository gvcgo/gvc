package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	downloader "github.com/moqsien/gvc/pkgs/fetcher"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Vlang struct {
	Conf *config.GVConfig
	*downloader.Downloader
	env *utils.EnvsHandler
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
		env:        utils.NewEnvsHandler(),
	}
	return
}

func (that *Vlang) download(force bool) string {
	vUrls := that.Conf.Vlang.VlangGiteeUrls
	fmt.Println("Choose your URL to download:")
	fmt.Println("1) Gitee (by default & fast in China);")
	fmt.Println("2) Github .")
	var choice string
	fmt.Scan(&choice)
	if choice == "2" {
		vUrls = that.Conf.Vlang.VlangUrls
	}
	that.Url = vUrls[runtime.GOOS]
	if that.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if force {
			os.RemoveAll(fpath)
		}
		if ok, _ := utils.PathIsExist(fpath); !ok || force {
			if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	}
	return ""
}

func (that *Vlang) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok && !force {
		fmt.Println("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.VlangRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.VlangFilesDir); err != nil {
		os.RemoveAll(config.VlangRootDir)
		os.RemoveAll(zipFilePath)
		fmt.Println("[Unarchive failed] ", err)
		return
	}
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok {
		that.CheckAndInitEnv()
	}
}

func (that *Vlang) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		vlangEnv := fmt.Sprintf(utils.VlangEnv, config.VlangRootDir)
		that.env.UpdateSub(utils.SUB_VLANG, vlangEnv)
	} else {
		envList := map[string]string{
			"PATH": config.VlangRootDir,
		}
		that.env.SetEnvForWin(envList)
	}
}
