package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	downloader "github.com/moqsien/gvc/pkgs/fetcher"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Typst struct {
	Conf *config.GVConfig
	*downloader.Downloader
	env *utils.EnvsHandler
}

func NewTypstVersion() (tv *Typst) {
	tv = &Typst{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
		env:        utils.NewEnvsHandler(),
	}
	return
}

func (that *Typst) download(force bool) string {
	vUrls := that.Conf.Typst.GiteeUrls
	fmt.Println("Choose your URL to download:")
	fmt.Println("1) Gitee (by default & fast in China);")
	fmt.Println("2) Github .")
	var choice string
	fmt.Scan(&choice)
	if choice == "2" {
		vUrls = that.Conf.Typst.GithubUrls
	}
	that.Url = vUrls[runtime.GOOS]
	suffix := utils.GetExt(that.Url)
	if that.Url != "" {
		fpath := filepath.Join(config.TypstFilesDir, fmt.Sprintf("typst%s", suffix))
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

func (that *Typst) renameDir() {
	dList, _ := os.ReadDir(config.TypstFilesDir)
	for _, d := range dList {
		if d.IsDir() && strings.Contains(d.Name(), "typst-") {
			os.Rename(filepath.Join(config.TypstFilesDir, d.Name()),
				filepath.Join(config.TypstRootDir))
		}
	}
}

func (that *Typst) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok && !force {
		fmt.Println("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.TypstRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.TypstFilesDir); err != nil {
		os.RemoveAll(config.TypstRootDir)
		os.RemoveAll(zipFilePath)
		fmt.Println("[Unarchive failed] ", err)
		return
	}
	that.renameDir()
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok {
		that.CheckAndInitEnv()
	} else {
		fmt.Println("Install typst failed!")
	}
}

func (that *Typst) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		typstEnv := fmt.Sprintf(utils.TypstEnv, config.TypstRootDir)
		that.env.UpdateSub(utils.SUB_TYPST, typstEnv)
	} else {
		envList := map[string]string{
			"PATH": config.TypstRootDir,
		}
		that.env.SetEnvForWin(envList)
	}
}
