package vproxy

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VmessList struct {
	Proxies []*Proxy `koanf:"proxies"`
	Date    string   `koanf:"date"`
	Total   int      `koanf:"count"`
	k       *koanf.Koanf
	parser  *yaml.YAML
	path    string
}

func NewVmessList(fileName ...string) (r *VmessList) {
	if len(fileName) == 0 {
		fmt.Println(">>>[FileName must be specified]")
		return
	}
	r = &VmessList{
		Proxies: make([]*Proxy, 0),
		k:       koanf.New("."),
		parser:  yaml.Parser(),
		path:    filepath.Join(config.ProxyFilesDir, fileName[0]),
	}
	return
}

func (that *VmessList) Today() string {
	return time.Now().Format("2006-01-02")
}

func (that *VmessList) GetDate() string {
	return that.Date
}

func (that *VmessList) CheckFilePath() (ok bool) {
	if ok, _ = utils.PathIsExist(that.path); !ok {
		fmt.Println("gvc[xray-core] is not ready, please check later.")
		fmt.Println("You can keep checking file existence: ", that.path)
	} else {
		fmt.Println("Find verified vmess list file: ", that.path)
	}
	return
}

func (that *VmessList) Reload() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		fmt.Println("ProxyList file does not exist: ", that.path)
		return
	}
	err := that.k.Load(file.Provider(that.path), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
}

func (that *VmessList) restore() {
	if ok, _ := utils.PathIsExist(config.ProxyFilesDir); !ok {
		os.MkdirAll(config.ProxyFilesDir, os.ModePerm)
	}
	that.k.Load(structs.Provider(*that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.path, b, 0666)
	}
}

func (that *VmessList) Update(proxies any) {
	pList, ok := proxies.([]*Proxy)
	if !ok {
		if rawProxyList, ok1 := proxies.([]RawProxy); ok1 {
			for _, p := range rawProxyList {
				if pxy, ok := p.(*Proxy); ok {
					pList = append(pList, pxy)
				}
			}
		} else {
			fmt.Println("Unsupported proxies.")
			return
		}
	}
	if len(pList) == 0 {
		fmt.Println("[Proxy List is empty]")
		return
	}
	that.Proxies = pList
	that.Date = that.Today()
	that.Total = len(that.Proxies)
	that.restore()
}

func (that *VmessList) Add(pxy *Proxy) {
	that.Reload()
	flag := false
	for _, p := range that.Proxies {
		if p.Uri == pxy.Uri {
			flag = true
			break
		}
	}
	if !flag {
		that.Date = that.Today()
		that.Proxies = append(that.Proxies, pxy)
		that.Total = len(that.Proxies)
		that.restore()
	}
}

func (that *VmessList) ShowProxyList() {
	that.Reload()

	color.Red.Println(that.Date)
	color.Yellow.Println(fmt.Sprintf("Total: %d", that.Total))
	for idx, p := range that.Proxies {
		rawUrl := p.GetUri()
		xo := &XrayVmessOutbound{}
		xo.ParseVmessUri(rawUrl)
		if xo.Address != "" {
			color.Cyan.Println(fmt.Sprintf("%d. %s:%d?path=%s (rtt: %dms)", idx, xo.Address, xo.Port, xo.Path, p.RTT))
		}
	}
}

func (that *VmessList) GetProxyList() []*Proxy {
	that.Reload()
	return that.Proxies
}

func (that *VmessList) ChooseFastest() *Proxy {
	that.Reload()
	if len(that.Proxies) == 0 {
		return nil
	}
	fastest := that.Proxies[0]
	for _, p := range that.Proxies {
		if p.RTT < fastest.RTT {
			fastest = p
		}
	}
	return fastest
}

func (that *VmessList) ChooseRandom() *Proxy {
	that.Reload()
	if len(that.Proxies) == 0 {
		return nil
	}
	r := rand.Intn(len(that.Proxies))
	return that.Proxies[r]
}

func (that *VmessList) ChooseByIndex(idx int) *Proxy {
	if idx < 0 || idx >= len(that.Proxies) {
		return that.ChooseRandom()
	}
	return that.Proxies[idx]
}
