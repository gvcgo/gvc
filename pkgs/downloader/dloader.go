package downloader

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	Url              string
	Timeout          time.Duration
	ManuallyRedirect bool
	PostBody         []byte
	Headers          map[string]string
}

func (that *Downloader) setHeaders(req *http.Request) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	if that.Headers == nil || len(that.Headers) == 0 {
		that.Headers = map[string]string{
			"User-Agent": ua,
		}
	}
	if req != nil {
		for key, value := range that.Headers {
			req.Header.Add(key, value)
		}
	}
}

func (that *Downloader) GetUrl() *http.Response {
	if that.Url != "" && utils.VerifyUrls(that.Url) {
		var (
			r   *http.Response
			err error
		)
		httpReq, err := http.NewRequest(http.MethodGet, that.Url, nil)
		if err != nil {
			fmt.Println("[New Request Failed] ", err)
			return nil
		}

		that.setHeaders(httpReq)
		if !that.ManuallyRedirect {
			r, err = http.DefaultClient.Do(httpReq)
		} else {
			r, err = http.DefaultTransport.RoundTrip(httpReq)
		}
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
		return 100
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

func (that *Downloader) PostUrl() (httpRsp *http.Response) {
	if that.Url != "" && len(that.PostBody) > 0 {

		reqBody := strings.NewReader(string(that.PostBody))
		httpReq, err := http.NewRequest(http.MethodPost, that.Url, reqBody)
		if err != nil {
			fmt.Println("[New Request Failed] ", err)
			return nil
		}
		that.setHeaders(httpReq)
		if that.ManuallyRedirect {
			httpRsp, err = http.DefaultTransport.RoundTrip(httpReq)
		} else {
			httpRsp, err = http.DefaultClient.Do(httpReq)
		}
		if err != nil {
			fmt.Println("[Request Failed] ", err)
			return nil
		}
		return httpRsp
	} else {
		fmt.Println("[Illegal URL] URL: ", that.Url)
		return nil
	}
}
