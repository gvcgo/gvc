package vctrl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

var ZigOSArchMap = map[string]string{
	"windows-x86_64":  "windows_amd64",
	"windows-aarch64": "windows_arm64",
	"macos-x86_64":    "darwin_amd64",
	"macos-aarch64":   "darwin_arm64",
	"linux-x86_64":    "linux_amd64",
	"linux-aarch64":   "linux_arm64",
}

// https://github.com/ziglang/zig
// https://ziglang.org/
type Zig struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
	zigList map[string]string
}

func NewZig() (z *Zig) {
	z = &Zig{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
		zigList: map[string]string{},
	}
	z.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *Zig) GetZigList() {
	if len(that.zigList) > 0 {
		return
	}
	that.fetcher.SetUrl(that.Conf.Zig.ZigDownloadUrl)
	that.fetcher.Timeout = time.Minute * 5
	if resp := that.fetcher.Get(); resp != nil {
		doc, err := goquery.NewDocumentFromReader(resp.RawBody())
		if err != nil {
			gprint.PrintError(fmt.Sprintf("Parse page errored: %+v", err))
		}
		if doc == nil {
			gprint.PrintError(fmt.Sprintf("Cannot parse html for %s", that.fetcher.Url))
			os.Exit(1)
		}
		// Latest stable version only.
		doc.Find("table").Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			if href != "" {
				for k, v := range ZigOSArchMap {
					if strings.Contains(href, k) {
						that.zigList[v] = href
					}
				}
			}
		})
	}
	// fmt.Printf("%+v\n", that.zigList)
}

func (that *Zig) Install() {

}

func (that *Zig) InstalZls() {

}
