package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GoInstalledVersion struct {
	Current   string   `koanf:"current"`
	Installed []string `koanf:"installed"`
}

type GV struct {
	Versions *GoInstalledVersion
	k        *koanf.Koanf
	parser   *yaml.YAML
	path     string
}

func NewGoVersion() *GV {
	return &GV{
		Versions: &GoInstalledVersion{},
		k:        koanf.New("."),
		parser:   yaml.Parser(),
		path:     GoInstalled,
	}
}

func (that *GV) Initiate() {
	dir := filepath.Dir(that.path)
	if ok, _ := utils.PahtIsExist(dir); !ok {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		} else {
			that.Versions = &GoInstalledVersion{}
			that.k.Load(structs.Provider(*that.Versions, "koanf"), nil)
			if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
				os.WriteFile(that.path, b, 0666)
			}
		}
	} else if ok2, _ := utils.PahtIsExist(that.path); ok2 {
		err := that.k.Load(file.Provider(that.path), that.parser)
		if err != nil {
			fmt.Println("[Config Load Failed] ", err)
			return
		}
		that.k.UnmarshalWithConf("", that.Versions, koanf.UnmarshalConf{Tag: "koanf"})
	} else {
		that.Versions = &GoInstalledVersion{}
		that.k.Load(structs.Provider(*that.Versions, "koanf"), nil)
		if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
			os.WriteFile(that.path, b, 0666)
		}
	}
	if ok, _ := utils.PahtIsExist(GoTarFilesPath); !ok {
		if err := os.Mkdir(GoTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PahtIsExist(GoUnTarFilesPath); !ok {
		if err := os.Mkdir(GoUnTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}
