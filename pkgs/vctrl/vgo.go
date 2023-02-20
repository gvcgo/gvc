package vctrl

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GoPackage struct {
	Url       string
	FileName  string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	CheckType string
}

type GoVersion struct {
	*downloader.Downloader
	Versions  map[string][]*GoPackage
	Doc       *goquery.Document
	Conf      *config.GVConfig
	ParsedUrl *url.URL
}

func NewGoVersion() (gv *GoVersion) {
	gv = &GoVersion{
		Versions:   make(map[string][]*GoPackage, 50),
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
	}
	gv.initeDirs()
	return
}

func (that *GoVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.DefaultGoRoot); !ok {
		if err := os.MkdirAll(config.DefaultGoRoot, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GoTarFilesPath); !ok {
		if err := os.MkdirAll(config.GoTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GoUnTarFilesPath); !ok {
		if err := os.MkdirAll(config.GoUnTarFilesPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *GoVersion) getDoc() {
	if len(that.Conf.Go.CompilerUrls) > 0 {
		that.Url = that.Conf.Go.CompilerUrls[0]
		var err error
		if that.ParsedUrl, err = url.Parse(that.Url); err != nil {
			panic(err)
		}
		that.Timeout = 30 * time.Second
		if resp := that.GetUrl(); resp != nil {
			var err error
			that.Doc, err = goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				fmt.Println("[parse page errored] ", err)
			}
			c := []byte{}
			resp.Body.Read(c)
			resp.Body.Close()
		}
	}
}

func (that *GoVersion) findPackages(table *goquery.Selection) (pkgs []*GoPackage) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		if strings.HasPrefix(href, "/") { // relative paths
			href = fmt.Sprintf("%s://%s%s", that.ParsedUrl.Scheme, that.ParsedUrl.Host, href)
		}
		pkgs = append(pkgs, &GoPackage{
			FileName:  td.Eq(0).Find("a").Text(),
			Url:       href,
			Kind:      strings.ToLower(td.Eq(1).Text()),
			OS:        utils.MapArchAndOS(td.Eq(2).Text()),
			Arch:      utils.MapArchAndOS(td.Eq(3).Text()),
			Size:      td.Eq(4).Text(),
			Checksum:  td.Eq(5).Text(),
			CheckType: alg,
		})
	})
	return pkgs
}

func (that *GoVersion) hasUnstableVersions() bool {
	if that.Doc == nil {
		that.getDoc()
	}
	label := that.Doc.Find("#unstable")
	if label == nil {
		return false
	}
	return that.Doc.Find("#unstable").Length() > 0
}

func (that *GoVersion) StableVersions() (err error) {
	var divs *goquery.Selection
	if that.hasUnstableVersions() {
		divs = that.Doc.Find("#stable").NextUntil("#unstable")
	} else {
		divs = that.Doc.Find("#stable").NextUntil("#archive")
	}
	divs.Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) UnstableVersions() (err error) {
	that.Doc.Find("#unstable").NextUntil("#archive").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) ArchivedVersions() (err error) {
	that.Doc.Find("#archive").Find("div.toggle").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) AllVersions() (err error) {
	if that.Doc == nil {
		that.getDoc()
	}
	err = that.StableVersions()
	if err != nil {
		return
	}
	err = that.ArchivedVersions()
	if err != nil {
		return
	}
	err = that.UnstableVersions()
	if err != nil {
		return
	}
	return
}

func (that *GoVersion) GetVersions() (vList []string) {
	for v := range that.Versions {
		vList = append(vList, v)
	}
	return vList
}

const (
	ShowAll      string = "1"
	ShowStable   string = "2"
	ShowUnstable string = "3"
)

func (that *GoVersion) ShowRemoteVersions(arg string) {
	var v *utils.VComparator
	if that.Doc == nil {
		that.getDoc()
	}
	switch arg {
	case ShowAll:
		if err := that.AllVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	case ShowStable:
		if err := that.StableVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	case ShowUnstable:
		if err := that.UnstableVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	default:
		fmt.Println("[Unknown show type] ", arg)
	}
}

func (that *GoVersion) findPackage(version string, kind ...string) (p *GoPackage) {
	k := "archive"
	if len(kind) > 0 {
		k = kind[0]
	}
	that.AllVersions()
	if vList, ok := that.Versions[version]; ok {
		for _, v := range vList {
			if v.OS == runtime.GOOS && v.Arch == runtime.GOARCH && v.Kind == k {
				p = v
			}
		}
	}
	return
}

