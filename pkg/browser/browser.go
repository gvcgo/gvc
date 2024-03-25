package browser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gtea/gtable"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
	"github.com/moond4rk/hackbrowserdata/browser"
)

func ShowSupportedBrowser() {
	bList := browser.ListBrowsers()
	columns := []gtable.Column{
		{Title: "supported browsers", Width: 150},
	}

	rows := []gtable.Row{}

	for _, bName := range bList {
		rows = append(rows, gtable.Row{
			gprint.CyanStr(bName),
		})
	}

	t := gtable.NewTable(
		gtable.WithColumns(columns),
		gtable.WithRows(rows),
		gtable.WithFocused(true),
		gtable.WithHeight(15),
		gtable.WithWidth(100),
	)
	t.Run()
}

func getBrowserDataDir() string {
	fp := filepath.Join(conf.GetGVCWorkDir(), "browser_data")
	os.MkdirAll(fp, os.ModePerm)
	return fp
}

func getTempDir() string {
	tp := filepath.Join(getBrowserDataDir(), "tmp")
	os.MkdirAll(tp, os.ModePerm)
	return tp
}

func getBrowser(bname string) browser.Browser {
	browsers, err := browser.PickBrowsers(bname, "")
	if err != nil || len(browsers) == 0 {
		gprint.PrintError("%+v", err)
		return nil
	}
	return browsers[0]
}

func supportedOrNot(bname string) bool {
	bList := browser.ListBrowsers()
	for _, b := range bList {
		if b == bname {
			return true
		}
	}
	return false
}

func copyFile(keepTemp bool) {
	dList, _ := os.ReadDir(getTempDir())
	for _, d := range dList {
		if !d.IsDir() {
			dName := strings.ToLower(d.Name())
			if strings.Contains(dName, "extension") || strings.Contains(dName, "password") || strings.Contains(dName, "bookmarks") {
				gutils.CopyAFile(
					filepath.Join(getTempDir(), dName),
					filepath.Join(getBrowserDataDir(), dName),
				)
			}
		}
	}
	if !keepTemp {
		os.RemoveAll(getTempDir())
	}
}

func SaveBrowserData(browserName string, keepTemp bool) {
	if !supportedOrNot(browserName) {
		gprint.PrintError("unsupported browser.")
		return
	}
	b := getBrowser(browserName)
	if b != nil {
		data, err := b.BrowsingData(true)
		if err != nil {
			gprint.PrintError("%+v", err)
		}
		data.Output(getTempDir(), b.Name(), "json")
		copyFile(keepTemp)
		gprint.PrintInfo("data dir: %s", getBrowserDataDir())
	}
}
