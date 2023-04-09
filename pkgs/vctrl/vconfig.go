package vctrl

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/studio-b12/gowebdav"
)

type WebdavConf struct {
	Host            string `koanf:"url"`
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	RemoteDir       string `koanf:"remote_dir"`
	LocalDir        string `koanf:"local_dir"`
	DefaultFilesUrl string `koanf:"default_files"`
}

type GVCWebdav struct {
	DavConf *WebdavConf
	Conf    *config.GVConfig
	k       *koanf.Koanf
	parser  *yaml.YAML
	client  *gowebdav.Client
}

func NewGVCWebdav() (gw *GVCWebdav) {
	gw = &GVCWebdav{
		Conf:   config.New(),
		k:      koanf.New("."),
		parser: yaml.Parser(),
	}
	return
}
