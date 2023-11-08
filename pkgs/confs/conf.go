package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/gvc/pkgs/utils"
)

func init() {
	utils.MakeDirs(GVCDir)
}

type GVConfig struct {
	GVCProxy *GVCReverseProxyConf `koanf:"gvc_proxy"`
	GVC      *GvcConf             `koanf:"gvc"`
	Hosts    *HostsConf           `koanf:"hosts"`
	Go       *GoConf              `koanf:"go"`
	Java     *JavaConf            `koanf:"java"`
	Gradle   *GradleConf          `koanf:"gradle"`
	Maven    *MavenConf           `koanf:"maven"`
	Rust     *RustConf            `koanf:"rust"`
	Code     *CodeConf            `koanf:"code"`
	Nodejs   *NodejsConf          `koanf:"nodejs"`
	Python   *PyConf              `koanf:"python"`
	NVim     *NVimConf            `koanf:"nvim"`
	NeoBox   *NeoboxConf          `koanf:"neobox"`
	Github   *GithubConf          `koanf:"github"`
	Cpp      *CppConf             `koanf:"cpp"`
	Homebrew *HomebrewConf        `koanf:"homebrew"`
	Vlang    *VlangConf           `koanf:"vlang"`
	Flutter  *FlutterConf         `koanf:"flutter"`
	Julia    *JuliaConf           `koanf:"julia"`
	Typst    *TypstConf           `koanf:"typst"`
	Webdav   *DavConf             `koanf:"dav"`
	Protobuf *ProtobufConf        `koanf:"protobuf"`
	GSudo    *GsudoConf           `koanf:"gsudo"`
	Docker   *DockerConf          `koanf:"docker"`
	GPT      *GPTConf             `koanf:"gpt"`
	path     string
	koanfer  *koanfer.JsonKoanfer
}

func New() (r *GVConfig) {
	kfer, _ := koanfer.NewKoanfer(GVConfigPath)
	r = &GVConfig{
		GVC:      NewGvcConf(),
		Hosts:    NewHostsConf(),
		Go:       NewGoConf(),
		Java:     NewJavaConf(),
		Gradle:   NewGradleConf(),
		Maven:    NewMavenConf(),
		Rust:     NewRustConf(),
		Code:     NewCodeConf(),
		Nodejs:   NewNodejsConf(),
		Python:   NewPyConf(),
		NeoBox:   NewNeoboxConf(),
		Github:   NewGithubConf(),
		Cpp:      NewCppConf(),
		Homebrew: NewHomebrewConf(),
		Vlang:    NewVlangConf(),
		Flutter:  NewFlutterConf(),
		Julia:    NewJuliaConf(),
		Typst:    NewTypstConf(),
		Webdav:   NewDavConf(),
		Protobuf: NewProtobuf(),
		GSudo:    NewGsudoConf(),
		Docker:   NewDockerConf(),
		GPT:      NewGPTConf(),
		path:     GVConfigPath,
		koanfer:  kfer,
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
		gprint.PrintWarning("Cannot find default config files.")
		cfm := confirm.NewConfirm(confirm.WithTitle("Use the default config files now?"))
		cfm.Run()
		if cfm.Result() {
			that.Reset()
		}
	}
}

func (that *GVConfig) SetDefault() {
	that.GVCProxy = NewReverseProxyConf()
	that.GVCProxy.Reset()
	that.GVC = NewGvcConf()
	that.GVC.Reset()
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
	that.NeoBox = NewNeoboxConf()
	that.NeoBox.Reset()
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
	that.Webdav = NewDavConf()
	that.Webdav.Reset()
	that.Protobuf = NewProtobuf()
	that.Protobuf.Reset()
	that.GSudo = NewGsudoConf()
	that.GSudo.Reset()
	that.Docker = NewDockerConf()
	that.Docker.Reset()
	that.GPT = NewGPTConf()
	that.GPT.Reset()
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
