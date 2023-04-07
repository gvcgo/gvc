package vctrl

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gookit/color"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
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
	fmt.Println("1) From oracle.com (Only latest versions are available.)")
	fmt.Println("2) From injdk.cn (Old version are available.)")
	choice := "1"
	fmt.Scan(&choice)
	switch choice {
	case "1":
		that.IsOfficial = true
	default:
		that.IsOfficial = false
	}
}

func (that *JDKVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.DefaultJavaRoot); !ok {
		os.RemoveAll(config.DefaultJavaRoot)
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
				if !strings.Contains(tArchive, "archive") {
					return
				}
				p := &JDKPackage{}
				p.Arch = utils.ParseArch(tUrl)
				p.OS = utils.ParsePlatform(tUrl)
				if p.Arch == "" || p.OS == "" || tUrl == "" {
					return
				}
				p.Size = tSize
				p.Url = tUrl
				if suffix := that.GetFileSuffix(p.Url); suffix != "" {
					p.FileName = fmt.Sprintf("jdk%s-%s_%s%s", vn, p.OS, p.Arch, suffix)
				} else {
					return
				}
				p.Checksum = that.GetSha(tSha)
				key := fmt.Sprintf("jdk%s", vn)
				that.Versions[key] = append(that.Versions[key], p)
			})
		})
	} else {
		that.Doc.Find("div#oracle-jdk").Find("div.col-sm-3").Each(func(i int, s *goquery.Selection) {
			vName := strings.ToLower(s.Find("span").Text())
			vName = strings.ReplaceAll(vName, "\n", "")
			vName = strings.ReplaceAll(vName, "\r", "")
			vName = strings.ReplaceAll(vName, " ", "")
			vName = strings.ReplaceAll(vName, "(lts)", "-lts")
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
			})
		})

		that.Doc.Find("#Kona").Find("div.col-sm-3").Each(func(i int, s *goquery.Selection) {
			vName := strings.ToLower(s.Find("span").Text())
			vName = strings.ReplaceAll(vName, "\n", "")
			vName = strings.ReplaceAll(vName, "\r", "")
			vName = strings.ReplaceAll(vName, " ", "")
			vName = strings.ReplaceAll(vName, "(lts)", "-lts")
			s.Find("li").Each(func(i int, ss *goquery.Selection) {
				if !strings.Contains(vName, "jdk8") {
					return
				}
				if !strings.Contains(vName, "lts") {
					vName = fmt.Sprintf("%s-%s", vName, "lts")
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
			})
		})
	}
}

func (that *JDKVersion) ShowVersions() {
	that.GetVersions()
	vList := []string{}
	for k := range that.Versions {
		vList = append(vList, k)
	}
	vList = sorts.SortJDKVersion(vList)
	fmt.Println(strings.Join(vList, " "))
}

func (that *JDKVersion) findVersion(version string) (p *JDKPackage) {
	var pList []*JDKPackage
	for k, v := range that.Versions {
		if strings.Contains(k, version) {
			pList = v
			break
		}
	}
	if len(pList) > 0 {
		for _, p := range pList {
			if p.Arch == runtime.GOARCH && p.OS == runtime.GOOS {
				return p
			}
		}
	}
	return
}

func (that *JDKVersion) download(version string) (r string) {
	that.GetVersions()

	if p := that.findVersion(version); p != nil {
		that.d.Url = p.Url
		that.d.Timeout = 100 * time.Minute
		fpath := filepath.Join(config.JavaTarFilesPath, p.FileName)
		if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if p.Checksum != "" {
				if ok := utils.CheckFile(fpath, "sha256", p.Checksum); ok {
					return fpath
				} else {
					os.RemoveAll(fpath)
				}
			} else {
				return fpath
			}
		} else {
			os.RemoveAll(fpath)
		}
	} else {
		fmt.Println("Invalid jdk version. ", version)
	}
	return
}

