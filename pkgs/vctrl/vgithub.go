package vctrl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type serctl struct {
	Uuid        string `json:"uuid"`
	DownloadUrl string `json:"download_url"`
}

type GhDownloader struct {
	UrlSerctl  string
	UrlGhProxy string
	Conf       *config.GVConfig
	path       string
	*downloader.Downloader
}

func NewGhDownloader() (gd *GhDownloader) {
	gd = &GhDownloader{
		UrlSerctl:  "https://d.serctl.com/api.rb?dl_start",
		UrlGhProxy: "https://ghproxy.com/%s",
		path:       filepath.Join(utils.GetHomeDir(), "Downloads"),
		Downloader: &downloader.Downloader{},
		Conf:       config.New(),
	}
	return
}

func (that *GhDownloader) sendSerctlPost(zipUrl string) {
	body := serctl{
		Uuid:        "",
		DownloadUrl: zipUrl,
	}
	that.PostBody, _ = json.Marshal(body)
	fmt.Println(string(that.PostBody))
	that.Url = that.UrlSerctl
	that.Headers = map[string]string{
		"referer":    "https://d.serctl.com/?dl_start",
		"origin":     "https://d.serctl.com",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
	}
	if resp := that.PostUrl(); resp != nil {
		content := []byte{}
		resp.Body.Read(content)
		fmt.Println(string(content))
		fmt.Println(resp.StatusCode)
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
	urlList := that.Conf.Proxy.GithubDownload
	if len(urlList) == 0 {
		fmt.Println("No github download acceleration available.")
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
			fmt.Println(err)
		}
	}
}
