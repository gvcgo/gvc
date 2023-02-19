package confs

import (
	"fmt"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/moqsien/gvc/pkgs/utils"
)

func init() {
	if ok, _ := utils.PathIsExist(GVCWorkDir); !ok {
		if err := os.MkdirAll(GVCWorkDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", GVCWorkDir, err)
		}
	}
}

type GVConfig struct {
	Hosts  *HostsConf  `koanf:"hosts"`
	Go     *GoConf     `koanf:"go"`
	Code   *CodeConf   `koanf:"code"`
	w      *WebdavConf `koanf:"webdav"`
	k      *koanf.Koanf
	parser *yaml.YAML
	path   string
}

func New() (r *GVConfig) {
	r = &GVConfig{
		Hosts:  NewHostsConf(),
		Go:     NewGoConf(),
		Code:   NewCodeConf(),
		w:      NewWebdavConf(),
		k:      koanf.New("."),
		parser: yaml.Parser(),
		path:   GVConfigPath,
	}
	r.initiate()
	return
}

func (that *GVConfig) initiate() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		that.w.Pull()
	}
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.Reload()
	} else {
		fmt.Println("[Cannot find default config files]")
	}
}

func (that *GVConfig) Reset() {
	os.RemoveAll(GVCBackupDir)
	that.Hosts = NewHostsConf()
	that.Hosts.Reset()
	that.Go = NewGoConf()
	that.Go.Reset()
	that.Code = NewCodeConf()
	that.Code.Reset()
	that.Restore()
}

func (that *GVConfig) Reload() {
	err := that.k.Load(file.Provider(that.path), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
}

func (that *GVConfig) Restore() {
	if ok, _ := utils.PathIsExist(GVCBackupDir); !ok {
		os.MkdirAll(GVCBackupDir, os.ModePerm)
	}
	that.k.Load(structs.Provider(*that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.path, b, 0666)
	}
}

func (that *GVConfig) ShowPath() {
	fmt.Println("[GVC] config file path: ", that.path)
}

func (that *GVConfig) Pull() {
	that.w.Pull()
}

func (that *GVConfig) Push() {
	that.w.Push()
}

func (that *GVConfig) UseDefautFiles() {
	that.w.GetDefaultFiles()
}

func (that *GVConfig) SetupWebdav() {
	that.w.SetConf()
}

func (that *GVConfig) ShowDavConfigPath() {
	that.w.ShowDavConfigPath()
}
