package vctrl

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/pterm/pterm"
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
	Versions map[string][]*FlutterPackage
	Json     *gjson.Json
	Conf     *config.GVConfig
	fetcher  *query.Fetcher
	env      *utils.EnvsHandler
	baseUrl  string
}

func NewFlutterVersion() (fv *FlutterVersion) {
	fv = &FlutterVersion{
		Versions: make(map[string][]*FlutterPackage, 500),
		Conf:     config.New(),
		fetcher:  query.NewFetcher(),
		env:      utils.NewEnvsHandler(),
	}
	fv.initeDirs()
	fv.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *FlutterVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.FlutterRootDir); !ok {
		os.RemoveAll(config.FlutterRootDir)
		if err := os.MkdirAll(config.FlutterRootDir, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterTarFilePath); !ok {
		if err := os.MkdirAll(config.FlutterTarFilePath, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterUntarFilePath); !ok {
		if err := os.MkdirAll(config.FlutterUntarFilePath, os.ModePerm); err != nil {
			tui.PrintError(err)
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

	that.fetcher.Url = fUrl
	if resp := that.fetcher.Get(); resp != nil {
		content, _ := io.ReadAll(resp.RawBody())
		that.Json = gjson.New(content)
	}
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
			if rChannel != "stable" || version == "" || strings.Contains(version, "hotfix") {
				continue
			}

			p := &FlutterPackage{}
			p.Url = j.GetString("archive")
			p.Arch = utils.ParseArch(j.GetString("dart_sdk_arch"))
			if p.Url == "" || p.Arch == "" {
				continue
			}
			p.OS = runtime.GOOS
			p.DartVersion = j.GetString("dart_sdk_version")
			p.Checksum = j.GetString("sha256")
			p.FileName = fmt.Sprintf("flutter-%s-%s-%s%s",
				version, p.OS, p.Arch, that.GetFileSuffix(p.Url))
			if len(that.Versions[version]) == 0 {
				that.Versions[version] = []*FlutterPackage{p}
			} else {
				that.Versions[version] = append(that.Versions[version], p)
			}
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
	fc := tui.NewFadeColors(res)
	fc.Println()
}

func (that *FlutterVersion) findPackage(version string) *FlutterPackage {
	for _, pk := range that.Versions[version] {
		if pk.Arch == runtime.GOARCH && pk.OS == runtime.GOOS {
			return pk
		}
	}
	return nil
}

func (that *FlutterVersion) download(version string) (r string) {
	if len(that.Versions) == 0 || that.baseUrl == "" {
		that.GetVersions()
	}

	if p := that.findPackage(version); p != nil {
		that.fetcher.Url, _ = url.JoinPath(that.baseUrl, p.Url)
		if !utils.VerifyUrls(that.fetcher.Url) {
			return
		}
		that.fetcher.Timeout = 100 * time.Minute
		fpath := filepath.Join(config.FlutterTarFilePath, p.FileName)
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
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
		tui.PrintError(fmt.Sprintf("Invalid Flutter version: %s", version))
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
				tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.FlutterRootDir); ok {
		os.RemoveAll(config.FlutterRootDir)
	}
	finder := utils.NewBinaryFinder(untarfile, "", "version")
	dir := finder.String()
	if dir != "" {
		if err := utils.MkSymLink(dir, config.FlutterRootDir); err != nil {
			tui.PrintError(fmt.Sprintf("Create link failed: %+v", err))
			return
		}
		if !that.env.DoesEnvExist(utils.SUB_FLUTTER) {
			that.CheckAndInitEnv()
		}
		tui.PrintSuccess(fmt.Sprintf("Use %s succeeded!", version))
	}
}

func (that *FlutterVersion) getCurrent() string {
	content, _ := os.ReadFile(filepath.Join(config.FlutterRootDir, "version"))
	return strings.TrimSpace(string(content))
}

func (that *FlutterVersion) ShowInstalled() {
	current := that.getCurrent()
	dList, _ := os.ReadDir(config.FlutterUntarFilePath)
	for _, d := range dList {
		if d.IsDir() {
			switch d.Name() {
			case current:
				fmt.Println(pterm.Yellow(fmt.Sprintf("%s <Current>", d.Name())))
			default:
				fmt.Println(pterm.Cyan(d.Name()))
			}
		}
	}
}

func (that *FlutterVersion) removeTarFile(version string) {
	fName := fmt.Sprintf("flutter-%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
	dList, _ := os.ReadDir(config.FlutterTarFilePath)
	for _, d := range dList {
		if !d.IsDir() && strings.Contains(d.Name(), fName) {
			os.RemoveAll(filepath.Join(config.FlutterTarFilePath, d.Name()))
		}
	}
}

func (that *FlutterVersion) RemoveVersion(version string) {
	current := that.getCurrent()
	if version == current {
		return
	}
	dList, _ := os.ReadDir(config.FlutterUntarFilePath)
	for _, d := range dList {
		if d.IsDir() && d.Name() == version {
			os.RemoveAll(filepath.Join(config.FlutterUntarFilePath, d.Name()))
			that.removeTarFile(version)
		}
	}
}

func (that *FlutterVersion) RemoveUnused() {
	current := that.getCurrent()
	dList, _ := os.ReadDir(config.FlutterUntarFilePath)
	for _, d := range dList {
		if d.IsDir() && d.Name() != current {
			os.RemoveAll(filepath.Join(config.FlutterUntarFilePath, d.Name()))
			that.removeTarFile(d.Name())
		}
	}
}
