package confs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

/*
Verifier port range
*/
type VPortRange struct {
	Start int `koanf:"start"`
	End   int `koanf:"end"`
}

type XtrayConf struct {
	FetcherUrl        string      `koanf:"fetcher_url"`
	WorkDir           string      `koanf:"work_dir"`
	RawProxyFile      string      `koanf:"raw_file"`
	PorxyFile         string      `koanf:"proxy_file"`
	PortRange         *VPortRange `koanf:"port_range"`
	Port              int         `koanf:"port"`
	TestUrl           string      `koanf:"test_url"`
	SwitchyOmegaUrl   string      `koanf:"omega_url"`
	GeoInfoUrl        string      `koanf:"geo_info_url"`
	Timeout           int         `koanf:"timeout"`
	VerifierCron      string      `koanf:"verifier_cron"`
	KeeperCron        string      `koanf:"keeper_cron"`
	StorageSqlitePath string      `koanf:"storage_sqlite_path"`
	StorageExportPath string      `koanf:"storage_export_path"`
	path              string
}

func NewXtrayConf() (r *XtrayConf) {
	r = &XtrayConf{
		path: ProxyFilesDir,
	}
	r.setup()
	return
}

func (that *XtrayConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *XtrayConf) Reset() {
	that.WorkDir = that.path
	that.FetcherUrl = "https://gitee.com/moqsien/test/raw/master/conf.txt"
	that.RawProxyFile = filepath.Join(that.WorkDir, "raw_proxy.json")
	that.PorxyFile = filepath.Join(that.WorkDir, "latest.json")
	that.PortRange = &VPortRange{2020, 2075}
	that.Port = 2019
	that.TestUrl = "https://www.google.com"
	that.SwitchyOmegaUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/switch-omega.zip"
	that.GeoInfoUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/geoinfo.zip"
	that.Timeout = 3
	// "@every 1h30m10s" https://pkg.go.dev/github.com/robfig/cron
	that.VerifierCron = "@every 2h"
	that.KeeperCron = "@every 3m"

	that.StorageSqlitePath = filepath.Join(ProxyFilesDir, "storage.db")
	that.StorageExportPath = filepath.Join(GVCBackupDir, "vpn_history.json")
}
