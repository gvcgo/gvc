package vctrl

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/downloader"
)

type Code struct {
	Conf *config.Conf
	*downloader.Downloader
}

type TypeMap map[string]string

var CodeType TypeMap = TypeMap{
	"windows-amd64": "win32-x64-archive",
	"windows-arm64": "win32-arm64-archive",
	"linux-amd64":   "linux-x64",
	"linux-arm64":   "linux-arm64",
	"darwin-amd64":  "darwin",
	"darwin-arm64":  "darwin-arm64",
}

func NewCode() *Code {
	return &Code{
		Conf: config.New(),
		Downloader: &downloader.Downloader{
			ManuallyRedirect: true,
		},
	}
}

func (that *Code) getRealUrl() (r string) {
	key := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	if t, ok := CodeType[key]; ok {
		that.Url = fmt.Sprintf("%s?build=%s&os=%s", that.Conf.Config.Code.DownloadUrl, "stable", t)
		that.Timeout = 30 * time.Second
		if resp := that.GetUrl(); resp != nil {
			location := resp.Header["Location"]
			if len(location) > 0 {
				r = strings.Replace(location[0], that.Conf.Config.Code.StableUrl, that.Conf.Config.Code.CdnUrl, 1)
			} else {
				fmt.Println("[Download failed] ", that.Url)
			}
		}
	}
	return
}

func (that *Code) Run() {
	dUrl := that.getRealUrl()
	fmt.Println(dUrl)
}
