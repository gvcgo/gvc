package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GoConfig struct {
	CompilerUrls []string `koanf:"compiler_urls"`
	Proxies      []string `koanf:"proxies"`
}

type CodeConfig struct {
	StableUrl      string   `koanf:"stable_url"`
	CdnUrl         string   `koanf:"cdn_url"`
	DownloadUrl    string   `koanf:"download_url"`
	ExtIdentifiers []string `koanf:"ext_identifiers"`
}

type GVConfig struct {
	SourceUrls  []string    `koanf:"source_urls"`
	HostFilters []string    `koanf:"host_filters"`
	ReqTimeout  int         `koanf:"req_timeout"`
	MaxAvgRtt   int         `koanf:"max_avg_rtt"`
	PingCount   int         `koanf:"ping_count"`
	WorkerNum   int         `koanf:"worker_num"`
	Go          *GoConfig   `koanf:"go_config"`
	Code        *CodeConfig `koanf:"vscode"`
}

var (
	dConf = &GVConfig{
		SourceUrls: []string{
			"https://www.foul.trade:3000/Johy/Hosts/raw/branch/main/hosts.txt",
			"https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts",
			"https://raw.hellogithub.com/hosts",
		},
		HostFilters: []string{
			"github",
		},
		ReqTimeout: 30,
		MaxAvgRtt:  400,
		PingCount:  10,
		WorkerNum:  100,
		Go: &GoConfig{
			CompilerUrls: []string{
				"https://golang.google.cn/dl/",
				"https://go.dev/dl/",
				"https://studygolang.com/dl",
			},
			Proxies: []string{
				"https://goproxy.cn,direct",
				"https://repo.huaweicloud.com/repository/goproxy/,direct",
			},
		},
		Code: &CodeConfig{
			StableUrl:   "az764295.vo.msecnd.net",
			CdnUrl:      "vscode.cdn.azure.cn",
			DownloadUrl: "https://code.visualstudio.com/sha/download",
			ExtIdentifiers: []string{
				"moqsien.easynotes",
				"doggy8088.go-extension-pack",
				"galkowskit.go-interface-annotations",
				"liuchao.go-struct-tag",
				"tabnine.tabnine-vscode",
				"gruntfuggly.todo-tree",
				"zxh404.vscode-proto3",
				"premparihar.gotestexplorer",
				"ms-python.python",
				"ms-python.vscode-pylance",
				"donjayamanne.python-environment-manager",
				"alefragnani.project-manager",
				"yzhang.markdown-all-in-one",
				"mhutchie.git-graph",
				"asvetliakov.vscode-neovim",
				"ms-ceintl.vscode-language-pack-zh-hans",
				"bracketpaircolordlw.bracket-pair-color-dlw",
				"rust-lang.rust-analyzer",
				"vue.volar",
				"joe-re.sql-language-server",
				"akamud.vscode-theme-onedark",
				"pkief.material-icon-theme",
			},
		},
	}
)

type Conf struct {
	Config *GVConfig
	k      *koanf.Koanf
	parser *yaml.YAML
	path   string
}

func New() *Conf {
	c := &Conf{
		Config: new(GVConfig),
		k:      koanf.New("."),
		parser: yaml.Parser(),
		path:   DefaultConfigPath,
	}
	c.Initiate()
	return c
}

func (that *Conf) Initiate() {
	dir := filepath.Dir(that.path)
	if ok, _ := utils.PathIsExist(dir); !ok {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		} else {
			that.Config = dConf
			that.k.Load(structs.Provider(*that.Config, "koanf"), nil)
			if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
				os.WriteFile(that.path, b, 0666)
			}
		}
	} else if ok2, _ := utils.PathIsExist(that.path); ok2 {
		err := that.k.Load(file.Provider(that.path), that.parser)
		if err != nil {
			fmt.Println("[Config Load Failed] ", err)
			return
		}
		that.k.UnmarshalWithConf("", that.Config, koanf.UnmarshalConf{Tag: "koanf"})
	} else {
		that.Config = dConf
		that.k.Load(structs.Provider(*that.Config, "koanf"), nil)
		if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
			os.WriteFile(that.path, b, 0666)
		}
	}
}

func (that *Conf) Reset() {
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.Config = dConf
		that.k.Load(structs.Provider(*that.Config, "koanf"), nil)
		if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
			os.WriteFile(that.path, b, 0666)
		}
	}
}

func (that *Conf) Restore() {
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.k.Load(structs.Provider(*that.Config, "koanf"), nil)
		if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
			os.WriteFile(that.path, b, 0666)
		}
	}
}

func (that *Conf) ShowConfigFilePath() {
	fmt.Println("[GVC Config File] path: ", that.path)
}

func (that *Conf) GetHostsFilePath() string {
	if strings.Contains(runtime.GOOS, "window") {
		return HostFilePathForWin
	}
	return HostFilePathForNix
}
