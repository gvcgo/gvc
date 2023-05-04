package vctrl

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/hackbrowser/browser"
	"github.com/moqsien/hackbrowser/item"
	"github.com/moqsien/hackbrowser/log"
)

type Browser struct {
	Conf *config.GVConfig
}

func NewBrowser() *Browser {
	return &Browser{
		Conf: config.New(),
	}
}

func (that *Browser) ShowSupportedBrowser() {
	bList := browser.ListBrowsers()
	fmt.Println(color.InCyan(strings.Join(bList, "  ")))
}

func (that *Browser) ShowBackupPath() {
	fmt.Println(color.InCyan(config.GVCBackupDir))
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

func (that *Browser) Save(browserName string, toPush bool) {
	if !that.isBrowserSupported(browserName) {
		fmt.Println(color.InRed("unsupported browser!"))
		return
	}

	log.SetVerbose()
	browsers, err := browser.PickBrowsers(browserName, "")
	if err != nil {
		log.Error(err)
		return
	}

	itemsToSave := []item.Item{
		item.FirefoxPassword,
		item.ChromiumPassword,
		item.YandexPassword,
		item.FirefoxExtension,
		item.ChromiumExtension,
		// item.FirefoxBookmark,
		// item.ChromiumBookmark,
	}

	for _, b := range browsers {
		b.OnlyToSave(itemsToSave)
		data, err := b.BrowsingData(true)
		if err != nil {
			log.Error(err)
		}
		data.Output(config.GVCBackupDir, b.Name(), "json")
	}

	if toPush {
		vconf := NewGVCWebdav()
		vconf.Push()
	}
}
