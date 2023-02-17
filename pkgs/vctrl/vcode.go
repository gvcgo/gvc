package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
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

func NewCode() (co *Code) {
	co = &Code{
		Conf: config.New(),
		Downloader: &downloader.Downloader{
			ManuallyRedirect: true,
		},
	}
	co.initeDirs()
	return
}

func (that *Code) initeDirs() {
	if ok, _ := utils.PathIsExist(config.CodeFileDir); !ok {
		if err := os.MkdirAll(config.CodeFileDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.CodeTarFileDir); !ok {
		if err := os.MkdirAll(config.CodeTarFileDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
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

func (that *Code) download() (r string) {
	dUrl := that.getRealUrl()
	if strings.HasSuffix(dUrl, ".zip") || strings.HasSuffix(dUrl, ".tar.gz") {
		nameList := strings.Split(dUrl, "/")
		fpath := filepath.Join(config.CodeTarFileDir, nameList[len(nameList)-1])
		that.Url = dUrl
		that.Timeout = 180 * time.Second
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size == 0 {
			fmt.Println("[VSCode download failed] ", fpath)
		} else {
			r = fpath
		}
	}
	return
}

func (that *Code) InstallForWin() {}

func (that *Code) InstallForMac() {}

func (that *Code) InstallForLinux() {}

func (that *Code) Run() {
	that.download()
}
