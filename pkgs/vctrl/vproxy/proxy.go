package vproxy

import (
	"github.com/gocolly/colly/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

type Proxy struct {
	Uri string `koanf:"uri"`
	Rtt int    `koanf:"rtt"`
}

type ProxyList struct {
	Proxies []*Proxy `koanf:"proxies"`
	k       *koanf.Koanf
	parser  *json.JSON
}

type Proxyer struct {
	Conf       *config.GVConfig
	XRay       *core.Instance
	XRayConfig *conf.Config
	ProxyList  *ProxyList
	c          *colly.Collector
}

func NewProxyer() (r *Proxyer) {
	r = &Proxyer{
		Conf: config.New(),
		c:    colly.NewCollector(),
		ProxyList: &ProxyList{
			Proxies: make([]*Proxy, 200),
			k:       koanf.New("."),
			parser:  json.Parser(),
		},
	}
	return
}
