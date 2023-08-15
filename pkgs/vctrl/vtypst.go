package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/pterm/pterm"
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
	tv.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Typst) download(force bool) string {
	selector := pterm.DefaultInteractiveSelect
	selector.DefaultText = "Choose your download URL"
	optionList := []string{
		"From Gitee[Default]",
		"From Github",
	}
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(optionList).Show()

	var vUrls map[string]string
	switch selectedOption {
	case optionList[1]:
		vUrls = that.Conf.Typst.GithubUrls
	default:
		vUrls = that.Conf.Typst.GiteeUrls
	}
	that.fetcher.Url = vUrls[runtime.GOOS]
	suffix := utils.GetExt(that.fetcher.Url)
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.TypstFilesDir, fmt.Sprintf("typst%s", suffix))
		if force {
			os.RemoveAll(fpath)
		}
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
		tui.PrintInfo("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.TypstRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.TypstFilesDir); err != nil {
		os.RemoveAll(config.TypstRootDir)
		os.RemoveAll(zipFilePath)
		tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
		return
	}
	that.renameDir()
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok {
		that.CheckAndInitEnv()
	} else {
		tui.PrintError("Install typst failed.")
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
