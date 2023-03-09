package vctrl

import (
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
)

type Proxy struct {
	*downloader.Downloader
	Conf      *config.GVConfig
	ProxyList *config.ProxyList
	date      string
}

func NewProxy() (p *Proxy) {
	return &Proxy{
		Conf: config.New(),
		ProxyList: &config.ProxyList{
			Proxies: []*config.Proxy{},
		},
		date:       time.Now().Format("2006-01-02"),
		Downloader: &downloader.Downloader{},
	}
}
