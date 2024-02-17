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
	NvimProxyWin string = `set HTTP_PROXY =%s
set HTTPS_PROXY=%s`
	NvimProxyUnix string = `export HTTP_PROXY=%s
export HTTPS_PROXY=%s`
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
						[]byte(fmt.Sprintf(NvimScriptWin, binPath, `%*`)),
						os.ModePerm,
					)
				} else {
					binPath := filepath.Join(nvimDir, "bin", "nvim")
					that.setupPrevillage(binPath)
					scriptPath := filepath.Join(config.NVimBinDir, "nvim")
					os.WriteFile(scriptPath,
						[]byte(fmt.Sprintf(NvimScriptUnix, binPath, `$*`)),
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
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath = filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath); size <= 0 {
			gprint.PrintError(fmt.Sprintf("Download %s failed.", that.fetcher.Url))
			return
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find tree-sitter package for %s", runtime.GOOS))
	}

	if ok, _ := utils.PathIsExist(config.NVimFileDir); ok && fpath != "" {
		if archive, err := arch.NewArchiver(fpath, config.NVimBinDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				that.renameTreeSitterBinary(treesitterBinPath)
				os.RemoveAll(treesitterBinPath) // removes tree-sitter.
				os.RemoveAll(fpath)
				gprint.PrintError("Unarchive failed: %+v", err)
				return
			} else {
				os.RemoveAll(treesitterBinPath)                // removes old tree-sitter.
				that.renameTreeSitterBinary(treesitterBinPath) // rename
				that.setupPrevillage(treesitterBinPath)
			}
		}
	}

	os.RemoveAll(fpath)

	// lazygit
	gprint.PrintInfo("Installing lazygit...")
	if runtime.GOOS == utils.Windows {
		scriptPath := filepath.Join(config.NVimBinDir, "lazygit.bat")
		os.WriteFile(scriptPath,
			[]byte(fmt.Sprintf(NvimScriptWin, "g git lazygit", `%*`)),
			os.ModePerm,
		)
	} else {
		scriptPath := filepath.Join(config.NVimBinDir, "nvim")
		os.WriteFile(scriptPath,
			[]byte(fmt.Sprintf(NvimScriptUnix, "g git lazygit", `$*`)),
			os.ModePerm,
		)
		that.setupPrevillage(scriptPath)
	}

	// fzf
	gprint.PrintInfo("Installing fzf...")

	uList = gh.ParseReleasesForGithubProject(that.Conf.NVim.FzFUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	fzfBinPath := filepath.Join(config.NVimBinDir, "fzf")
	if runtime.GOOS == utils.Windows {
		fzfBinPath += ".exe"
	}

	fpath = ""
	if that.fetcher.Url != "" {
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath = filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath); size <= 0 {
			gprint.PrintError(fmt.Sprintf("Download %s failed.", that.fetcher.Url))
			that.SetEnvs()
			return
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find fzf package for %s", runtime.GOOS))
	}

	if ok, _ := utils.PathIsExist(config.NVimFileDir); ok && fpath != "" {
		if archive, err := arch.NewArchiver(fpath, config.NVimBinDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				os.RemoveAll(fzfBinPath) // removes fzf dir.
				gprint.PrintError("Unarchive failed: %+v", err)
				that.SetEnvs()
				return
			} else {
				os.RemoveAll(fzfBinPath) // removes old fzf.
				that.setupPrevillage(fzfBinPath)
			}
		}
	}

	os.RemoveAll(fpath)
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
		gprint.PrintError(fmt.Sprintf("Cannot find neovide package for %s", runtime.GOOS))
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

// Finds neovim config dir.
func (that *NeoVim) FindConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	var configDir string
	/*
		https://neovim.io/doc/user/starting.html#standard-path
		Unix:         ~/.config                   ~/.config/nvim
		Windows:      ~/AppData/Local             ~/AppData/Local/nvim
	*/
	if runtime.GOOS != utils.Windows {
		configDir = filepath.Join(homeDir, ".config", "nvim")
	} else {
		configDir = filepath.Join(homeDir, "AppData", "Local", "nvim")
	}
	return configDir
}

