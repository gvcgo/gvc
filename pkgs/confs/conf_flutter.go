package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type FlutterConf struct {
	DefaultURLs  map[string]string `koanf:"default_urls"`
	OfficialURLs map[string]string `koanf:"official_urls"`
	TsingHuaUrl  string            `koanf:"tsing_hua_url"`
	NjuniUrl     string            `koanf:"njuni_url"`
	path         string
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
			tui.PrintError(err)
		}
	}
}

func (that *FlutterConf) Reset() {
	that.DefaultURLs = map[string]string{
		utils.Windows:      "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_windows.json",
		utils.MacOS:        "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_linux.json",
		utils.Linux:        "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_macos.json",
		"hosted_url":       "https://pub.flutter-io.cn",
		"storage_base_url": "https://storage.flutter-io.cn",
		"git_url":          "https://gitee.com/mirrors/Flutter.git",
	}

	that.OfficialURLs = map[string]string{
		utils.Windows:      "https://storage.googleapis.com/flutter_infra_release/releases/releases_windows.json",
		utils.MacOS:        "https://storage.googleapis.com/flutter_infra_release/releases/releases_macos.json",
		utils.Linux:        "https://storage.googleapis.com/flutter_infra_release/releases/releases_linux.json",
		"hosted_url":       "https://pub.dartlang.org",
		"storage_base_url": "https://storage.googleapis.com",
		"git_url":          "https://github.com/flutter/flutter.git",
	}
	that.TsingHuaUrl = "https://mirrors.cnnic.cn/flutter/flutter_infra_release/releases/"
	that.NjuniUrl = "https://mirrors.nju.edu.cn/flutter/flutter_infra_release/releases/"
}
