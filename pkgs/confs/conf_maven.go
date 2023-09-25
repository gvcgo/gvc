package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type MavenConf struct {
	ApacheUrl3    string `koanf:"apache_url3"`
	ApacheUrl4    string `koanf:"apache_url4"`
	UrlPattern    string `koanf:"url_pattern"`
	ShaUrlPattern string `koanf:"sha_url_pattern"`
	path          string
}

func NewMavenConf() (r *MavenConf) {
	r = &MavenConf{
		path: MavenRoot,
	}
	r.setup()
	return
}

func (that *MavenConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *MavenConf) Reset() {
	that.ApacheUrl3 = "https://dlcdn.apache.org/maven/maven-3/"
	that.ApacheUrl4 = "https://dlcdn.apache.org/maven/maven-4/"
	that.UrlPattern = "%s%s/binaries/apache-maven-%s-bin.tar.gz"
	that.ShaUrlPattern = "%s%s/binaries/apache-maven-%s-bin.tar.gz.sha512"
}
