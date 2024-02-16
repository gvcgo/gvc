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
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
	"github.com/mholt/archiver/v3"
)

const (
	NvimDirName   string = "nvim"
	NvimScriptWin string = `@echo off
%s %s`
	NvimScriptUnix string = `#!/bin/sh
%s %s
`
)

type NeoVim struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
}

func NewNeoVim() (nv *NeoVim) {
	nv = &NeoVim{
		Conf:    config.New(),
		env:     utils.NewEnvsHandler(),
		fetcher: request.NewFetcher(),
	}

	nv.setup()
	nv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *NeoVim) setup() {
	utils.MakeDirs(
		config.NVimFileDir,
		config.NVimBinDir,
	)
}

// Sets envs for neovim and dependencies.
func (that *NeoVim) SetEnvs() {
	if ok, _ := utils.PathIsExist(config.NVimBinDir); !ok {
		return
	}

	if runtime.GOOS == utils.Windows {
		that.env.SetEnvForWin(map[string]string{
			"PATH": config.NVimBinDir,
		})
	} else {
		that.env.UpdateSub(utils.SUB_NVIM,
			fmt.Sprintf(utils.NVimEnv, config.NVimBinDir))
	}
}

// Renames nvim dir.
func (that *NeoVim) RenameNvimDir() {
	dList, _ := os.ReadDir(config.NVimFileDir)
	for _, d := range dList {
		if d.IsDir() && strings.Contains(d.Name(), NvimDirName) && d.Name() != NvimDirName {
			os.Rename(
				filepath.Join(config.NVimFileDir, d.Name()),
				filepath.Join(config.NVimFileDir, NvimDirName),
			)
			return
		}
	}
}

// Adds +x previlledge.
func (that *NeoVim) setupPrevillage(binPath string) {
	if ok, _ := utils.PathIsExist(binPath); ok {
		utils.ExecuteSysCommand(false, "chmod", "+x", binPath)
	}
}

// Installs the latest stable version of neovim.
func (that *NeoVim) InstallNeovim() {
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.NVim.NvimUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	var fpath string

	nvimDir := filepath.Join(config.NVimFileDir, NvimDirName)
	if that.fetcher.Url != "" {
		utils.ClearDir(nvimDir)
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath = filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath); size <= 0 {
			gprint.PrintError(fmt.Sprintf("Download %s failed.", that.fetcher.Url))
			return
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find nvim package for %s", runtime.GOOS))
	}

	if ok, _ := utils.PathIsExist(config.NVimFileDir); ok && fpath != "" {
		dst := config.NVimFileDir
		if err := archiver.Unarchive(fpath, dst); err != nil {
			that.RenameNvimDir()
			os.RemoveAll(nvimDir)
			gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
		} else {
			that.RenameNvimDir()
			if runtime.GOOS == utils.Windows {
				binPath := filepath.Join(nvimDir, "bin", "nvim.exe")
				scriptPath := filepath.Join(config.NVimBinDir, "nvim.bat")
				os.WriteFile(scriptPath,
					[]byte(fmt.Sprintf("%s %s", binPath, `%*`)),
					os.ModePerm,
				)
			} else {
				binPath := filepath.Join(nvimDir, "bin", "nvim")
				that.setupPrevillage(binPath)
				scriptPath := filepath.Join(config.NVimBinDir, "nvim")
				os.WriteFile(scriptPath,
					[]byte(fmt.Sprintf("%s %s", binPath, `$*`)),
					os.ModePerm,
				)
				that.setupPrevillage(scriptPath)
			}
			that.SetEnvs()
		}
		os.RemoveAll(fpath)
	}
}

// Installs the latest stable version of neovide.
func (that *NeoVim) InstallNeovide() {

}

// Installs tree-sitter, fzf, and lazygit for neovim.
func (that *NeoVim) InstallNeovimDependencies() {

}

// Installs gnvim config for neovim.
// https://github.com/gvcgo/gnvim
func (that *NeoVim) InstallGnvimConfig() {

}

// Enables/Diables gvc proxy for neovim.
func (that *NeoVim) ToggleProxy() {

}

// Removes all neovim related files.
func (that *NeoVim) RemoveNeovim() {

}
