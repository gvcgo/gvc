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
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
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
	fetcher  *request.Fetcher
	dir      string
	env      *utils.EnvsHandler
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
		fetcher:  request.NewFetcher(),
		env:      utils.NewEnvsHandler(),
	}
	nv.initeDirs()
	nv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *NodeVersion) initeDirs() {
	utils.MakeDirs(config.NodejsRoot, config.NodejsTarFiles, config.NodejsUntarFiles)
}

func (that *NodeVersion) getSuffix() string {
	suffix := ".tar.gz"
	if runtime.GOOS == utils.Windows {
		suffix = ".zip"
	}
	return suffix
}

func (that *NodeVersion) getVersions() (r []string) {
	that.fetcher.Url = that.Conf.Nodejs.CompilerUrl
	if _, err := url.Parse(that.fetcher.Url); err != nil {
		gprint.PrintError("%+v", err)
		os.Exit(1)
	}

	if resp := that.fetcher.Get(); resp != nil {
		content, _ := io.ReadAll(resp.RawBody())
		if err := json.Unmarshal(content, &that.vList); err != nil {
			gprint.PrintError(fmt.Sprintf("Parse content from %s failed.", that.fetcher.Url))
			return
		}
	}

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
	fc := gprint.NewFadeColors(that.getVersions())
	fc.Println()
}

func (that *NodeVersion) download(version string) string {
	if len(that.vList) == 0 {
		that.getVersions()
	}
	if v, ok := that.Versions[version]; ok {
		that.fetcher.Url = v.VUrl
		if that.fetcher.Url == "" {
			return ""
		}
		that.Doc = nil
		if resp := that.fetcher.Get(); resp != nil {
			that.Doc, _ = goquery.NewDocumentFromReader(resp.RawBody())
		}
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
			that.fetcher.Url, _ = url.JoinPath(v.VUrl, "SHASUMS256.txt")
			if resp := that.fetcher.Get(); resp != nil {
				content, _ := io.ReadAll(resp.RawBody())
				sumList := strings.Split(string(content), "\n")
				nameList := strings.Split(v.Url, "/")
				for _, vl := range sumList {
					if strings.Contains(vl, nameList[len(nameList)-1]) {
						v.Checksum = strings.Trim(strings.Split(vl, " ")[0], " ")
					}
				}
			}

			if v.Checksum != "" {
				that.fetcher.Url = v.Url
				that.fetcher.Timeout = 100 * time.Minute
				fpath := filepath.Join(config.NodejsTarFiles, v.FileName)
				that.fetcher.SetThreadNum(8)
				if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
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
		nodeEnv := fmt.Sprintf(utils.NodeEnv, nodeHome)
		that.env.UpdateSub(utils.SUB_NODE, nodeEnv)
	} else {
		envList := map[string]string{
			"NODE_HOME": nodeHome,
			"PATH":      nodeHome,
		}
		that.env.SetEnvForWin(envList)
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
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
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
		gprint.PrintError(fmt.Sprintf("Create link failed: %+v", err))
		return
	}
	vFilePath := filepath.Join(that.dir, "version.txt")
	if ok, _ := utils.PathIsExist(vFilePath); !ok {
		vFile, err := os.OpenFile(vFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			gprint.PrintError(fmt.Sprintf("Open file failed: %+v", err))
			return
		}
		defer vFile.Close()
		io.Copy(vFile, bytes.NewBuffer([]byte(version)))
	}
	that.setEnv(config.NodejsRoot)
	that.setNpm()
	gprint.PrintSuccess(fmt.Sprintf("Use %s succeeded!", version))
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
					gprint.Yellow("%s <Current>", v.Name())
					continue
				}
				gprint.Cyan(v.Name())
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
