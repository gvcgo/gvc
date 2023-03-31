package vproxy

import (
	"strings"

	"github.com/gocolly/colly/v2"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Proxy struct {
	Uri string `koanf:"uri"`
}

func (that *Proxy) GetUri() string {
	return that.Uri
}

type ProxyFetcher struct {
	ProxyList Proxies
	Type      ProxyType
	Conf      *config.GVConfig
	c         *colly.Collector
	filter    map[string]struct{}
}

func NewProxyFetcher(typ ...ProxyType) (r *ProxyFetcher) {
	if len(typ) == 0 || typ[0] == "vmess" {
		r = &ProxyFetcher{
			ProxyList: NewVmessList(),
			Conf:      config.New(),
			c:         colly.NewCollector(),
			filter:    map[string]struct{}{},
		}
		r.Type = "vmess"
	}
	return
}

func (that *ProxyFetcher) parseProxy(body []byte) any {
	r := string(body)
	if that.Type == Vmess {
		if !strings.Contains(r, "vmess") {
			r = utils.DecodeBase64(r)
		}
		result := []*Proxy{}
		if strings.Contains(r, "vmess") {
			for _, p := range strings.Split(r, "\n") {
				pUrl := strings.Trim(p, " ")
				if !strings.HasPrefix(pUrl, "vmess") {
					// fmt.Println(pUrl)
					continue
				}
				if _, ok := that.filter[pUrl]; !ok {
					that.filter[pUrl] = struct{}{}
					result = append(result, &Proxy{
						Uri: pUrl,
					})
				}
			}
		}
		return result
	}
	return nil
}

func (that *ProxyFetcher) GetProxyList(force bool) {
	that.ProxyList.Reload()
	if that.Type == Vmess {
		if that.ProxyList.Today() != that.ProxyList.GetDate() || force {
			that.filter = map[string]struct{}{}
			pList := []*Proxy{}
			for _, url := range that.Conf.Proxy.GetSubUrls() {
				// that.collector.SetRequestTimeout(5 * time.Second)
				that.c.OnResponse(func(r *colly.Response) {
					res := that.parseProxy(r.Body)
					result, ok := res.([]*Proxy)
					if ok {
						pList = append(pList, result...)
					}
				})
				that.c.Visit(url)
			}
			that.ProxyList.Update(pList)
		}
	}
	that.ProxyList.Reload()
}
