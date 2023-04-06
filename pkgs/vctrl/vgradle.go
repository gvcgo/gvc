package vctrl

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GradlePackage struct {
	Version  string
	Url      string
	Checksum string
}

type GradleVersion struct {
	Versions map[string]*GradlePackage
	Doc      *goquery.Document
	Conf     *config.GVConfig
	d        *downloader.Downloader
}

func NewGradleVersion() (gv *GradleVersion) {
	gv = &GradleVersion{
		Versions: make(map[string]*GradlePackage, 100),
		Conf:     config.New(),
		d:        &downloader.Downloader{},
	}
	gv.initeDirs()
	return gv
}

func (that *GradleVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.GradleRoot); !ok {
		os.RemoveAll(config.GradleRoot)
		if err := os.MkdirAll(config.GradleRoot, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleTarFilePath); !ok {
		if err := os.MkdirAll(config.GradleTarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleUntarFilePath); !ok {
		if err := os.MkdirAll(config.GradleUntarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}
