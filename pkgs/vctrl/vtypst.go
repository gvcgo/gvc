package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/request"
	"github.com/mholt/archiver/v3"
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
)

type Typst struct {
	Conf    *config.GVConfig
	fetcher *request.Fetcher
	env     *utils.EnvsHandler
}

func NewTypstVersion() (tv *Typst) {
	tv = &Typst{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	tv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *Typst) download(force bool) string {
	// vUrls := that.Conf.Typst.GithubUrls
	// if runtime.GOOS == utils.Windows {
	// 	that.fetcher.Url = vUrls[runtime.GOOS]
	// } else {
	// 	that.fetcher.Url = vUrls[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)]
	// }
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.Typst.TypstUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	suffix := utils.GetExt(that.fetcher.Url)
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.TypstFilesDir, fmt.Sprintf("typst%s", suffix))
		if force {
			os.RemoveAll(fpath)
		}
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(2)
		if ok, _ := utils.PathIsExist(fpath); !ok || force {
			if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
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
		gprint.PrintInfo("Typst is already installed.")
		return
	} else {
		os.RemoveAll(config.TypstRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.TypstFilesDir); err != nil {
		os.RemoveAll(config.TypstRootDir)
		os.RemoveAll(zipFilePath)
		gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
		return
	}
	that.renameDir()
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok {
		that.CheckAndInitEnv()
	} else {
		gprint.PrintError("Install typst failed.")
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
