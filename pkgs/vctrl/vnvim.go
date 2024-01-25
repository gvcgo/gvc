package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mholt/archiver/v3"
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
	that.setenv()
	return
}

func (that *NVim) setenv() {
	binPath := that.getBinPath()
	if ok, _ := utils.PathIsExist(binPath); ok {
		if runtime.GOOS == utils.Windows {
			that.env.SetEnvForWin(map[string]string{
				"PATH": binPath,
			})
		} else {
			nvimEnv := fmt.Sprintf(utils.NVimEnv, binPath)
			that.env.UpdateSub(utils.SUB_NVIM, nvimEnv)
		}
	}
}

func (that *NVim) Install() {
	that.download()
}

// TODO: neovide.

/*
TODO: synchronize nvim conf files to remote repo.

neovim conf files.
*/

// TODO: astronvim installation.
