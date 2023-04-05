package vctrl

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

var AllowedSuffixes = []string{
	".zip",
	".tar.gz",
	".tar.bz2",
	".tar.xz",
}

type JDKPackage struct {
	Url      string
	FileName string
	OS       string
	Arch     string
	Size     string
	Checksum string
}

type JDKVersion struct {
	IsOfficial bool
	Versions   map[string][]*JDKPackage
	Doc        *goquery.Document
	Conf       *config.GVConfig
	c          *colly.Collector
	d          *downloader.Downloader
	dir        string
	env        *utils.EnvsHandler
}

func NewJDKVersion() (jv *JDKVersion) {
	jv = &JDKVersion{
		Versions: make(map[string][]*JDKPackage, 100),
		Conf:     config.New(),
		c:        colly.NewCollector(),
		d:        &downloader.Downloader{},
		env:      utils.NewEnvsHandler(),
	}
	jv.initeDirs()
	return
}

func (that *JDKVersion) ChooseResource() {
	fmt.Println("Choose a JDK download resource: ")
	fmt.Println("1) From injdk.cn (Faster in china. Default.)")
	fmt.Println("2) From oracle.com (Only latest versions are available.)")
	choice := "1"
	fmt.Scan(&choice)
	switch choice {
	case "2":
		that.IsOfficial = true
	default:
		that.IsOfficial = false
	}
}

func (that *JDKVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.DefaultJavaRoot); !ok {
		if err := os.MkdirAll(config.DefaultJavaRoot, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.JavaTarFilesPath); !ok {
		if err := os.MkdirAll(config.JavaTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.JavaUnTarFilesPath); !ok {
		if err := os.MkdirAll(config.JavaUnTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *JDKVersion) getDoc() {
	jUrl := that.Conf.Java.JDKUrl
	if that.IsOfficial {
		jUrl = that.Conf.Java.CompilerUrl
	}
	if !utils.VerifyUrls(jUrl) {
		return
	}
	that.c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	})
	that.c.Visit(jUrl)
}

func (that *JDKVersion) GetSha(sUrl string) (res string) {
	if that.IsOfficial {
		c := colly.NewCollector()
		c.OnResponse(func(r *colly.Response) {
			res = string(r.Body)
		})
		c.Visit(sUrl)
	}
	return
}

func (that *JDKVersion) GetFileSuffix(fName string) string {
	for _, k := range AllowedSuffixes {
		if strings.HasSuffix(fName, k) {
			return k
		}
	}
	return ""
}

func (that *JDKVersion) GetVersions() {
	if that.Doc == nil {
		that.getDoc()
	}
	if that.IsOfficial {
		that.Doc.Find("ul.rw-inpagetabs").First().Find("li").Each(func(i int, s *goquery.Selection) {
			v, _ := s.Find("a").Attr("href")
			sList := strings.Split(v, "java")
			vn := sList[len(sList)-1]
			that.Doc.Find(fmt.Sprintf("div#java%s", vn)).After("nav").Find("table").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					return
				}
				tArchive := strings.ToLower(s.Find("td").Eq(0).Text())
				tArchive = strings.ReplaceAll(tArchive, " ", "")
				tSize := s.Find("td").Eq(1).Text()
				tUrl, _ := s.Find("td").Eq(2).Find("a").Eq(0).Attr("href")
				tSha, _ := s.Find("td").Eq(2).Find("a").Eq(1).Attr("href")
				if strings.Contains(tArchive, Platform[runtime.GOARCH]) && strings.Contains(tArchive, "archive") {
					if !strings.Contains(tUrl, Platform[runtime.GOOS]) {
						return
					}
					p := &JDKPackage{}
					p.Arch = runtime.GOARCH
					p.OS = runtime.GOOS
					p.Size = tSize
					p.Url = tUrl
					if suffix := that.GetFileSuffix(p.Url); suffix != "" {
						p.FileName = fmt.Sprintf("jdk%s-%s_%s%s", vn, p.OS, p.Arch, suffix)
					} else {
						return
					}
					p.Checksum = that.GetSha(tSha)
					that.Versions[vn] = append(that.Versions[vn], p)
				}
			})
		})
	} else {
		that.Doc.Find("div#oracle-jdk").Find("div.col-sm-3").Each(func(i int, s *goquery.Selection) {
			vName := strings.ToLower(s.Find("span").Text())
			vName = strings.ReplaceAll(vName, "\n", "")
			vName = strings.ReplaceAll(vName, "\r", "")
			vName = strings.ReplaceAll(vName, " ", "")
			vName = strings.ReplaceAll(vName, "(lts)", "-lts")
			fmt.Println(vName)
			s.Find("li").Each(func(i int, ss *goquery.Selection) {
				if strings.Contains(vName, "jdk8") {
					return
				}
				p := &JDKPackage{}
				fileName := strings.ReplaceAll(strings.ToLower(ss.Find("a").Text()), " ", "")
				p.Arch = utils.ParseArch(fileName)
				p.OS = utils.ParsePlatform(fileName)
				if p.Arch == "" || p.OS == "" {
					return
				}
				if suffix := that.GetFileSuffix(fileName); suffix != "" {
					p.FileName = fmt.Sprintf("%s-%s_%s%s", vName, p.OS, p.Arch, suffix)
				} else {
					return
				}
				p.Url = strings.ReplaceAll(ss.Find("a").AttrOr("href", ""), " ", "")
				if p.Url == "" {
					return
				}
				that.Versions[vName] = append(that.Versions[vName], p)
				fmt.Println(p)
			})
		})

		that.Doc.Find("#Kona").Find("div.col-sm-3").Each(func(i int, s *goquery.Selection) {
			vName := strings.ToLower(s.Find("span").Text())
			vName = strings.ReplaceAll(vName, "\n", "")
			vName = strings.ReplaceAll(vName, "\r", "")
			vName = strings.ReplaceAll(vName, " ", "")
			vName = strings.ReplaceAll(vName, "(lts)", "-lts")
			fmt.Println(vName)
			s.Find("li").Each(func(i int, ss *goquery.Selection) {
				if !strings.Contains(vName, "jdk8") {
					return
				}
				p := &JDKPackage{}
				fileName := strings.ReplaceAll(strings.ToLower(ss.Find("a").Text()), " ", "")
				p.Arch = utils.ParseArch(fileName)
				p.OS = utils.ParsePlatform(fileName)
				if p.Arch == "" || p.OS == "" {
					return
				}
				if suffix := that.GetFileSuffix(fileName); suffix != "" {
					p.FileName = fmt.Sprintf("%s-%s_%s%s", vName, p.OS, p.Arch, suffix)
				} else {
					return
				}
				p.Url = strings.ReplaceAll(ss.Find("a").AttrOr("href", ""), " ", "")
				if p.Url == "" {
					return
				}
				that.Versions[vName] = append(that.Versions[vName], p)
				fmt.Println(p)
			})
		})
	}
}
