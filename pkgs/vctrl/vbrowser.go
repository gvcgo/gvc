package vctrl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/hackbrowser/browser"
	"github.com/gvcgo/hackbrowser/item"
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
	"github.com/gvcgo/gvc/pkgs/utils/bkm"
)

type Browser struct {
	Conf  *config.GVConfig
	fPath string // files path.
}

func NewBrowser() *Browser {
	b := &Browser{
		Conf:  config.New(),
		fPath: filepath.Join(config.GetGVCWorkDir(), "browser_files"),
	}
	if ok, _ := utils.PathIsExist(b.fPath); !ok {
		os.MkdirAll(b.fPath, os.ModePerm)
	}
	return b
}

func (that *Browser) ShowSupportedBrowser() {
	bList := browser.ListBrowsers()
	fc := gprint.NewFadeColors(bList)
	fc.Println()
}

func (that *Browser) ShowBackupPath() {
	fc := gprint.NewFadeColors("Browser data restore dir: " + config.GVCBackupDir)
	fc.Println()
}

func (that *Browser) isBrowserSupported(name string) bool {
	bList := browser.ListBrowsers()
	for _, bName := range bList {
		if bName == name {
			return true
		}
	}
	return false
}

func (that *Browser) getBrowser(browserName string) browser.Browser {
	browsers, err := browser.PickBrowsers(browserName, "")
	if err != nil {
		gprint.PrintError("%+v", err)
		return nil
	}
	return browsers[0]
}

func (that *Browser) clearTempFiles() {
	fPathList := []string{
		item.TempChromiumKey,
		item.TempChromiumPassword,
		item.TempChromiumCookie,
		item.TempChromiumBookmark,
		item.TempChromiumHistory,
		item.TempChromiumDownload,
		item.TempChromiumCreditCard,
		item.TempChromiumLocalStorage,
		item.TempChromiumSessionStorage,
		item.TempChromiumExtension,
		item.TempYandexPassword,
		item.TempYandexCreditCard,
		item.TempFirefoxKey4,
		item.TempFirefoxPassword,
		item.TempFirefoxCookie,
		item.TempFirefoxBookmark,
		item.TempFirefoxHistory,
		item.TempFirefoxDownload,
		item.TempFirefoxLocalStorage,
		item.TempFirefoxSessionStorage,
		item.TempFirefoxCreditCard,
		item.TempFirefoxExtension,
	}
	for _, f := range fPathList {
		if ok, _ := utils.PathIsExist(f); ok && f != "" {
			os.RemoveAll(f)
		}
	}
}

func (that *Browser) save(browserName string) {
	if !that.isBrowserSupported(browserName) {
		gprint.PrintError("unsupported browser!")
		return
	}

	itemsToSave := []item.Item{
		item.FirefoxPassword,
		item.ChromiumPassword,
		item.YandexPassword,
		item.FirefoxExtension,
		item.ChromiumExtension,
	}

	b := that.getBrowser(browserName)
	if b == nil {
		return
	}
	b.OnlyToSave(itemsToSave)
	data, err := b.BrowsingData(true)
	if err != nil {
		gprint.PrintError("%+v", err)
	}
	data.Output(that.fPath, b.Name(), "json")

	b.CopyBookmark()

	bType := bkm.Chrome
	copyPath := item.TempChromiumBookmark
	if browserName == "firefox" {
		bType = bkm.Firefox
		copyPath = item.TempFirefoxBookmark
	}
	toSavePath := filepath.Join(that.fPath, fmt.Sprintf("%s_bookmarks.html", browserName))
	n := bkm.NewBkmTree(bType, copyPath, toSavePath)
	n.SaveHtml()
	that.clearTempFiles()
}

/*
zip and upload files for browsers.
*/
func (that *Browser) HandleBrowserFiles(browserName string, toDownload bool) {
	repoSyncer := NewSynchronizer()
	remoteFileName := "browserdata.zip"
	if toDownload {
		// download and deploy.
		repoSyncer.DownloadFile(
			that.fPath,
			remoteFileName,
			EncryptByZip,
		)
	} else {
		// zip and upload.
		that.save(browserName)
		repoSyncer.UploadFile(
			that.fPath,
			remoteFileName,
			EncryptByZip,
		)
	}
}
