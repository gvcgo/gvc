package confs

import (
	"os"
	"path/filepath"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
	neoconf "github.com/moqsien/neobox/pkgs/conf"
)

type NeoboxConf struct {
	NeoConf *neoconf.NeoBoxConf `koanf:"neobox_conf"`
	path    string
}

func NewNeoboxConf() (r *NeoboxConf) {
	r = &NeoboxConf{
		path:    ProxyFilesDir,
		NeoConf: &neoconf.NeoBoxConf{},
	}
	r.setup()
	return
}

func (that *NeoboxConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *NeoboxConf) Reset() {
	if that.NeoConf == nil {
		that.NeoConf = &neoconf.NeoBoxConf{}
	}
	that.NeoConf.NeoWorkDir = that.path
	that.NeoConf.NeoLogFileDir = filepath.Join(that.path, "neobox_logs")
	that.NeoConf.AssetDir = that.path
	that.NeoConf.XLogFileName = "neobox_client.log"
	that.NeoConf.SockFilesDir = that.path
	that.NeoConf.RawUriURL = "https://gitlab.com/moqsien/neobox_resources/-/raw/main/conf.txt"
	that.NeoConf.RawUriFileName = "neobox_raw_proxies.json"
	that.NeoConf.ParsedFileName = "neobox_parsed_proxies.json"
	that.NeoConf.PingedFileName = "neobox_pinged_proxies.json"
	that.NeoConf.MaxPingers = 100
	that.NeoConf.MaxAvgRTT = 600
	that.NeoConf.VerifiedFileName = "neobox_verified_proxies.json"
	that.NeoConf.VerifierPortRange = &neoconf.PortRange{
		Min: 4000,
		Max: 4050,
	}
	that.NeoConf.VerificationUri = "https://www.google.com"
	that.NeoConf.VerificationTimeout = 3
	that.NeoConf.VerificationCron = "@every 2h"
	that.NeoConf.NeoBoxClientInPort = 2019
	that.NeoConf.GeoInfoUrls = map[string]string{
		"geoip.dat":   "https://gitlab.com/moqsien/neobox_resources/-/raw/main/geoip.dat",
		"geosite.dat": "https://gitlab.com/moqsien/neobox_resources/-/raw/main/geosite.dat",
		"geoip.db":    "https://gitlab.com/moqsien/neobox_resources/-/raw/main/geoip.db",
		"geosite.db":  "https://gitlab.com/moqsien/neobox_resources/-/raw/main/geosite.db",
	}
	that.NeoConf.NeoBoxKeeperCron = "@every 3m"
	that.NeoConf.HistoryVpnsFileDir = GVCBackupDir
}
