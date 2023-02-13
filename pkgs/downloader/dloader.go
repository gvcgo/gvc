package downloader

import (
	"fmt"
	"net/http"
	"time"

	"github.com/moqsien/gvc/pkgs/utils"
)

type Downloader struct {
	Url     string
	Timeout time.Duration
}

func (that *Downloader) Download() (resp *http.Response) {
	if that.Url != "" && utils.VerifyUrls(that.Url) {
		fmt.Println(that.Url)
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
