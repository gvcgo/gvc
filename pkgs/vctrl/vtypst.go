package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Typst struct {
	Conf    *config.GVConfig
	fetcher *query.Fetcher
	env     *utils.EnvsHandler
}

func NewTypstVersion() (tv *Typst) {
	tv = &Typst{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	tv.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Typst) download(force bool) string {
	vUrls := that.Conf.Typst.GiteeUrls
	fmt.Println(color.InGreen("Choose your URL to download:"))
	fmt.Println(color.InGreen("1) Gitee (by default & fast in China);"))
	fmt.Println(color.InGreen("2) Github ."))
	var choice string
	fmt.Scan(&choice)
	if choice == "2" {
		vUrls = that.Conf.Typst.GithubUrls
	}
	that.fetcher.Url = vUrls[runtime.GOOS]
	suffix := utils.GetExt(that.fetcher.Url)
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.TypstFilesDir, fmt.Sprintf("typst%s", suffix))
		if force {
			os.RemoveAll(fpath)
		}
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
		fmt.Println(color.InGreen("Vlang is already installed."))
		return
	} else {
		os.RemoveAll(config.TypstRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.TypstFilesDir); err != nil {
		os.RemoveAll(config.TypstRootDir)
		os.RemoveAll(zipFilePath)
		fmt.Println(color.InRed("[Unarchive failed] "), err)
		return
	}
	that.renameDir()
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok {
		that.CheckAndInitEnv()
	} else {
		fmt.Println(color.InRed("Install typst failed!"))
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