func (that *GoVersion) download(version string) (r string) {
	p := that.findPackage(version)
	if p != nil {
		fName := fmt.Sprintf("go-%s-%s.%s%s", version, p.OS, p.Arch, utils.GetExt(p.FileName))
		fpath := filepath.Join(config.GoTarFilesPath, fName)
		that.Url = p.Url
		that.Timeout = 180 * time.Second
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if ok := that.checkFile(p, fpath); ok {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	} else {
		fmt.Println("Cannot find version:", version, ".")
	}
	return
}

func (that *GoVersion) checkFile(p *GoPackage, fpath string) (r bool) {
	return utils.CheckFile(fpath, p.CheckType, p.Checksum)
}

func (that *GoVersion) CheckAndInitEnv() {
	if runtime.GOOS != "windows" {
		envar := fmt.Sprintf(config.GoUnixEnv,
			that.Conf.Go.Proxies[0],
			fmt.Sprintf("%s:%s:$PATH", "$GOPATH/bin", "$GOROOT/bin"))
		utils.SetUnixEnv(envar)
	} else {
		envarList := map[string]string{
			"GOROOT":  config.DefaultGoRoot,
			"GOPATH":  config.DefaultGoPath,
			"GOBIN":   filepath.Join(config.DefaultGoPath, "bin"),
			"GOPROXY": that.Conf.Go.Proxies[0],
		}
		for key, value := range envarList {
			utils.SetWinEnv(key, value)
		}
		ePath := fmt.Sprintf("%s;%s",
			filepath.Join(config.DefaultGoPath, "bin"),
			filepath.Join(config.DefaultGoRoot, "bin"))
		utils.SetWinEnv("PATH", ePath)
		fmt.Println("set go envs successed!")
	}
	if ok, _ := utils.PathIsExist(config.DefaultGoPath); !ok {
		os.MkdirAll(config.DefaultGoPath, os.ModePerm)
	}
}

func (that *GoVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.GoUnTarFilesPath, version)
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.DefaultGoRoot); ok {
		os.RemoveAll(config.DefaultGoRoot)
	}
	if err := utils.MkSymLink(filepath.Join(untarfile, "go"), config.DefaultGoRoot); err != nil {
		fmt.Println("[Create link failed] ", err)
		return
	}
	if !that.isGoEnvsAvailable() {
		that.CheckAndInitEnv()
	}
	fmt.Println("Use", version, "successed!")
}

func (that *GoVersion) isGoEnvsAvailable() (r bool) {
	var ePath string
	if runtime.GOOS == "windows" {
		ePath = os.Getenv("Path")
	} else {
		ePath = os.Getenv("PATH")
	}
	if strings.Contains(ePath, filepath.Join(config.DefaultGoRoot, "bin")) {
		r = true
	}
	return
}

func (that *GoVersion) getCurrent() (current string) {
	vFile := filepath.Join(config.DefaultGoRoot, "VERSION")
	if ok, _ := utils.PathIsExist(vFile); ok {
		if data, err := os.ReadFile(vFile); err == nil {
			return strings.TrimLeft(string(data), "go")
		}
	}
	return
}

func (that *GoVersion) ShowInstalled() {
	current := that.getCurrent()
	installedList, err := os.ReadDir(config.GoUnTarFilesPath)
	if err != nil {
		fmt.Println("[Read dir failed] ", err)
		return
	}
	for _, v := range installedList {
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

func (that *GoVersion) parseTarFileName(name string) (v string) {
	v = strings.Split(name, "-")[1]
	return
}

func (that *GoVersion) RemoveUnused() {
	current := that.getCurrent()
	installedList, err := os.ReadDir(config.GoUnTarFilesPath)
	if err != nil {
		fmt.Println("[Read dir failed] ", err)
		return
	}
	tarFiles, _ := os.ReadDir(config.GoTarFilesPath)
	for _, v := range installedList {
		if current == v.Name() {
			continue
		}
		os.RemoveAll(filepath.Join(config.GoUnTarFilesPath, v.Name()))
		for _, vInfo := range tarFiles {
			if v.Name() == that.parseTarFileName(vInfo.Name()) {
				os.Remove(filepath.Join(config.GoTarFilesPath, vInfo.Name()))
			}
		}
	}
}

func (that *GoVersion) RemoveVersion(version string) {
	current := that.getCurrent()
	if current != version {
		tarFiles, _ := os.ReadDir(config.GoTarFilesPath)
		os.RemoveAll(filepath.Join(config.GoUnTarFilesPath, version))
		for _, vInfo := range tarFiles {
			if version == that.parseTarFileName(vInfo.Name()) {
				os.Remove(filepath.Join(config.GoTarFilesPath, vInfo.Name()))
			}
		}
	}
}
