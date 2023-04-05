package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type FlutterConf struct {
	FlutterCNUrls map[string]string `koanf:"flutter_cn_urls"`
	FlutterENUrls map[string]string `koanf:"flutter_en_urls"`
	TsingHuaUrl   string            `koanf:"tsing_hua_url"`
	NjuniUrl      string            `koanf:"njuni_url"`
	path          string
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
	that.FlutterCNUrls = map[string]string{
		utils.MacOS:   "https://flutter.cn/docs/get-started/install/macos",
		utils.Linux:   "https://flutter.cn/docs/get-started/install/linux",
		utils.Windows: "https://flutter.cn/docs/get-started/install/windows",
	}

	that.FlutterENUrls = map[string]string{
		utils.MacOS:   "https://docs.flutter.dev/get-started/install/macos",
		utils.Linux:   "https://docs.flutter.dev/get-started/install/linux",
		utils.Windows: "https://docs.flutter.dev/get-started/install/windows",
	}

	that.TsingHuaUrl = "https://mirrors.cnnic.cn/flutter/flutter_infra_release/releases/"
	that.NjuniUrl = "https://mirrors.nju.edu.cn/flutter/flutter_infra_release/releases/"
}
