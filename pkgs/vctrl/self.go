package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Self struct {
	Conf *config.GVConfig
	env  *utils.EnvsHandler
}

func NewSelf() (s *Self) {
	s = &Self{
		Conf: config.New(),
		env:  utils.NewEnvsHandler(),
	}
	s.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *Self) setEnv() {
	if runtime.GOOS != utils.Windows {
		that.env.UpdateSub(utils.SUB_GVC, fmt.Sprintf(utils.GvcEnv, config.GVCDir))
	} else {
		that.env.SetEnvForWin(map[string]string{"PATH": config.GVCDir})
	}
}

func (that *Self) setShortcut() {
	switch runtime.GOOS {
	case utils.Windows:
		fPath := filepath.Join(config.GVCDir, "gvc.exe")
		if ok, _ := utils.PathIsExist(fPath); ok {
			newPath := filepath.Join(config.GVCDir, "g.exe")
			os.RemoveAll(newPath)
			utils.CopyFile(fPath, newPath)
		}
	default:
		fPath := filepath.Join(config.GVCDir, "gvc")
		if ok, _ := utils.PathIsExist(fPath); ok {
			utils.MkSymLink(fPath, filepath.Join(config.GVCDir, "g"))
		}
	}
}

func (that *Self) Install() {
	utils.MakeDirs(config.GVCDir)
	ePath, _ := os.Executable()
	if strings.Contains(ePath, filepath.Join(utils.GetHomeDir(), ".gvc")) && !strings.Contains(ePath, "bin_temp") {
		// call the installed exe is not allowed.
		return
	}
	name := filepath.Base(ePath)
	if strings.HasSuffix(ePath, "/gvc") || strings.HasSuffix(ePath, "gvc.exe") {
		if _, err := utils.CopyFile(ePath, filepath.Join(config.GVCDir, name)); err != nil {
			gprint.PrintError("%+v", err)
		}
		that.setEnv()
		that.setShortcut()
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
		if ok, _ := utils.PathIsExist(config.GVCInstallDir); ok {
			os.RemoveAll(config.GVCInstallDir)
		}
		if ok, _ := utils.PathIsExist(config.GVCDir); ok {
			os.RemoveAll(config.GVCDir)
		}
	} else {
		gprint.PrintInfo("Remove has been aborted.")
	}
}

func (that *Self) ShowPath() {
	content := fmt.Sprintf(
		"GVCDir: %s\nGVCInstallDir:%s\nGVConfPath: %s\nDAVConfPath: %s",
		config.GVCDir,
		config.GVCInstallDir,
		config.GVConfigPath,
		config.GVCWebdavConfigPath,
	)
	bp := gprint.NewBlockPrinter(
		content,
		gprint.WithAlign(lipgloss.Left),
		gprint.WithForeground("#FAFAFA"),
		gprint.WithBackground("#874BFD", "#7D56F4"),
		gprint.WithPadding(2, 6),
		gprint.WithHeight(5),
		gprint.WithWidth(78),
		gprint.WithBold(true),
		gprint.WithItalic(true),
	)
	bp.Println()
}

func (that *Self) CheckLatestVersion(currentVersion string) {
	gUrl := "https://github.com/moqsien/gvc/releases/latest"
	gUrl = that.Conf.GVCProxy.WrapUrl(gUrl)
	fetcher := request.NewFetcher()
	fetcher.SetUrl(gUrl)
	fetcher.Timeout = 30 * time.Second
	resp := fetcher.Get()
	if resp == nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		gprint.PrintError("%+v", err)
		return
	}

	latest := doc.Find("span.css-truncate-target").Find("span").Text()
	if currentVersion == strings.TrimSpace(latest) {
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
	dUrl := that.Conf.GVC.Urls[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)]
	if dUrl != "" {
		fPath := filepath.Join(config.GVCBinTempDir, "gvc.zip")
		fetcher := request.NewFetcher()
		dUrl = that.Conf.GVCProxy.WrapUrl(dUrl)
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
