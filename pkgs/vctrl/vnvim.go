package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mholt/archiver/v3"
	arch "github.com/moqsien/goutils/pkgs/archiver"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
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
}

func NewNVim() (nv *NVim) {
	nv = &NVim{
		fetcher:   request.NewFetcher(),
		Conf:      config.New(),
		checksum:  "",
		checktype: "sha256",
		env:       utils.NewEnvsHandler(),
	}
	nv.setup()
	nv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *NVim) setup() {
	utils.MakeDirs(config.NVimFileDir)
}

func (that *NVim) getBinPath() string {
	dfinder := utils.NewBinaryFinder()
	dfinder.SetStartDir(config.NVimFileDir)
	dfinder.SetParentDirName("bin")
	fName := "nvim"
	if runtime.GOOS == utils.Windows {
		fName = "nvim.exe"
	}
	dfinder.SetUniqueFileName(fName)
	return dfinder.String()
}

func (that *NVim) download() (r string) {
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.NVim.NvimUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	if that.fetcher.Url != "" {
		utils.ClearDir(config.NVimFileDir)
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath := filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			r = fpath
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find nvim package for %s", runtime.GOOS))
	}
	if ok, _ := utils.PathIsExist(config.NVimFileDir); ok && r != "" {
		dst := config.NVimFileDir
		if err := archiver.Unarchive(r, dst); err != nil {
			os.RemoveAll(filepath.Dir(that.getBinPath()))
			os.RemoveAll(r)
			gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
			return
		}
		os.RemoveAll(r)
	} else {
		os.MkdirAll(config.NVimFileDir, os.ModePerm)
		return
	}
	return
}

func (that *NVim) setenv() {
	nvimBinPath := that.getBinPath()
	neovideBinPath := that.FindNeovideBinary()
	if ok, _ := utils.PathIsExist(nvimBinPath); ok {
		exist, _ := utils.PathIsExist(neovideBinPath)
		if runtime.GOOS == utils.Windows {
			that.env.SetEnvForWin(map[string]string{
				"PATH": nvimBinPath,
			})
			if exist {
				that.env.SetEnvForWin(map[string]string{
					"PATH": neovideBinPath,
				})
			}
		} else {
			nvimEnv := fmt.Sprintf(utils.NVimEnv, nvimBinPath)
			if exist {
				neovideEnv := fmt.Sprintf(utils.NVimEnv, neovideBinPath)
				nvimEnv = fmt.Sprintf("%s\n%s", nvimEnv, neovideEnv)
			}
			that.env.UpdateSub(utils.SUB_NVIM, nvimEnv)
		}
	}
}

func (that *NVim) Install() {
	that.download()
	that.setenv()
}

func (that *NVim) downloadNeovide() (r string) {
	gh := NewGhDownloader()
	uList := gh.ParseReleasesForGithubProject(that.Conf.NVim.NeovideUrl)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(uList[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)])
	if that.fetcher.Url != "" {
		utils.ClearDir(config.NeovideBinDir)
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		fpath := filepath.Join(config.NVimFileDir, filepath.Base(that.fetcher.Url))
		if size := that.fetcher.GetAndSaveFile(fpath, true); size > 0 {
			r = fpath
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find nvim package for %s", runtime.GOOS))
	}
	return
}

func (that *NVim) FindNeovideBinary() string {
	dfinder := utils.NewBinaryFinder()
	dfinder.SetStartDir(config.NVimFileDir)
	dName := "neovide"
	if runtime.GOOS == utils.MacOS {
		dName = "MacOS"
	}
	dfinder.SetParentDirName(dName)
	fName := "neovide"
	if runtime.GOOS == utils.Windows {
		fName = "neovide.exe"
	}
	dfinder.SetUniqueFileName(fName)
	return dfinder.String()
}

func (that *NVim) InstallNeovide() {
	srcPath := that.downloadNeovide()
	if srcPath == "" {
		return
	}
	if runtime.GOOS == utils.MacOS {
		if archive, err := arch.NewArchiver(srcPath, config.NVimFileDir, false); err == nil {
			_, err = archive.UnArchive()
			if err != nil {
				gprint.PrintError("unarchive failed: %+v", err)
				return
			}
		}
		dmgPath := filepath.Join(config.NVimFileDir, "neovide.dmg")
		if ok, _ := utils.PathIsExist(dmgPath); ok {
			os.RemoveAll(srcPath)
			srcPath = dmgPath
		} else {
			gprint.PrintError("neovide.dmg not found: %s", dmgPath)
		}
	}
	if archive, err := arch.NewArchiver(srcPath, config.NeovideBinDir, false); err == nil {
		_, err = archive.UnArchive()
		defer os.RemoveAll(srcPath)
		if err != nil {
			gprint.PrintError("unarchive failed: %+v", err)
			return
		}
		gprint.PrintSuccess("download successed: %s", srcPath)
	}
	that.setenv()
}

/*
TODO: synchronize nvim conf files to remote repo.

neovim conf files.
*/

// TODO: astronvim installation.
