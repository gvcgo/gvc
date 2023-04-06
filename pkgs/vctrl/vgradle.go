package vctrl

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
)

type GradlePackage struct {
	Version  string
	Url      string
	Checksum string
	FileName string
}

type GradleVersion struct {
	Versions map[string]*GradlePackage
	sha      map[string]string
	Doc      *goquery.Document
	Conf     *config.GVConfig
	d        *downloader.Downloader
	c        *colly.Collector
}

func NewGradleVersion() (gv *GradleVersion) {
	gv = &GradleVersion{
		Versions: make(map[string]*GradlePackage, 100),
		sha:      make(map[string]string, 100),
		Conf:     config.New(),
		d:        &downloader.Downloader{},
		c:        colly.NewCollector(),
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

func (that *GradleVersion) getDoc() {
	gUrl := that.Conf.Gradle.OfficialUrl
	if !utils.VerifyUrls(gUrl) {
		return
	}
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	})
	that.c.Visit(gUrl)
}

func (that *GradleVersion) getSha() {
	that.c = colly.NewCollector()
	cUrl := that.Conf.Gradle.OfficialCheckUrl
	if !utils.VerifyUrls(cUrl) {
		return
	}
	that.Doc = nil
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	})
	that.c.Visit(cUrl)
	if that.Doc != nil {
		that.Doc.Find("h3.u-text-with-icon").Each(func(i int, s *goquery.Selection) {
			version := s.Find("a").AttrOr("id", "")
			if version == "" {
				return
			}
			shaCode := s.Next().Find("li").Eq(0).Find("code").Text()
			if shaCode != "" {
				that.sha[version] = shaCode
			}
		})
	}
}

func (that *GradleVersion) shaCode(version string) (code string) {
	if len(that.sha) == 0 {
		that.getSha()
	}
	for k, v := range that.sha {
		if strings.Contains(k, version) {
			return v
		}
	}
	return
}

func (that *GradleVersion) getVersions() {
	if that.Doc == nil {
		that.getDoc()
	}
	if that.Doc != nil {
		that.Doc.Find("div.indent").Each(func(i int, s *goquery.Selection) {
			aLabel := s.Find("li").Eq(0).Find("a").Eq(0)
			p := &GradlePackage{}
			p.Url = aLabel.AttrOr("href", "")
			p.Version = aLabel.AttrOr("data-version", "")
			if p.Url == "" || p.Version == "" {
				return
			}
			p.Checksum = that.shaCode(p.Version)
			p.FileName = fmt.Sprintf("gradle-%s.zip", p.Version)
			that.Versions[p.Version] = p
		})
	}
}

func (that *GradleVersion) ShowVersions() {
	if len(that.Versions) == 0 {
		that.getVersions()
	}
	vList := []string{}
	for k := range that.Versions {
		vList = append(vList, k)
	}
	res := sorts.SortGoVersion(vList)
	fmt.Println(strings.Join(res, "  "))
}
