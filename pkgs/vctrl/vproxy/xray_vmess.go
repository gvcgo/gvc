package vproxy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VmessProxy struct {
	Uri string `koanf:"uri"`
	RTT int    `koanf:"rtt"`
}

type VmessList struct {
	Proxies []*VmessProxy `koanf:"proxies"`
	Date    string        `koanf:"date"`
	k       *koanf.Koanf
	parser  *yaml.YAML
	path    string
}

func NewVmessList() (r *VmessList) {
	r = &VmessList{
		Proxies: make([]*VmessProxy, 2000),
		k:       koanf.New("."),
		parser:  yaml.Parser(),
		path:    filepath.Join(config.ProxyFilesDir, "proxies-raw-vmess.yml"),
	}
	return
}

func (that *VmessList) Today() string {
	return time.Now().Format("2006-01-02")
}

func (that *VmessList) GetDate() string {
	return that.Date
}

func (that *VmessList) Reload() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		fmt.Println("ProxyList file does not exist.")
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
	pList, ok := proxies.([]*VmessProxy)
	if !ok {
		fmt.Println("Unsupported proxies.")
		return
	}
	if len(pList) == 0 {
		fmt.Println("[Proxy List is empty]")
		return
	}
	that.Proxies = pList
	that.Date = that.Today()
	that.restore()
}
