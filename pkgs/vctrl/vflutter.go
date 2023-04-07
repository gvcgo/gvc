package vctrl

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
)

type FlutterPackage struct {
	Url         string
	FileName    string
	OS          string
	Arch        string
	DartVersion string
	Checksum    string
}

type FlutterVersion struct {
	Versions map[string]*FlutterPackage
	Json     *gjson.Json
	Conf     *config.GVConfig
	c        *colly.Collector
	d        *downloader.Downloader
	env      *utils.EnvsHandler
	baseUrl  string
}

func NewFlutterVersion() (fv *FlutterVersion) {
	fv = &FlutterVersion{
		Versions: make(map[string]*FlutterPackage, 500),
		Conf:     config.New(),
		c:        colly.NewCollector(),
		d:        &downloader.Downloader{},
		env:      utils.NewEnvsHandler(),
	}
	fv.initeDirs()
	return
}

func (that *FlutterVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.FlutterRootDir); !ok {
		os.RemoveAll(config.FlutterRootDir)
		if err := os.MkdirAll(config.FlutterRootDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterTarFilePath); !ok {
		if err := os.MkdirAll(config.FlutterTarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterUntarFilePath); !ok {
		if err := os.MkdirAll(config.FlutterUntarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *FlutterVersion) getJson() {
	var fUrl string
	switch runtime.GOOS {
	case utils.Windows:
		fUrl = that.Conf.Flutter.WinUrl
	case utils.Linux:
		fUrl = that.Conf.Flutter.LinuxUrl
	case utils.MacOS:
		fUrl = that.Conf.Flutter.MacosUrl
	default:
	}
	if !utils.VerifyUrls(fUrl) {
		return
	}
	that.c.OnResponse(func(r *colly.Response) {
		that.Json = gjson.New(r.Body)
	})
	that.c.Visit(fUrl)
	if that.Json != nil {
		that.baseUrl = that.Json.GetString("base_url")
	}
}

func (that *FlutterVersion) GetFileSuffix(fName string) string {
	for _, k := range AllowedSuffixes {
		if strings.HasSuffix(fName, k) {
			return k
		}
	}
	return ""
}

func (that *FlutterVersion) GetVersions() {
	if that.Json == nil {
		that.getJson()
	}
	if that.Json != nil {
		rList := that.Json.GetArray("releases")
		for _, release := range rList {
			j := gjson.New(release)
			rChannel := j.GetString("channel")
			version := j.GetString("version")
			if rChannel != "stable" || version == "" || strings.Contains(version, "+") {
				continue
			}

			p := &FlutterPackage{}
			p.Url = j.GetString("archive")
			p.OS = runtime.GOOS
			p.Arch = utils.ParseArch(j.GetString("dart_sdk_arch"))
			p.DartVersion = j.GetString("dart_sdk_version")
			p.Checksum = j.GetString("sha256")
			p.FileName = fmt.Sprintf("flutter-%s-%s-%s%s",
				version, p.OS, p.Arch, that.GetFileSuffix(p.Url))
			that.Versions[version] = p
		}
	}
}

func (that *FlutterVersion) ShowVersions() {
	if len(that.Versions) == 0 {
		that.GetVersions()
	}
	vList := []string{}
	for k := range that.Versions {
		vList = append(vList, k)
	}
	res := sorts.SortGoVersion(vList)
	fmt.Println(strings.Join(res, "  "))
}

func (that *FlutterVersion) download(version string) (r string) {
	if len(that.Versions) == 0 || that.baseUrl == "" {
		that.GetVersions()
	}

	if p := that.Versions[version]; p != nil {
		that.d.Url, _ = url.JoinPath(that.baseUrl, p.Url)
		if !utils.VerifyUrls(that.d.Url) {
			return
		}
		that.d.Timeout = 100 * time.Minute
		fpath := filepath.Join(config.FlutterTarFilePath, p.FileName)
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

func (that *FlutterVersion) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		flutterEnv := fmt.Sprintf(utils.FlutterEnv,
			config.FlutterRootDir,
			that.Conf.Flutter.HostedUrl,
			that.Conf.Flutter.StorageBaseUrl,
			that.Conf.Flutter.GitUrl)
		that.env.UpdateSub(utils.SUB_FLUTTER, flutterEnv)
	} else {
		envList := map[string]string{
			"PUB_HOSTED_URL":           that.Conf.Flutter.HostedUrl,
			"FLUTTER_STORAGE_BASE_URL": that.Conf.Flutter.StorageBaseUrl,
			"FLUTTER_GIT_URL":          that.Conf.Flutter.GitUrl,
			"PATH":                     filepath.Join(config.FlutterRootDir, "bin"),
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *FlutterVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.FlutterUntarFilePath, version)
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterRootDir); ok {
		os.RemoveAll(config.FlutterRootDir)
	}
	finder := utils.NewBinaryFinder(untarfile, "bin")
	dir := finder.String()
	if dir != "" {
		if err := utils.MkSymLink(dir, config.FlutterRootDir); err != nil {
			fmt.Println("[Create link failed] ", err)
			return
		}
		if !that.env.DoesEnvExist(utils.SUB_FLUTTER) {
			that.CheckAndInitEnv()
		}
		fmt.Println("Use", version, "succeeded!")
	}
}
