package vctrl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type NodePackage struct {
	Url      string
	VUrl     string
	FileName string
	Lts      string
	OS       string
	Arch     string
	Checksum string
}

var PlatformList map[string]string = map[string]string{
	"darwin":  "darwin",
	"windows": "win",
	"linux":   "linux",
	"amd64":   "x64",
	"arm64":   "arm64",
}

type nV struct {
	Version string `json:"version"`
	Lts     any    `json:"lts"`
	Date    string `json:"date"`
}

type NodeVersion struct {
	c        *colly.Collector
	d        *downloader.Downloader
	dir      string
	Versions map[string]*NodePackage
	vList    []*nV
	Conf     *config.GVConfig
	Doc      *goquery.Document
}

func NewNodeVersion() (nv *NodeVersion) {
	nv = &NodeVersion{
		Versions: make(map[string]*NodePackage, 50),
		Conf:     config.New(),
		vList:    []*nV{},
		c:        colly.NewCollector(),
		d:        &downloader.Downloader{},
	}
	nv.initeDirs()
	return
}

func (that *NodeVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.NodejsRoot); !ok {
		if err := os.MkdirAll(config.NodejsRoot, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.NodejsTarFiles); !ok {
		if err := os.MkdirAll(config.NodejsTarFiles, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.NodejsUntarFiles); !ok {
		if err := os.MkdirAll(config.NodejsUntarFiles, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *NodeVersion) getSuffix() string {
	suffix := ".tar.gz"
	if runtime.GOOS == "windows" {
		suffix = ".zip"
	}
	return suffix
}

func (that *NodeVersion) getVersions() (r []string) {
	if that.Conf.Nodejs.CompilerUrl != "" {
		that.c.OnResponse(func(r *colly.Response) {
			if err := json.Unmarshal(r.Body, &that.vList); err != nil {
				fmt.Println(err)
			}
		})
	}
	if _, err := url.Parse(that.Conf.Nodejs.CompilerUrl); err != nil {
		panic(err)
	}
	that.c.Visit(that.Conf.Nodejs.CompilerUrl)

	for i, v := range that.vList {
		lts := that.parseLTS(v.Lts)
		if i == 0 || lts != "" {
			r = append(r, v.Version)
			p := &NodePackage{}
			p.VUrl, _ = url.JoinPath(that.Conf.Nodejs.ReleaseUrl, v.Version)
			p.Arch = runtime.GOARCH
			p.OS = runtime.GOOS
			p.Lts = lts
			that.Versions[v.Version] = p
			if lts != "" {
				v.Version = fmt.Sprintf("%s(%s)", v.Version, lts)
			}
			p.FileName = fmt.Sprintf("nodejs%s-%s-%s%s",
				v.Version, p.OS, p.Arch, that.getSuffix())
		}
	}
	return
}

func (that *NodeVersion) parseLTS(v any) (r string) {
	if val, ok := v.(bool); ok {
		if val {
			r = "yes"
		}
		return
	}
	r, _ = v.(string)
	return
}

func (that *NodeVersion) ShowVersions() {
	fmt.Println(strings.Join(that.getVersions(), "  "))
}

func (that *NodeVersion) download(version string) string {
	if len(that.vList) == 0 {
		that.getVersions()
	}
	that.c = colly.NewCollector()
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	})

	if v, ok := that.Versions[version]; ok {
		that.c.Visit(v.VUrl)
		if that.Doc != nil {
			that.Doc.Find("a")
			that.Doc.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")

				if strings.Contains(href, PlatformList[runtime.GOOS]) && strings.Contains(href, PlatformList[runtime.GOARCH]) && strings.HasSuffix(href, that.getSuffix()) {
					v.Url, _ = url.JoinPath(v.VUrl, href)
				}
			})
		}
		if v.Url != "" {
			sumUrl, _ := url.JoinPath(v.VUrl, "SHASUMS256.txt")
			that.c = colly.NewCollector()
			that.c.OnResponse(func(r *colly.Response) {
				sumList := strings.Split(string(r.Body), "\n")
				nameList := strings.Split(v.Url, "/")
				for _, vl := range sumList {
					if strings.Contains(vl, nameList[len(nameList)-1]) {
						v.Checksum = strings.Trim(strings.Split(vl, " ")[0], " ")
					}
				}
			})
			that.c.Visit(sumUrl)
			if v.Checksum != "" {
				that.d.Url = v.Url
				that.d.Timeout = 100 * time.Minute
				fpath := filepath.Join(config.NodejsTarFiles, v.FileName)
				if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
					if ok := utils.CheckFile(fpath, "sha256", v.Checksum); ok {
						return fpath
					} else {
						os.RemoveAll(fpath)
					}
				}
			}
		}
	}
	return ""
}

func (that *NodeVersion) setEnv(nodeHome string) {
	if runtime.GOOS != "windows" {
		envar := fmt.Sprintf(config.NodejsEnvPattern, nodeHome)
		utils.SetUnixEnv(envar)
	} else {
		utils.SetWinEnv("NODE_HOME", nodeHome)
		utils.SetWinEnv("Path", nodeHome)
	}
}

func (that *NodeVersion) findDir(untarfile string) {
	if rd, err := os.ReadDir(untarfile); err == nil {
		for _, d := range rd {
			if d.IsDir() && d.Name() == "bin" {
				if ok, _ := utils.PathIsExist(filepath.Join(untarfile, "bin/node")); ok {
					that.dir = untarfile
				}
			} else if !d.IsDir() && d.Name() == "node.exe" {
				that.dir = untarfile
			} else if d.IsDir() {
				that.findDir(filepath.Join(untarfile, d.Name()))
			}
		}
	}
}

func (that *NodeVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.NodejsUntarFiles, version)
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.NodejsRoot); ok {
		os.RemoveAll(config.NodejsRoot)
	}
	that.findDir(untarfile)
	if that.dir == "" {
		return
	}
	if err := utils.MkSymLink(that.dir, config.NodejsRoot); err != nil {
		fmt.Println("[Create link failed] ", err)
		return
	}
	that.setEnv(config.NodejsRoot)
	fmt.Println("Use", version, "successed!")
}
