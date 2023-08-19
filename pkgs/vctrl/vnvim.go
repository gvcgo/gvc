package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mholt/archiver/v3"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type NVim struct {
	Conf      *config.GVConfig
	checksum  string
	checktype string
	env       *utils.EnvsHandler
	fetcher   *request.Fetcher
	checker   *SumChecker
}

func NewNVim() (nv *NVim) {
	nv = &NVim{
		fetcher:   request.NewFetcher(),
		Conf:      config.New(),
		checksum:  "",
		checktype: "sha256",
		env:       utils.NewEnvsHandler(),
	}
	nv.checker = NewSumChecker(nv.Conf)
	nv.setup()
	nv.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *NVim) setup() {
	utils.MakeDirs(config.NVimFileDir)
}

func (that *NVim) download() (r string) {
	nurl, ok := that.Conf.NVim.Urls[runtime.GOOS]
	if ok {
		utils.ClearDir(config.NVimFileDir)
		that.fetcher.Url = nurl.Url
		that.fetcher.Timeout = 120 * time.Second
		that.fetcher.SetThreadNum(1)
		fpath := filepath.Join(config.NVimFileDir, fmt.Sprintf("%s%s", nurl.Name, nurl.Ext))
		if !that.checker.IsUpdated(fpath, that.fetcher.Url) {
			tui.PrintInfo("Current version is already the latest.")
			r = fpath
			return
		}
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			r = fpath
		}
	} else {
		tui.PrintError(fmt.Sprintf("Cannot find nvim package for %s", runtime.GOOS))
	}
	if ok, _ := utils.PathIsExist(config.NVimFileDir); ok && r != "" {
		dst := config.NVimFileDir
		if err := archiver.Unarchive(r, dst); err != nil {
			os.RemoveAll(filepath.Dir(that.getBinaryPath()))
			os.RemoveAll(r)
			tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
			return
		}
		os.RemoveAll(r)
	} else {
		os.MkdirAll(config.NVimFileDir, os.ModePerm)
		return
	}
	that.setenv()
	that.initiatePlugins()
	return
}

func (that *NVim) getBinaryPath() (r string) {
	nurl := that.Conf.NVim.Urls[runtime.GOOS]
	r = filepath.Join(config.NVimFileDir, nurl.Name, "bin")
	if runtime.GOOS == utils.Windows {
		utils.MkSymLink(filepath.Join(r, "nvim.exe"), filepath.Join("nvim"))
	}
	return r
}

func (that *NVim) setenv() {
	if ok, _ := utils.PathIsExist(that.getBinaryPath()); ok {
		if runtime.GOOS == utils.Windows {
			that.env.SetEnvForWin(map[string]string{
				"PATH": that.getBinaryPath(),
			})
		} else {
			nvimEnv := fmt.Sprintf(utils.NVimEnv, that.getBinaryPath())
			that.env.UpdateSub(utils.SUB_NVIM, nvimEnv)
		}
		that.setInitFile()
	}
}

func (that *NVim) setInitFile() {
	dst := config.GetNVimInitPath()
	if ok, _ := utils.PathIsExist(dst); ok {
		tui.PrintInfo(fmt.Sprintf("Neovim init file already exists: %s", dst))
		return
	}
	dir_ := filepath.Dir(config.NVimInitBackupPath)
	if ok, _ := utils.PathIsExist(dir_); !ok {
		os.MkdirAll(dir_, os.ModePerm)
	}
	utils.CopyFile(config.NVimInitBackupPath, dst)
}

func (that *NVim) initiatePlugins() {
	that.fetcher.Url = that.Conf.NVim.PluginsUrl
	that.fetcher.Timeout = 120 * time.Second
	fpath := filepath.Join(config.NVimFileDir, "nvim-plugins.zip")
	if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
		if ok, _ := utils.PathIsExist(fpath); ok {
			archiver.Unarchive(fpath, config.NVimFileDir)
			os.Remove(fpath)
		}
	}
	if iList, err := os.ReadDir(config.NVimFileDir); err == nil {
		for _, info := range iList {
			if info.IsDir() && (info.Name() == "autoload" || info.Name() == "plugged") {
				shortcut := filepath.Join(config.GetNVimPlugDir(), info.Name())
				if ok, _ := utils.PathIsExist(shortcut); ok {
					continue
				}
				utils.MkSymLink(filepath.Join(config.NVimFileDir, info.Name()), shortcut)
			}
		}
	}
}

func (that *NVim) Install() {
	that.download()
}
