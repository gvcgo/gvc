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

type Vlang struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
	checker *SumChecker
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	vl.checker = NewSumChecker(vl.Conf)
	vl.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Vlang) download(force bool) string {
	selector := pterm.DefaultInteractiveSelect
	selector.DefaultText = "Choose your download URL"
	optionList := []string{
		"From Gitlab[Default]",
		"From Github",
	}
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(optionList).Show()

	var vUrls map[string]string
	switch selectedOption {
	case optionList[1]:
		vUrls = that.Conf.Vlang.VlangUrls
	default:
		vUrls = that.Conf.Vlang.VlangGiteeUrls
	}
	that.fetcher.Url = vUrls[runtime.GOOS]
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if strings.Contains(that.fetcher.Url, "gitlab.com") && !that.checker.IsUpdated(fpath, that.fetcher.Url) {
			tui.PrintInfo("Current version is already the latest.")
			return fpath
		}
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
		} else if ok && !force {
			return fpath
		}
	}
	return ""
}

func (that *Vlang) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok && !force {
		tui.PrintInfo("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.VlangRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.VlangFilesDir); err != nil {
		os.RemoveAll(config.VlangRootDir)
		os.RemoveAll(zipFilePath)
		tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
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
