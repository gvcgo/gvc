package vctrl

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Qingluan/merkur"
	"github.com/gocolly/colly/v2"
	config "github.com/moqsien/gvc/pkgs/confs"
)

type Proxy struct {
	Conf      *config.GVConfig
	ProxyList *config.ProxyList
	date      string
	collector *colly.Collector
	filter    map[string]struct{}
}

func NewProxy() (p *Proxy) {
	return &Proxy{
		Conf: config.New(),
		ProxyList: &config.ProxyList{
			Proxies: []*config.Proxy{},
		},
		date:      time.Now().Format("2006-01-02"),
		collector: colly.NewCollector(),
	}
}

func (that *Proxy) GetProxyList() {
	that.filter = map[string]struct{}{}
	for _, url := range that.Conf.Proxy.GetSubUrls() {
		// that.collector.SetRequestTimeout(5 * time.Second)
		that.collector.OnResponse(func(r *colly.Response) {
			that.parseProxy(r.Body)
		})
		that.collector.Visit(url)
	}
	fmt.Println(len(that.ProxyList.Proxies))
	that.verify()
	// for _, p := range that.ProxyList.Proxies {
	// 	fmt.Println(p.Url)
	// }
}

func (that *Proxy) parseProxy(body []byte) {
	r := string(body)
	if !strings.Contains(r, "vmess") {
		s, _ := base64.StdEncoding.DecodeString(r)
		r = string(s)
	}
	if strings.Contains(r, "vmess") {
		that.ProxyList.Date = that.date
		for _, p := range strings.Split(r, "\n") {
			pUrl := strings.Trim(p, " ")
			if _, ok := that.filter[pUrl]; !ok {
				that.filter[pUrl] = struct{}{}
				that.ProxyList.Proxies = append(that.ProxyList.Proxies, &config.Proxy{
					Url: pUrl,
				})
			}
		}
	}
}

func (that *Proxy) verify() {
	n := 0
	for _, p := range that.ProxyList.Proxies {
		if d, err := merkur.NewDialerByURI(p.Url); err == nil {
			c := d.ToHttpClient(10)
			_, err := c.Get("https://www.google.com")
			if err == nil {
				n++
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(n)
}

func (that *Proxy) Run() {
	that.GetProxyList()
}
