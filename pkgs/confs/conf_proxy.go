package confs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Proxy struct {
	Url string `koanf:"url"`
	RTT int64  `koanf:"rtt"`
}

type ProxyList struct {
	Date    string   `koanf:"date"`
	Proxies []*Proxy `koanf:"proxies"`
}

type ProxyConf struct {
	SubUrls []string `koanf:"suburls"`
	path    string
	k       *koanf.Koanf
	parser  *yaml.YAML
}

func NewProxyConf() (r *ProxyConf) {
	r = &ProxyConf{
		path:   ProxyFilesDir,
		k:      koanf.New("."),
		parser: yaml.Parser(),
	}
	r.setup()
	return r
}

func (that *ProxyConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *ProxyConf) Reset() {
	that.SubUrls = []string{
		`https://clashnode.com/wp-content/uploads/%s.txt`,
		`https://nodefree.org/dy/%s.txt`,
		"https://raw.fastgit.org/freefq/free/master/v2",
		"https://raw.githubusercontent.com/mfuu/v2ray/master/v2ray",
		"https://sub.nicevpn.top/long",
		"https://raw.githubusercontent.com/ermaozi/get_subscribe/main/subscribe/v2ray.txt",
		"https://raw.githubusercontent.com/tbbatbb/Proxy/master/dist/v2ray.config.txt",
		"https://raw.githubusercontent.com/vveg26/get_proxy/main/dist/v2ray.config.txt",
		"https://freefq.neocities.org/free.txt",
		"https://ghproxy.com/https://raw.githubusercontent.com/kxswa/k/k/base64",
	}
}

func (that *ProxyConf) GetSubUrls() []string {
	for idx, url := range that.SubUrls {
		if strings.Contains(url, `%s`) {
			that.SubUrls[idx] = fmt.Sprintf(url, time.Now().Format("2006/01/20060102"))
		}
	}
	return that.SubUrls
}

func (that *ProxyConf) LoadProxies() (r *ProxyList) {
	err := that.k.Load(file.Provider(that.path), that.parser)
	if err != nil {
		fmt.Println("[Proxies Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", r, koanf.UnmarshalConf{Tag: "koanf"})
	return
}

func (that *ProxyConf) RestoreProxies(p *ProxyList) {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		os.MkdirAll(that.path, os.ModePerm)
	}
	that.k.Load(structs.Provider(*p, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.path, b, 0666)
	}
}
