package confs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/moqsien/gvc/pkgs/utils"
)

type ProxyCronConfig struct {
	Hours   int `koanf:"hours"`
	Minutes int `koanf:"minutes"`
}

type ProxyConf struct {
	SubUrls         []string         `koanf:"suburls"`
	VerifyUrl       string           `koanf:"verify_url"`
	InboundPort     int              `koanf:"inbound_port"`
	VerifyPortRange []int            `koanf:"verify_port_range"`
	Crontab         *ProxyCronConfig `koanf:"crontab"`
	GeoIpUrl        string           `koanf:"geo_ip_url"`
	SwitchOmegaUrl  string           `koanf:"switch_mega_url"`
	GithubDownload  []string         `koanf:"github_download"`
	MaxRTT          int              `koanf:"max_rtt"`
	PingPort        int              `koanf:"ping_port"`
	path            string
	k               *koanf.Koanf
	parser          *yaml.YAML
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
		"https://gitlab.com/mianfeifq/share/-/raw/master/data2023036.txt",
		"https://raw.fastgit.org/freefq/free/master/v2",
		"https://raw.githubusercontent.com/mfuu/v2ray/master/v2ray",
		"https://sub.nicevpn.top/long",
		"https://raw.githubusercontent.com/ermaozi/get_subscribe/main/subscribe/v2ray.txt",
		"https://raw.githubusercontent.com/tbbatbb/Proxy/master/dist/v2ray.config.txt",
		"https://raw.githubusercontent.com/vveg26/get_proxy/main/dist/v2ray.config.txt",
		"https://freefq.neocities.org/free.txt",
		"https://ghproxy.com/https://raw.githubusercontent.com/kxswa/k/k/base64",
	}
	that.VerifyUrl = "https://www.google.com"
	that.InboundPort = 2019
	that.VerifyPortRange = []int{2020, 2060}
	that.Crontab = &ProxyCronConfig{
		Hours:   1,
		Minutes: 30,
	}

	that.GeoIpUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/geoip.zip"
	that.SwitchOmegaUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/switch-omega.zip"
	that.GithubDownload = []string{
		"https://ghproxy.com/",
		"https://d.serctl.com/?dl_start",
	}
	that.MaxRTT = 3
	that.PingPort = 4156
}

func (that *ProxyConf) GetSubUrls() []string {
	for idx, url := range that.SubUrls {
		if strings.Contains(url, `%s`) {
			that.SubUrls[idx] = fmt.Sprintf(url, time.Now().Format("2006/01/20060102"))
		}
	}
	return that.SubUrls
}

func (that *ProxyConf) GetVerifyPorts() (result []int) {
	start, end := 2020, 2050
	if len(that.VerifyPortRange) == 1 {
		start, end = that.VerifyPortRange[0], that.VerifyPortRange[0]
	} else if len(that.VerifyPortRange) == 2 {
		start, end = func(input []int) (int, int) {
			if input[0] > input[1] {
				return input[1], input[0]
			}
			return input[0], input[1]
		}(that.VerifyPortRange)
	}
	for i := start; i < end; i++ {
		result = append(result, i)
	}
	return
}

func (that *ProxyConf) GetCrontabStr() (r string) {
	if that.Crontab != nil && (that.Crontab.Hours != 0 || that.Crontab.Minutes != 0) {
		r = fmt.Sprintf("@every %vh%vm", that.Crontab.Hours, that.Crontab.Minutes)
	} else {
		r = "@every 2h"
	}
	return
}
