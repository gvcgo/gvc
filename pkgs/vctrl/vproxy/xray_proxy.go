package vproxy

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gogf/gf/encoding/gjson"
	futils "github.com/moqsien/free/pkgs/utils"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Proxy struct {
	Uri string `koanf:"uri"`
	RTT int64  `koanf:"rtt"`
}

func (that *Proxy) GetUri() string {
	return that.Uri
}

func (that *Proxy) SetRTT(rtt int64) {
	that.RTT = rtt
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
			ProxyList: NewVmessList("proxies-raw-vmess.yml"),
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
		// force to fetch new proxy list
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

func (that *ProxyFetcher) GetProxies(force bool) {
	that.ProxyList.Reload()
	if that.Type == Vmess {
		if that.ProxyList.Today() != that.ProxyList.GetDate() || force {
			that.c.OnResponse(func(r *colly.Response) {
				jsonStr, _ := futils.DefaultCrypt.AesDecrypt(r.Body)
				j := gjson.New(jsonStr)
				rawList := j.GetArray("vmess.list")
				vList := []*Proxy{}
				for _, v := range rawList {
					if vmessUri, ok := v.(string); ok {
						vList = append(vList, &Proxy{
							Uri: vmessUri,
						})
					}
				}
				that.ProxyList.Update(vList)
			})
			that.c.Visit("https://gitee.com/moqsien/test/raw/master/conf.txt")
		}
	}
	that.ProxyList.Reload()
}