func (that *JDKVersion) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		javaEnv := fmt.Sprintf(utils.JavaEnv, config.DefaultJavaRoot)
		that.env.UpdateSub(utils.SUB_JDK, javaEnv)
	} else {
		classPath := filepath.Join(config.DefaultJavaRoot, "lib")
		envList := map[string]string{
			"JAVA_HOME":  config.DefaultJavaRoot,
			"CLASS_PATH": filepath.Join(config.DefaultJavaRoot, "lib"),
			"PATH": fmt.Sprintf("%s;%s;%s", filepath.Join(config.DefaultJavaRoot, "bin"),
				filepath.Join(classPath, "tools.jar"), filepath.Join(classPath, "dt.jar")),
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *JDKVersion) findDir(untarfile string) {
	if rd, err := os.ReadDir(untarfile); err == nil {
		for _, d := range rd {
			if d.IsDir() && d.Name() == "bin" {
				that.dir = untarfile
			} else if d.IsDir() {
				that.findDir(filepath.Join(untarfile, d.Name()))
			}
		}
	}
}

func (that *JDKVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.JavaUnTarFilesPath, version)
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.DefaultJavaRoot); ok {
		os.RemoveAll(config.DefaultJavaRoot)
	}
	that.findDir(untarfile)
	if that.dir == "" {
		fmt.Println("[Can not find binaries] ", untarfile)
		return
	}

	if err := utils.MkSymLink(that.dir, config.DefaultJavaRoot); err != nil {
		fmt.Println("[Create link failed] ", err)
		return
	}
	if !that.env.DoesEnvExist(utils.SUB_JDK) {
		that.CheckAndInitEnv()
	}
	fmt.Println("Use", version, "succeeded!")
}

func (that *JDKVersion) getCurrent() (version string) {
	fpath := filepath.Join(config.DefaultJavaRoot, "release")
	content, _ := os.ReadFile(fpath)
	if len(content) == 0 {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "JAVA_VERSION=") {
			version = strings.ReplaceAll(strings.Split(line, "=")[1], `"`, "")
			version = strings.Split(version, ".")[0]
			version = fmt.Sprintf("jdk%s", version)
		}
	}
	return
}

func (that *JDKVersion) ShowInstalled() {
	current := that.getCurrent()
	dList, _ := os.ReadDir(config.JavaUnTarFilesPath)
	for _, d := range dList {
		if !strings.Contains(d.Name(), "jdk") {
			continue
		}
		if current == d.Name() {
			s := fmt.Sprintf("%s <Current>", d.Name())
			color.Yellow.Println(s)
			continue
		}
		color.Cyan.Println(d.Name())
	}
}

func (that *JDKVersion) removeTarFile(version string) {
	fNameStr := fmt.Sprintf("%s-%s_%s", version, runtime.GOOS, runtime.GOARCH)
	fNameStr1 := fmt.Sprintf("%s-lts-%s_%s", version, runtime.GOOS, runtime.GOARCH)
	dList, _ := os.ReadDir(config.JavaTarFilesPath)
	for _, d := range dList {
		if strings.Contains(d.Name(), fNameStr) || strings.Contains(d.Name(), fNameStr1) {
			os.RemoveAll(filepath.Join(config.JavaTarFilesPath, d.Name()))
		}
	}
}

func (that *JDKVersion) RemoveVersion(version string) {
	if !strings.HasPrefix(version, "jdk") {
		version = fmt.Sprintf("jdk%s", version)
	}
	current := that.getCurrent()
	if version != current {
		os.RemoveAll(filepath.Join(config.JavaUnTarFilesPath, version))
		that.removeTarFile(version)
	}
}

func (that *JDKVersion) RemoveUnused() {
	current := that.getCurrent()
	dList, _ := os.ReadDir(config.JavaUnTarFilesPath)
	for _, d := range dList {
		fmt.Println(d.Name())
		if current != d.Name() && strings.Contains(d.Name(), "jdk") {
			os.RemoveAll(filepath.Join(config.JavaUnTarFilesPath, d.Name()))
			that.removeTarFile(d.Name())
		}
	}
}
