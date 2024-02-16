package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	arch "github.com/gvcgo/goutils/pkgs/archiver"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/request"
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
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
		if archive, err := arch.NewArchiver(fpath, config.NVimFileDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				that.RenameNvimDir()
				os.RemoveAll(nvimDir) // removes nvim dir.
				gprint.PrintError("Unarchive failed: %+v", err)
				return
			} else {
				os.RemoveAll(nvimDir) // removes old neovim.
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
		}
		os.RemoveAll(fpath)
	}
}

func (that *NeoVim) renameTreeSitterBinary(toBinPath string) {
	dList, _ := os.ReadDir(config.NVimBinDir)
	for _, d := range dList {
		if !d.IsDir() && strings.Contains(d.Name(), "tree-sitter") {
			os.Rename(
				filepath.Join(config.NVimBinDir, d.Name()),
				filepath.Join(config.NVimBinDir, toBinPath),
			)
			break
		}
	}
}

// Installs tree-sitter, fzf, and lazygit for neovim.
func (that *NeoVim) InstallNeovimDependencies() {
	// tree-sitter
	gprint.PrintInfo("Installing tree-sitter...")
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.NVim.TreeSitterUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	treesitterBinPath := filepath.Join(config.NVimBinDir, "tree-sitter")
	if runtime.GOOS == utils.Windows {
		treesitterBinPath += ".exe"
	}

	var fpath string
	if that.fetcher.Url != "" {
		os.RemoveAll(treesitterBinPath) // removes old tree-sitter
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
		if archive, err := arch.NewArchiver(fpath, config.NVimFileDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				that.renameTreeSitterBinary(treesitterBinPath)
				os.RemoveAll(treesitterBinPath) // removes nvim dir.
				gprint.PrintError("Unarchive failed: %+v", err)
				return
			} else {
				os.RemoveAll(treesitterBinPath)                // removes old neovim.
				that.renameTreeSitterBinary(treesitterBinPath) // rename
				that.setupPrevillage(treesitterBinPath)
			}
		}
	}

	// lazygit
	gprint.PrintInfo("Installing lazygit...")

	// fzf
	gprint.PrintInfo("Installing fzf...")

	that.SetEnvs()
}

// Installs the latest stable version of neovide.
func (that *NeoVim) InstallNeovide() {
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.NVim.NeovideUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	var fpath string

	if that.fetcher.Url != "" {
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath = filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath, true); size <= 0 {
			gprint.PrintError(fmt.Sprintf("Download %s failed.", that.fetcher.Url))
			return
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find nvim package for %s", runtime.GOOS))
	}

	if ok, _ := utils.PathIsExist(config.NVimFileDir); !ok || fpath == "" {
		return
	}

	if runtime.GOOS == utils.MacOS {
		if archive, err := arch.NewArchiver(fpath, config.NVimFileDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				os.RemoveAll(fpath)
				gprint.PrintError("unarchive failed: %+v", err)
				return
			}
		}
		// find .dmg file for MacOS
		dList, _ := os.ReadDir(config.NVimFileDir)
		var dmgPath string
		for _, d := range dList {
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".dmg") {
				dmgPath = filepath.Join(config.NVimFileDir, d.Name())
				break
			}
		}
		os.RemoveAll(fpath)
		fpath = dmgPath
	}

	// removes old neovide files.
	dList, _ := os.ReadDir(config.NVimBinDir)
	for _, d := range dList {
		if strings.Contains(strings.ToLower(d.Name()), "neovide") {
			os.RemoveAll(filepath.Join(config.NVimBinDir, d.Name()))
		}
	}

	if archive, err := arch.NewArchiver(fpath, config.NVimBinDir, false); err == nil {
		_, err = archive.UnArchive()
		defer os.RemoveAll(fpath)
		if err != nil {
			gprint.PrintError("unarchive failed: %+v", err)
			return
		}
		gprint.PrintSuccess("download successed: %s", fpath)
	}

	// fix for MacOS: move binary to $NVimBinDir
	if runtime.GOOS == utils.MacOS {
		dList, _ := os.ReadDir(config.NVimBinDir)
		var appPath string
		for _, d := range dList {
			if strings.HasSuffix(d.Name(), ".app") {
				appPath = filepath.Join(config.NVimBinDir, d.Name())
				break
			}
		}
		if appPath != "" {
			utils.CopyFile(
				filepath.Join(appPath, "Contents", "MacOS", "neovide"),
				filepath.Join(config.NVimBinDir, "neovide"),
			)
			os.RemoveAll(appPath)
		}
	}

	if runtime.GOOS != utils.Windows {
		that.setupPrevillage(filepath.Join(config.NVimBinDir, "neovide"))
	}

	that.SetEnvs()
	os.RemoveAll(fpath)
}

// Installs gnvim config for neovim.
// https://github.com/gvcgo/gnvim
func (that *NeoVim) InstallGnvimConfig() {

}

// Enables/Diables gvc proxy for neovim.
func (that *NeoVim) ToggleProxy() {
	// for git

	// for neovim

}

// Removes all neovim related files.
func (that *NeoVim) RemoveNeovim() {
	if runtime.GOOS != utils.Windows {
		that.env.RemoveSub(utils.SUB_NVIM)
	}
	utils.ClearDir(config.NVimFileDir)
}
