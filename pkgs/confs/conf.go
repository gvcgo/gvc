package confs

import (
	"fmt"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/utils"
	xutils "github.com/moqsien/xtray/pkgs/utils"
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
	Typst    *TypstConf    `koanf:"typst"`
	Chatgpt  *ChatgptConf  `koanf:"chatgpt"`
	Webdav   *DavConf      `koanf:"dav"`
	path     string
	koanfer  *xutils.Koanfer
	// k        *koanf.Koanf
	// parser   *yaml.YAML
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
		Typst:    NewTypstConf(),
		Chatgpt:  NewGptConf(),
		Webdav:   NewDavConf(),
		path:     GVConfigPath,
		koanfer:  xutils.NewKoanfer(GVConfigPath),
		// k:        koanf.New("."),
		// parser:   yaml.Parser(),
	}
	r.initiate()
	return
}

func (that *GVConfig) initiate() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		that.SetDefault()
		that.Restore()
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

func (that *GVConfig) SetDefault() {
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
	that.Typst = NewTypstConf()
	that.Typst.Reset()
	that.Chatgpt = NewGptConf()
	that.Chatgpt.Reset()
	that.Webdav = NewDavConf()
	that.Webdav.Reset()
}

func (that *GVConfig) Reset() {
	os.RemoveAll(GVConfigPath)
	that.SetDefault()
	that.Restore()
}

func (that *GVConfig) Reload() {
	that.koanfer.Load(that)
}

func (that *GVConfig) Restore() {
	that.koanfer.Save(that)
}
