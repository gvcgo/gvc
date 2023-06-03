package vctrl

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GhDownloader struct {
	UrlSerctl  string
	UrlGhProxy string
	Conf       *config.GVConfig
	path       string
	fetcher    *request.Fetcher
}

func NewGhDownloader() (gd *GhDownloader) {
	gd = &GhDownloader{
		UrlSerctl:  "https://d.serctl.com/api.rb?dl_start",
		UrlGhProxy: "https://ghproxy.com/%s",
		path:       filepath.Join(utils.GetHomeDir(), "Downloads"),
		fetcher:    request.NewFetcher(),
		Conf:       config.New(),
	}
	return
}

func (that *GhDownloader) sendSerctlPost(zipUrl string) {
	that.fetcher.PostBody = map[string]interface{}{
		"uuid":        "",
		"downloadUrl": zipUrl,
	}
	that.fetcher.Url = that.UrlSerctl
	that.fetcher.Headers = map[string]string{
		"referer":    "https://d.serctl.com/?dl_start",
		"origin":     "https://d.serctl.com",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
	}
	if resp := that.fetcher.Post(); resp != nil {
		content := []byte{}
		resp.RawBody().Read(content)
		fmt.Println(string(content))
		fmt.Println(resp.RawResponse.StatusCode)
	}
}

func (that *GhDownloader) Download(zipUrl string) {
	if !strings.Contains(zipUrl, "github") {
		fmt.Println("[Illegal url] ", zipUrl)
		return
	}
	that.sendSerctlPost(zipUrl)
}

func (that *GhDownloader) OpenByBrowser(chosen int) {
	urlList := that.Conf.Github.AccelUrls
	if len(urlList) == 0 {
		tui.PrintError("No github download acceleration available.")
		return
	}
	var gUrl string
	if chosen >= len(urlList) {
		gUrl = urlList[0]
	} else {
		gUrl = urlList[chosen]
	}
	if gUrl != "" {
		var cmd *exec.Cmd
		if runtime.GOOS == utils.MacOS {
			cmd = exec.Command("open", gUrl)
		} else if runtime.GOOS == utils.Linux {
			cmd = exec.Command("x-www-browser", gUrl)
		} else if runtime.GOOS == utils.Windows {
			cmd = exec.Command("cmd", "/c", "start", gUrl)
		} else {
			return
		}
		if err := cmd.Run(); err != nil {
			tui.PrintError(fmt.Sprintf("Execution failed: %+v", err))
		}
	}
}
