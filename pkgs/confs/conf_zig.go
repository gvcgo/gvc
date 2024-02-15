package confs

import (
	"os"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/pkgs/utils"
)

type ZigConf struct {
	ZigDownloadUrl string `json,koanf:"zig_download_url"`
	ZlsDownloadUrl string `json,koanf:"zls_download_url"`
	path           string
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
	that.ZlsDownloadUrl = "https://github.com/zigtools/zls/releases/latest/"
}
