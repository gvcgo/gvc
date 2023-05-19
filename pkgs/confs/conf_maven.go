package confs

import (
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
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
			tui.PrintError(err)
		}
	}
}

func (that *MavenConf) Reset() {
	that.ApacheUrl3 = "https://dlcdn.apache.org/maven/maven-3/"
	that.ApacheUrl4 = "https://dlcdn.apache.org/maven/maven-4/"
	that.UrlPattern = "%s%s/binaries/apache-maven-%s-bin.tar.gz"
	that.ShaUrlPattern = "%s%s/binaries/apache-maven-%s-bin.tar.gz.sha512"
}
