package confs

import (
	"fmt"
	"net/http"
	"time"
)

type GithubConf struct {
	DownProxy string `koanf:"down_proxy"`
}

func NewGithubConf() (ghc *GithubConf) {
	ghc = &GithubConf{
		DownProxy: "https://ghproxy.com/",
	}
	return
}

func (that *GithubConf) Reset() {
	that.DownProxy = "https://ghproxy.com/"
}

func (that *GithubConf) testDownProxy() (r bool) {
	if _, err := (&http.Client{Timeout: 20 * time.Second}).Get(that.DownProxy); err == nil {
		r = true
	}
	return
}

func (that *GithubConf) GetDownUrl(oUrl string) (nUrl string) {
	nUrl = oUrl
	if that.testDownProxy() {
		nUrl = fmt.Sprintf("%s%s", that.DownProxy, oUrl)
	}
	return
}
