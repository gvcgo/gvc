package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type FlutterConf struct {
	WinUrl         string `koanf:"win_url"`
	LinuxUrl       string `koanf:"linux_url"`
	MacosUrl       string `koanf:"macos_url"`
	TsingHuaUrl    string `koanf:"tsing_hua_url"`
	NjuniUrl       string `koanf:"njuni_url"`
	HostedUrl      string `koanf:"hosted_url"`
	StorageBaseUrl string `koanf:"storage_base_url"`
	GitUrl         string `koanf:"git_url"`
	path           string
}

func NewFlutterConf() (r *FlutterConf) {
	r = &FlutterConf{
		path: FlutterFilesDir,
	}
	r.setup()
	return r
}

func (that *FlutterConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *FlutterConf) Reset() {
	that.WinUrl = "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_windows.json"
	that.LinuxUrl = "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_linux.json"
	that.MacosUrl = "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_macos.json"
	that.TsingHuaUrl = "https://mirrors.cnnic.cn/flutter/flutter_infra_release/releases/"
	that.NjuniUrl = "https://mirrors.nju.edu.cn/flutter/flutter_infra_release/releases/"
	that.HostedUrl = "https://pub.flutter-io.cn"
	that.StorageBaseUrl = "https://storage.flutter-io.cn"
	that.GitUrl = "https://gitee.com/mirrors/Flutter.git"
}
