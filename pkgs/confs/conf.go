package confs

import (
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	xutils "github.com/moqsien/xtray/pkgs/utils"
	"github.com/pterm/pterm"
)

func init() {
	utils.MakeDirs(GVCWorkDir)
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
	Xtray    *XtrayConf    `koanf:"xtray"`
	Github   *GithubConf   `koanf:"github"`
	Cpp      *CppConf      `koanf:"cpp"`
	Homebrew *HomebrewConf `koanf:"homebrew"`
	Vlang    *VlangConf    `koanf:"vlang"`
	Flutter  *FlutterConf  `koanf:"flutter"`
	Julia    *JuliaConf    `koanf:"julia"`
	Typst    *TypstConf    `koanf:"typst"`
	Chatgpt  *ChatgptConf  `koanf:"chatgpt"`
	Webdav   *DavConf      `koanf:"dav"`
	path     string
	koanfer  *xutils.Koanfer
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
		Xtray:    NewXtrayConf(),
		Github:   NewGithubConf(),
		Cpp:      NewCppConf(),
		Homebrew: NewHomebrewConf(),
		Vlang:    NewVlangConf(),
		Flutter:  NewFlutterConf(),
		Julia:    NewJuliaConf(),
		Typst:    NewTypstConf(),
		Chatgpt:  NewGptConf(),
		Webdav:   NewDavConf(),
		path:     GVConfigPath,
		koanfer:  xutils.NewKoanfer(GVConfigPath),
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
		tui.PrintWarning("Cannot find default config files.")
		if ok, _ := pterm.DefaultInteractiveConfirm.Show("Use the default config files now?"); ok {
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
	that.Xtray = NewXtrayConf()
	that.Xtray.Reset()
	that.Github = NewGithubConf()
	that.Github.Reset()
	that.Cpp = NewCppConf()
	that.Cpp.Reset()
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