// Installs gnvim config for neovim.
// https://github.com/gvcgo/gnvim
func (that *NeoVim) InstallGnvimConfig() {
	configDir := that.FindConfigDir()
	gnvimDir := filepath.Join(configDir, "lua", "gnvim")

	if ok, _ := utils.PathIsExist(gnvimDir); ok {
		gprint.PrintInfo("updating gnvim...")
		utils.ExecuteSysCommand(false, "g", "git", "pull", "-d", gnvimDir)
		return
	}

	parentDir := filepath.Dir(configDir)
	// backup old neovim configs.
	gprint.PrintInfo("backuping old configs...")
	if ok, _ := utils.PathIsExist(configDir); ok {
		zipName := fmt.Sprintf("nvim_config_%v.zip", time.Now())
		if archive, err := arch.NewArchiver(configDir, parentDir, false); err == nil {
			archive.SetZipName(zipName)
			err = archive.ZipDir()
			if err != nil {
				gprint.PrintError("Zip dir error: %+v", err)
				return
			}
		}
		os.RemoveAll(configDir)
	}

	// install gnvim.
	utils.ExecuteSysCommand(
		false,
		"g",
		"git",
		"clone",
		"-d",
		parentDir,
		that.Conf.NVim.GNvimUrl,
	)
	os.Rename(filepath.Join(parentDir, "gnvim"), configDir)
}

// Enables/Diables gvc proxy for neovim.
func (that *NeoVim) ToggleProxy() {
	// for git
	utils.ExecuteSysCommand(
		false,
		"g",
		"git",
		"toggle-ssh-proxy",
	)

	nvimDir := filepath.Join(config.NVimFileDir, NvimDirName)
	if exist, _ := utils.PathIsExist(nvimDir); !exist {
		gprint.PrintError("neovim is not installed by gvc.")
		return
	}

	homeDir, _ := os.UserHomeDir()
	confPath := filepath.Join(homeDir, ".ssh", "config")
	backupConfPath := filepath.Join(homeDir, ".ssh", "config.bak")
	// for neovim
	ok, _ := utils.PathIsExist(confPath)
	ok1, _ := utils.PathIsExist(backupConfPath)
	if !ok && !ok1 {
		gprint.PrintWarning(`Please set proxy for ssh using "g git ssh-proxy-fix".`)
		return
	}

	content, _ := os.ReadFile(filepath.Join(config.GVCDir, DefaultProxyFileName))
	if len(content) == 0 {
		gprint.PrintError("No proxy is set for gvc.")
		gprint.PrintInfo(`Set a proxy using "g git proxy http://127.0.0.1:port".`)
		return
	}

	pxy := string(content)
	switch runtime.GOOS {
	case utils.Windows:
		binPath := filepath.Join(nvimDir, "bin", "nvim.exe")
		scriptPath := filepath.Join(config.NVimBinDir, "nvim.bat")

		content := fmt.Sprintf(NvimScriptWin, binPath, `%*`)
		if ok {
			// enables proxy.
			proxyStr := fmt.Sprintf(NvimProxyWin, pxy, pxy)
			content = fmt.Sprintf(NvimScriptWin, proxyStr+"\n"+binPath, `%*`)
		}
		os.WriteFile(scriptPath,
			[]byte(content),
			os.ModePerm,
		)
	default:
		binPath := filepath.Join(nvimDir, "bin", "nvim")
		that.setupPrevillage(binPath)
		scriptPath := filepath.Join(config.NVimBinDir, "nvim")

		content := fmt.Sprintf(NvimScriptUnix, binPath, `$*`)
		if ok {
			// enables proxy.
			proxyStr := fmt.Sprintf(NvimProxyUnix, pxy, pxy)
			content = fmt.Sprintf(NvimScriptUnix, proxyStr+"\n"+binPath, `$*`)
		}
		os.WriteFile(scriptPath,
			[]byte(content),
			os.ModePerm,
		)
		that.setupPrevillage(scriptPath)
	}
}

// Removes all neovim related files.
func (that *NeoVim) RemoveNeovim() {
	if runtime.GOOS != utils.Windows {
		that.env.RemoveSub(utils.SUB_NVIM)
	}
	utils.ClearDir(config.NVimFileDir)
}

// Backups/Deploys neovim config.
func (that *NeoVim) HandleNeovimConfig(toDownload bool) {
	fPath := that.FindConfigDir()
	remoteFileName := "neovim_configs.zip"
	repoSyncer := NewSynchronizer()
	if toDownload {
		// download and deploy.
		repoSyncer.DownloadFile(
			fPath,
			remoteFileName,
			EncryptByZip,
		)
	} else {
		repoSyncer.UploadFile(
			fPath,
			remoteFileName,
			EncryptByZip,
		)
	}
}
