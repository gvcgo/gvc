package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Self struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	checker *SumChecker
}

func NewSelf() (s *Self) {
	s = &Self{
		Conf: config.New(),
		env:  utils.NewEnvsHandler(),
	}
	s.checker = NewSumChecker(s.Conf)
	s.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Self) setEnv() {
	if runtime.GOOS != utils.Windows {
		that.env.UpdateSub(utils.SUB_GVC, fmt.Sprintf(utils.GvcEnv, config.GVCWorkDir))
	} else {
		that.env.SetEnvForWin(map[string]string{"PATH": config.GVCWorkDir})
	}
}

func (that *Self) setShortcut() {
	switch runtime.GOOS {
	case utils.Windows:
		fPath := filepath.Join(config.GVCWorkDir, "gvc.exe")
		if ok, _ := utils.PathIsExist(fPath); ok {
			newPath := filepath.Join(config.GVCWorkDir, "g.exe")
			os.RemoveAll(newPath)
			utils.CopyFile(fPath, newPath)
		}
	default:
		fPath := filepath.Join(config.GVCWorkDir, "gvc")
		if ok, _ := utils.PathIsExist(fPath); ok {
			utils.MkSymLink(fPath, filepath.Join(config.GVCWorkDir, "g"))
		}
	}
}

func (that *Self) Install() {
	utils.MakeDirs(config.GVCWorkDir)
	ePath, _ := os.Executable()
	if strings.Contains(ePath, filepath.Join(utils.GetHomeDir(), ".gvc")) && !strings.Contains(ePath, "bin_temp") {
		// call the installed exe is not allowed.
		return
	}
	name := filepath.Base(ePath)
	if strings.HasSuffix(ePath, "/gvc") || strings.HasSuffix(ePath, "gvc.exe") {
		if _, err := utils.CopyFile(ePath, filepath.Join(config.GVCWorkDir, name)); err == nil {
			that.setEnv()
			that.setShortcut()
		}
		// reset config file to default.
		that.Conf.SetDefault()
		that.Conf.Restore()
	}
	nPath := filepath.Join(filepath.Dir(ePath), ".neobox_encrypt_key.json")
	if ok, _ := utils.PathIsExist(nPath); ok {
		os.RemoveAll(nPath)
	}
}

func (that *Self) Uninstall() {
	cfm := confirm.NewConfirm(confirm.WithTitle("To remove gvc?"))
	cfm.Run()

	if cfm.Result() {
		that.env.RemoveSubs()

		cfmDAV := confirm.NewConfirm(confirm.WithTitle("Save config files to WebDAV before removing gvc?"))
		cfmDAV.Run()

		if cfmDAV.Result() {
			dav := NewGVCWebdav()
			dav.GatherAndPushSettings()
		}
		if ok, _ := utils.PathIsExist(config.GVCWorkDir); ok {
			os.RemoveAll(config.GVCWorkDir)
		}
	} else {
		gprint.PrintInfo("Remove has been aborted.")
	}
}

func (that *Self) ShowPath() {
	content := fmt.Sprintf(
		"Installation Dir: %s\nGVC Config Path: %s\nGVC Webdav Config Path: %s",
		config.GVCWorkDir,
		config.GVConfigPath,
		config.GVCWebdavConfigPath,
	)
	gprint.PrintlnByDefault(content)
}

func (that *Self) CheckLatestVersion(currentVersion string) {
	latest := that.checker.GetLatestGVCVersion()
	if currentVersion == latest {
		gprint.PrintInfo(fmt.Sprintf("Current version: %s is already the latest.", currentVersion))
		return
	}
	cfm := confirm.NewConfirm(confirm.WithTitle("To download the latest version for GVC or not?"))
	cfm.Run()
	if cfm.Result() {
		that.download()
	}
}

func (that *Self) download() {
	dUrl := that.Conf.GVC.GitlabUrls[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)]
	if dUrl != "" {
		fPath := filepath.Join(config.GVCBinTempDir, "gvc.zip")
		fetcher := request.NewFetcher()
		fetcher.SetUrl(dUrl)
		fetcher.Timeout = 20 * time.Minute
		temBinPath := filepath.Join(config.GVCBinTempDir, "gvc")
		if runtime.GOOS == utils.Windows {
			temBinPath += ".exe"
		}
		// remove old files before get a new one
		os.RemoveAll(fPath)
		os.RemoveAll(temBinPath)
		if err := fetcher.DownloadAndDecompress(fPath, config.GVCBinTempDir, true); err != nil {
			gprint.PrintError("%+v", err)
		} else {
			gprint.PrintSuccess(config.GVCBinTempDir)
		}
	}
}
