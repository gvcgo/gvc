package downloader

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	Url     string
	Timeout time.Duration
}

func (that *Downloader) GetUrl() (resp *http.Response) {
	if that.Url != "" && utils.VerifyUrls(that.Url) {
		r, err := (&http.Client{Timeout: that.Timeout}).Get(that.Url)
		if err != nil {
			fmt.Println("[Request Errored] URL: ", that.Url, ", err: ", err)
			return nil
		}
		return r
	} else {
		fmt.Println("[Illegal URL] URL: ", that.Url)
		return nil
	}
}

func (that *Downloader) GetFile(name string, flag int, perm fs.FileMode, force ...bool) (size int64) {
	forceToDownload := false
	if len(force) > 0 && force[0] {
		forceToDownload = true
	}
	if ok, _ := utils.PathIsExist(name); ok && !forceToDownload {
		fmt.Println("[Downloader] File already exists.")
		return
	}
	resp := that.GetUrl()
	if resp == nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		fmt.Println("[open file failed] ", err)
		return
	}
	defer f.Close()
	var dst io.Writer
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(ansi.NewAnsiStdout(), "\n")
		}),
	)
	_ = bar.RenderBlank()
	dst = io.MultiWriter(f, bar)
	size, _ = io.Copy(dst, resp.Body)
	return
}
