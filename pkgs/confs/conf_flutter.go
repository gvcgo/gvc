package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type FlutterConf struct {
	DefaultURLs          map[string]string `koanf:"default_urls"`
	OfficialURLs         map[string]string `koanf:"official_urls"`
	TsingHuaUrl          string            `koanf:"tsing_hua_url"`
	NjuniUrl             string            `koanf:"njuni_url"`
	AndroidCMDTools      map[string]string `koanf:"android_cmd_tools"`
	AndroidCMDToolsUrlCN string            `koanf:"android_cmd_tools_cn_url"`
	AndroidCMDTooolsUrl  string            `koanf:"android_cmd_toools_url"`
	AndroidCN            string            `koanf:"android_cn_url"`
	Android              string            `koanf:"android_url"`
	path                 string
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *FlutterConf) Reset() {
	that.DefaultURLs = map[string]string{
		utils.Windows:      "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_windows.json",
		utils.MacOS:        "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_macos.json",
		utils.Linux:        "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_linux.json",
		"hosted_url":       "https://pub.flutter-io.cn",
		"storage_base_url": "https://storage.flutter-io.cn",
		"git_url":          "https://mirrors.tuna.tsinghua.edu.cn/git/flutter-sdk.git",
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
	that.AndroidCMDTools = map[string]string{
		utils.Windows: "https://googledownloads.cn/android/repository/commandlinetools-win-10406996_latest.zip",
		utils.MacOS:   "https://googledownloads.cn/android/repository/commandlinetools-mac-10406996_latest.zip",
		utils.Linux:   "https://googledownloads.cn/android/repository/commandlinetools-linux-10406996_latest.zip",
	}
	that.AndroidCMDToolsUrlCN = "https://googledownloads.cn/android/repository/"
	that.AndroidCMDTooolsUrl = "https://dl.google.com/android/repository/"
	that.AndroidCN = "https://developer.android.google.cn/studio?hl=zh-cn"
	that.Android = "https://developer.android.com/studio"
}
