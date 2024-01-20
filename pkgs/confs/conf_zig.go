package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type ZigConf struct {
	ZigDownloadUrl  string            `json,koanf:"zig_download_url"`
	ZlsDownloadUrls map[string]string `json,koanf:"zls_download_urls"`
	path            string
}

func NewZigConf() (z *ZigConf) {
	z = &ZigConf{
		path: ZigFilesDir,
	}
	z.setup()
	return
}

func (that *ZigConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *ZigConf) Reset() {
	that.ZigDownloadUrl = "https://ziglang.org/download/"
	that.ZlsDownloadUrls = map[string]string{
		"linux_amd64":  "https://github.com/zigtools/zls/releases/latest/download/zls-x86_64-linux.tar.gz",
		"linux_arm64":  "https://github.com/zigtools/zls/releases/latest/download/zls-aarch64-linux.tar.gz",
		"darwin_amd64": "https://github.com/zigtools/zls/releases/latest/download/zls-x86_64-macos.tar.gz",
		"darwin_arm64": "https://github.com/zigtools/zls/releases/latest/download/zls-aarch64-macos.tar.gz",
		"windows":      "https://github.com/zigtools/zls/releases/latest/download/zls-x86_64-windows.zip",
	}
}
