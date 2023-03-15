package vctrl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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
	if runtime.GOOS == utils.Windows {
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
		if v.Version == "" {
			continue
		}
		p := &NodePackage{}
		p.VUrl, _ = url.JoinPath(that.Conf.Nodejs.ReleaseUrl, v.Version)
		p.Arch = runtime.GOARCH
		p.OS = runtime.GOOS
		that.Versions[v.Version] = p
		lts := that.parseLTS(v.Lts)
		if i == 0 || lts != "" {
			// Show only lts versions.
			r = append(r, v.Version)
			p.Lts = lts
			if lts != "" {
				v.Version = fmt.Sprintf("%s(%s)", v.Version, lts)
			}
		}
		p.FileName = fmt.Sprintf("nodejs%s-%s-%s%s",
			v.Version, p.OS, p.Arch, that.getSuffix())
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
	if runtime.GOOS != utils.Windows {
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
	var tarfile string
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile = that.download(version); tarfile != "" {
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
	vFilePath := filepath.Join(that.dir, "version.txt")
	if ok, _ := utils.PathIsExist(vFilePath); !ok {
		vFile, err := os.OpenFile(vFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("open file err = %v\n", err)
			return
		}
		defer vFile.Close()
		io.Copy(vFile, bytes.NewBuffer([]byte(version)))
	}
	that.setEnv(config.NodejsRoot)
	that.setNpm()
	fmt.Println("Use", version, "successed!")
}

func (that *NodeVersion) getCurrent() (v string) {
	vPath := filepath.Join(config.NodejsRoot, "version.txt")
	if ok, _ := utils.PathIsExist(vPath); ok {
		f, _ := os.Open(vPath)
		content, _ := io.ReadAll(f)
		v = string(content)
	}
	return
}

func (that *NodeVersion) ShowInstalled() {
	current := that.getCurrent()
	if rd, err := os.ReadDir(config.NodejsUntarFiles); err == nil {
		for _, v := range rd {
			if v.IsDir() {
				if current == v.Name() {
					s := fmt.Sprintf("%s <Current>", v.Name())
					color.Yellow.Println(s)
					continue
				}
				color.Cyan.Println(v.Name())
			}
		}
	}
}

func (that *NodeVersion) RemoveVersion(version string) {
	current := that.getCurrent()
	if version == "all" {
		if rd, err := os.ReadDir(config.NodejsUntarFiles); err == nil {
			for _, v := range rd {
				if v.IsDir() {
					if current == v.Name() {
						continue
					}
					os.RemoveAll(filepath.Join(config.NodejsUntarFiles, v.Name()))
				}
			}
		}
		if rd, err := os.ReadDir(config.NodejsTarFiles); err == nil {
			for _, v := range rd {
				if !v.IsDir() && !strings.Contains(v.Name(), current) {
					os.RemoveAll(filepath.Join(config.NodejsTarFiles, v.Name()))
				}
			}
		}
	} else if version != current {
		os.RemoveAll(filepath.Join(config.NodejsUntarFiles, version))
		if rd, err := os.ReadDir(config.NodejsTarFiles); err == nil {
			for _, v := range rd {
				if !v.IsDir() && strings.Contains(v.Name(), version) {
					os.RemoveAll(filepath.Join(config.NodejsTarFiles, v.Name()))
				}
			}
		}
	}
}

func (that *NodeVersion) setNpm() {
	var binPath string
	if runtime.GOOS == utils.Windows {
		binPath = filepath.Join(config.NodejsRoot, "npm")
	} else {
		binPath = filepath.Join(config.NodejsRoot, "bin/npm")
	}
	if ok, _ := utils.PathIsExist(binPath); ok {
		// npm config set registry=http://registry.npm.taobao.org
		utils.RunCommand(binPath, "config", "set", fmt.Sprintf("registry=%s", that.Conf.Nodejs.ProxyUrls[0]))
		utils.RunCommand(binPath, "config", "set", "prefix", config.NodejsGlobal)
		utils.RunCommand(binPath, "config", "set", "cache", config.NodejsCache)
	}
}
