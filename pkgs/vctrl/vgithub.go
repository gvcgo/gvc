package vctrl

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/pterm/pterm"
)

type GhDownloader struct {
	Conf     *config.GVConfig
	path     string
	fetcher  *request.Fetcher
	releases map[string]string
}

func NewGhDownloader() (gd *GhDownloader) {
	gd = &GhDownloader{
		path:     filepath.Join(utils.GetHomeDir(), "Downloads"),
		fetcher:  request.NewFetcher(),
		Conf:     config.New(),
		releases: make(map[string]string),
	}
	return
}

func (that *GhDownloader) findFileName(dUrl string) (name string) {
	if strings.Contains(dUrl, "/archive") {
		sList := strings.Split(dUrl, "github.com/")
		if len(sList) == 2 {
			s := sList[1]
			sList = strings.Split(s, "/")
			if len(sList) >= 2 {
				return fmt.Sprintf("%s_code.zip", sList[1])
			}
		}
		return "source_code.zip"
	} else {
		sList := strings.Split(dUrl, "/")
		return fmt.Sprintf("%s_code.zip", sList[len(sList)-1])
	}
}

func (that *GhDownloader) downloadArchive(githubProjectUrl string) {
	// example: https://github.com/moqsien/gvc/archive/refs/heads/main.zip
	mainZipUrl := githubProjectUrl + "/archive/refs/heads/main.zip"
	client := resty.New()
	client.SetTimeout(time.Minute * 3)
	if resp, err := client.R().SetDoNotParseResponse(true).Head(mainZipUrl); err == nil {
		dUrl := mainZipUrl
		if resp.RawResponse.ContentLength <= 0 {
			dUrl = githubProjectUrl + "/archive/refs/heads/master.zip"
		}
		fPath := filepath.Join(that.path, that.findFileName(dUrl))
		that.fetcher.SetUrl(that.Conf.Github.DownProxy + dUrl)
		that.fetcher.Timeout = 30 * time.Minute
		tui.PrintInfo(fmt.Sprintf("[>>>] %s", dUrl))
		that.fetcher.GetFile(fPath, true)
	} else {
		tui.PrintError(err)
	}
}

func (that *GhDownloader) getCurrentTag(githubProjectUrl string) (tag string) {
	// example: https://github.com/moqsien/gvc/releases/latest
	dUrl := githubProjectUrl + "/releases/latest"
	client := resty.New()
	client.SetTimeout(time.Minute * 3)
	if resp, err := client.R().SetDoNotParseResponse(true).Head(that.Conf.Github.DownProxy + dUrl); err == nil {
		_url := resp.RawResponse.Request.URL.String()
		sList := strings.Split(_url, "/")
		return sList[len(sList)-1]
	}
	tui.PrintInfo("Latest tag: ", tag)
	return
}

func (that *GhDownloader) downloadBinary(githubProjectUrl string) {
	// example: https://github.com/moqsien/gvc/releases/expanded_assets/v1.3.1
	if tag := that.getCurrentTag(githubProjectUrl); tag != "" {
		that.fetcher.Url = that.Conf.Github.DownProxy + githubProjectUrl + fmt.Sprintf("/releases/expanded_assets/%s", tag)
		that.fetcher.Timeout = time.Minute * 3
		if resp := that.fetcher.Get(); resp != nil {
			if doc, err := goquery.NewDocumentFromReader(resp.RawBody()); err == nil && doc != nil {
				doc.Find("ul").Find("a").Each(func(i int, s *goquery.Selection) {
					if _url := s.AttrOr("href", ""); _url != "" {
						if filename := s.Find("span").First().Text(); filename != "" && !strings.Contains(filename, "Source code") {
							that.releases[filename], _ = url.JoinPath("https://github.com", _url)
						}
					}
				})
			}
		}
		if len(that.releases) > 0 {
			options := []string{}
			for opt := range that.releases {
				options = append(options, opt)
			}
			selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
			dUrl := that.releases[selectedOption]
			tui.PrintInfo("[Download] ", dUrl)
			that.fetcher.SetUrl(that.Conf.Github.DownProxy + dUrl)
			that.fetcher.SetThreadNum(4)
			that.fetcher.Timeout = 30 * time.Minute
			fPath := filepath.Join(that.path, selectedOption)
			if size := that.fetcher.GetAndSaveFile(fPath, true); size > 0 {
				tui.PrintSuccess(fPath)
			}
		}
	}
}

func (that *GhDownloader) Download(githubProjectUrl string, getSourceCode bool) {
	// exampler: https://github.com/moqsien/gvc
	if !strings.Contains(githubProjectUrl, "github.com/") {
		return
	}
	githubProjectUrl = strings.Split(githubProjectUrl, "/archive")[0]
	githubProjectUrl = strings.Split(githubProjectUrl, "/releases")[0]
	githubProjectUrl = strings.TrimRight(githubProjectUrl, "/")
	if getSourceCode {
		that.downloadArchive(githubProjectUrl)
	} else {
		that.downloadBinary(githubProjectUrl)
	}
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
