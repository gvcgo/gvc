package confs

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/moqsien/gvc/pkgs/utils"
)

func init() {
	if ok, _ := utils.PathIsExist(GVCWorkDir); !ok {
		if err := os.MkdirAll(GVCWorkDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", GVCWorkDir, err)
		}
	}
}

type GVConfig struct {
	Hosts    *HostsConf    `koanf:"hosts"`
	Go       *GoConf       `koanf:"go"`
	Java     *JavaConf     `koanf:"java"`
	Gradle   *GradleConf   `koanf:"gradle"`
	Maven    *MavenConf    `koanf:"maven"`
	Rust     *RustConf     `koanf:"rust"`
	Code     *CodeConf     `koanf:"code"`
	Nodejs   *NodejsConf   `koanf:"nodejs"`
	Python   *PyConf       `koanf:"python"`
	NVim     *NVimConf     `koanf:"nvim"`
	Proxy    *ProxyConf    `koanf:"proxy"`
	Github   *GithubConf   `koanf:"github"`
	Cygwin   *CygwinConf   `koanf:"cygwin"`
	Homebrew *HomebrewConf `koanf:"homebrew"`
	Vlang    *VlangConf    `koanf:"vlang"`
	Flutter  *FlutterConf  `koanf:"flutter"`
	Julia    *JuliaConf    `koanf:"julia"`
	w        *WebdavConf   `koanf:"webdav"`
	k        *koanf.Koanf
	parser   *yaml.YAML
	path     string
}

func New() (r *GVConfig) {
	r = &GVConfig{
		Hosts:    NewHostsConf(),
		Go:       NewGoConf(),
		Java:     NewJavaConf(),
		Gradle:   NewGradleConf(),
		Maven:    NewMavenConf(),
		Rust:     NewRustConf(),
		Code:     NewCodeConf(),
		Nodejs:   NewNodejsConf(),
		Python:   NewPyConf(),
		Proxy:    NewProxyConf(),
		Github:   NewGithubConf(),
		Cygwin:   NewCygwinConf(),
		Homebrew: NewHomebrewConf(),
		Vlang:    NewVlangConf(),
		Flutter:  NewFlutterConf(),
		Julia:    NewJuliaConf(),
		w:        NewWebdavConf(),
		k:        koanf.New("."),
		parser:   yaml.Parser(),
		path:     GVConfigPath,
	}
	r.initiate()
	return
}

func (that *GVConfig) initiate() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		that.w.Pull()
	}
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.Reload()
	} else {
		fmt.Println("[Cannot find default config files!]")
		fmt.Println("Do you want to use the default config files?[yes/N]")
		var r string
		fmt.Scan(&r)
		r = strings.ToLower(r)
		if r == "yes" || r == "y" {
			that.Reset()
		}
	}
}

func (that *GVConfig) Reset() {
	os.RemoveAll(GVCBackupDir)
	that.Hosts = NewHostsConf()
	that.Hosts.Reset()
	that.Go = NewGoConf()
	that.Go.Reset()
	that.Java = NewJavaConf()
	that.Java.Reset()
	that.Gradle = NewGradleConf()
	that.Gradle.Reset()
	that.Maven = NewMavenConf()
	that.Maven.Reset()
	that.Rust = NewRustConf()
	that.Rust.Reset()
	that.Code = NewCodeConf()
	that.Code.Reset()
	that.Nodejs = NewNodejsConf()
	that.Nodejs.Reset()
	that.Python = NewPyConf()
	that.Python.Reset()
	that.NVim = NewNVimConf()
	that.NVim.Reset()
	that.Proxy = NewProxyConf()
	that.Proxy.Reset()
	that.Github = NewGithubConf()
	that.Github.Reset()
	that.Cygwin = NewCygwinConf()
	that.Cygwin.Reset()
	that.Homebrew = NewHomebrewConf()
	that.Homebrew.Reset()
	that.Vlang = NewVlangConf()
	that.Vlang.Reset()
	that.Flutter = NewFlutterConf()
	that.Flutter.Reset()
	that.Julia = NewJuliaConf()
	that.Julia.Reset()
	that.Restore()
}

func (that *GVConfig) Reload() {
	err := that.k.Load(file.Provider(that.path), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
}

func (that *GVConfig) Restore() {
	if ok, _ := utils.PathIsExist(GVCBackupDir); !ok {
		os.MkdirAll(GVCBackupDir, os.ModePerm)
	}
	that.k.Load(structs.Provider(*that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.path, b, 0666)
	}
}

func (that *GVConfig) ShowPath() {
	fmt.Println("[GVC] config file path: ", that.path)
}

func (that *GVConfig) Pull() {
	that.w.Pull()
}

func (that *GVConfig) Push() {
	that.w.Push()
}

func (that *GVConfig) UseDefautFiles() {
	that.w.GetDefaultFiles()
}

func (that *GVConfig) SetupWebdav() {
	that.w.SetConf()
}

func (that *GVConfig) ShowDavConfigPath() {
	that.w.ShowDavConfigPath()
}
